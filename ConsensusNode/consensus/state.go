package consensus

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/p2pnetwork"
	"consensusNode/signature"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Consensus interface {
	StartConsensus()
	PrePrepare()
	Prepare()
	Commit()
}

type Stage int

// number different kinds of stage types
const (
	Idle Stage = iota
	PrePrepared
	Prepared
	Committed
)

// stage.string()
func (s Stage) String() string {
	switch s {
	case Idle:
		return "Idle"
	case PrePrepared:
		return "PrePrepared"
	case Prepared:
		return "Prepared"
	case Committed:
		return "Committed"
	}
	return "Unknown"
}

// timer
const StateTimerOut = 5 * time.Second
const MaxStateMsgNO = 100

type RequestTimer struct {
	*time.Ticker
	IsOk bool
}

// initialize timer
func newRequestTimer() *RequestTimer {
	tick := time.NewTicker(StateTimerOut)
	tick.Stop()
	return &RequestTimer{
		Ticker: tick,
		IsOk:   false,
	}
}

// start a timer
func (rt *RequestTimer) tick() {
	if rt.IsOk {
		return
	}
	rt.Reset(StateTimerOut)
	rt.IsOk = true
}

// stop a timer
func (rt *RequestTimer) tack() {
	rt.IsOk = false
	rt.Stop()
}

type EngineStatus int8

// number different kinds of EngineStatus types
const (
	Serving EngineStatus = iota
	ViewChanging
)

// EngineStatus.string()
func (es EngineStatus) String() string {
	switch es {
	case Serving:
		return "Server consensus......"
	case ViewChanging:
		return "Changing views......"
	}
	return "Unknown"
}

type StateEngine struct {
	NodeID      int64 `json:"nodeID"`
	CurViewID   int64 `json:"viewID"`
	CurSequence int64 `json:"curSeq"`
	PrimaryID   int64 `json:"primaryID"`
	nodeStatus  EngineStatus
	SrvHub      *net.TCPListener

	Timer    *RequestTimer
	P2pWire  p2pnetwork.P2pNetwork
	MsgChan  <-chan *message.ConMessage
	nodeChan chan<- *message.RequestRecord

	MiniSeq int64 `json:"miniSeq"`
	MaxSeq  int64 `json:"maxSeq"`
	//index is seqID
	msgLogs        map[int64]*NormalLog
	TimeCommitment map[string]string
	sCache         *VCCache
}

func InitConsensus(id int64, cChan chan<- *message.RequestRecord) *StateEngine {
	fmt.Printf("===>Service is Listening at[%d]\n", util.PortByID(id))

	ch := make(chan *message.ConMessage, MaxStateMsgNO)
	p2p := p2pnetwork.NewSimpleP2pLib(id, ch)
	se := &StateEngine{
		NodeID:         id,
		CurViewID:      0,
		CurSequence:    0,
		MiniSeq:        0,
		MaxSeq:         64,
		Timer:          newRequestTimer(),
		P2pWire:        p2p,
		MsgChan:        ch,
		nodeChan:       cChan,
		msgLogs:        make(map[int64]*NormalLog),
		sCache:         NewVCCache(),
		SrvHub:         new(net.TCPListener),
		TimeCommitment: make(map[string]string),
	}
	se.PrimaryID = se.CurViewID % message.TotalNodeNum

	if se.PrimaryID == se.NodeID {
		go se.WriteRandomOutput()
	}

	return se
}

func (s *StateEngine) WriteRandomOutput() {
	time.Sleep(25 * time.Second)
	fmt.Println("start wirte config")
	// generate random init input
	message := []byte("hello world")
	randomNum := signature.Digest(message)

	config.WriteOutput(string(randomNum))
}

// receive and handle consensus message
func (s *StateEngine) StartConsensus(sig chan interface{}) {
	s.nodeStatus = Serving

	for {
		select {
		case <-s.Timer.C:
			s.Timer.tack()
			if s.NodeID != s.PrimaryID {
				go s.sendUnionMsg()
			}
			fmt.Println(time.Now().Unix())
			fmt.Printf("======>[Node%d]Stop Receive messages\n", s.NodeID)
			for key, value := range s.TimeCommitment {
				fmt.Println("key is", key)
				fmt.Println("value is", value)
			}
		case conMsg := <-s.MsgChan:
			switch conMsg.Typ {
			case message.MTSubmit,
				message.MTApprove:
				if s.nodeStatus != Serving {
					fmt.Println("======>[ERROR]node is not in service status now......")
					continue
				}
				if err := s.procConsensusMsg(conMsg); err != nil {
					fmt.Println(err)
				}
			case message.MTViewChange,
				message.MTNewView:
				// if err := s.procManageMsg(conMsg); err != nil {
				// 	fmt.Println(err)
				// }
			}
		}
	}
}

// wait for request from client in UDP channel
func (s *StateEngine) WaitTC(sig chan interface{}) {
	defer func() {
		if r := recover(); r != nil {
			sig <- r
		}
	}()

	buf := make([]byte, 2048)
	for {
		conn, err := s.SrvHub.AcceptTCP()
		if err != nil {
			// fmt.Printf("===>[ERROR]TCP:%s\n", err)
			continue
		}
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("===>[ERROR]Service received err:%s\n", err)
			continue
		}

		// get message from entropy node
		msgFromEntropyNode := &message.ConMessage{}
		if err := json.Unmarshal(buf[:n], msgFromEntropyNode); err != nil {
			fmt.Printf("===>[ERROR]Service message parse err:%s\n", err)
			continue
		}
		entropyMsg := &message.EntropyMessage{}
		if err := json.Unmarshal(msgFromEntropyNode.Payload, entropyMsg); err != nil {
			fmt.Printf("===>[ERROR] Invalid[%s] Entropy message[%s]", err, msgFromEntropyNode)
		}

		// get entropy node's public key and verify signature
		pub, err := x509.ParsePKIXPublicKey(entropyMsg.PublicKey)
		if err != nil {
			fmt.Printf("===>[ERROR]Key message parse err:%s", err)
			continue
		}
		entropyPK := pub.(*ecdsa.PublicKey)
		verify := signature.VerifySig(msgFromEntropyNode.Payload, msgFromEntropyNode.Sig, entropyPK)
		if !verify {
			fmt.Printf("===>[ERROR]Verify new entropy message Signature failed, From Entropy Node[%d]\n", entropyMsg.ClientID)
			continue
		}
		fmt.Printf("======>[NewEntropyMsg]Verify success\n")

		// verify the VRF result
		previousOutput := config.GetPreviousInput()
		verify = signature.VerifySig([]byte(previousOutput), entropyMsg.VRFResult, entropyPK)
		if !verify {
			fmt.Printf("===>[ERROR]Verify new VRF result failed, From Entropy Node[%d]\n", entropyMsg.ClientID)
			continue
		}
		fmt.Printf("======>[VRFresult]Verify success\n")

		// verify the VRF and difficuly
		difficulty := config.GetDifficulty()
		vrfResultBinary := util.BytesToBinaryString(entropyMsg.VRFResult)
		vrfResultTail, err := strconv.Atoi(vrfResultBinary[len(vrfResultBinary)-1:])
		if err != nil {
			fmt.Printf("===>[ERROR]Get vrfResultTail, From Entropy Node[%d]\n", entropyMsg.ClientID)
			continue
		}
		if vrfResultTail != difficulty {
			fmt.Println("Cheater!!!!")
			fmt.Printf("===>[ERROR]Verify difficulty failed, From Entropy Node[%d]\n", entropyMsg.ClientID)
			continue
		}

		s.TimeCommitment[string(util.Digest(entropyMsg.TimeCommitment))] = entropyMsg.TimeCommitment
		// s.TimeCommitment = append(s.TimeCommitment, entropyMsg.TimeCommitment)
		fmt.Printf("===>Entropy message from Node[%d],time commitment is:%s\n", entropyMsg.ClientID, entropyMsg.TimeCommitment)
	}
}

// watch config
// when previousout changes, start a new round
func (s *StateEngine) WatchConfig(id int64, sig chan interface{}) {
	previousOutput := string(config.GetPreviousInput())
	fmt.Println("init output", previousOutput)

	// set config file
	viper.SetConfigFile("../config.yml")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// lock file
		f, err := os.Open("../lock")
		if err != nil {
			panic(err)
		}
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
			log.Println("add share lock in no block failed", err)
		}
		fmt.Println(time.Now())
		fmt.Println("Config Change")

		// 	// time.Sleep(1 * time.Second)
		// time.Sleep(200 * time.Millisecond)
		// 	// config.ReadConfig()
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		// newnewConfig := viper.AllSettings()
		// fmt.Printf("All settings #4 %+v\n\n", newnewConfig)

		newOutput := string(config.GetPreviousInput())
		if previousOutput != newOutput && newOutput != "" {
			fmt.Println("output change", newOutput)

			locAddr := net.TCPAddr{
				Port: util.EntropyPortByID(id),
			}
			srvHub, err := net.ListenTCP("tcp4", &locAddr)
			if err != nil {
				panic(err)
			}
			s.SrvHub = srvHub
			s.SrvHub.SetDeadline(time.Now().Add(5 * time.Second))

			go s.WaitTC(sig)
			previousOutput = newOutput
			s.Timer.tick()
		}

		// unlock file
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
			log.Println("unlock share lock failed", err)
		}
	})
}

// backups send union message
func (s *StateEngine) sendUnionMsg() {
	// new submit message
	// send submit message to primary node
	var tc []string
	for _, value := range s.TimeCommitment {
		tc = append(tc, value)
	}
	submit := &message.Submit{
		CollectTC: tc,
	}
	sk := s.P2pWire.GetMySecretkey()
	sMsg := message.CreateConMsg(message.MTSubmit, submit, sk, s.NodeID)
	conn := s.P2pWire.GetPrimaryConn(s.PrimaryID)
	if err := s.P2pWire.SendUniqueNode(conn, sMsg); err != nil {
		panic(err)
	}
}

// primary union TC
func (s *StateEngine) unionTC(msg *message.ConMessage) (err error) {
	// time.Sleep(1 * time.Second)
	fmt.Printf("======>[Union]Current Union Message from Node[%d]\n", msg.From)

	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[UnionERROR]Key message parse err:%s", err))
	}
	normalPublicKey := pub.(*ecdsa.PublicKey)
	curve := normalPublicKey.Curve

	// unmarshal public key
	x, y := elliptic.Unmarshal(curve, []byte(nodeConfig[msg.From].Pk))
	newPublicKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	// verify signature
	fmt.Println("payload", msg.Payload)
	fmt.Println("pk", nodeConfig[msg.From].Pk)
	verify := signature.VerifySig(msg.Payload, msg.Sig, newPublicKey)
	if !verify {
		fmt.Printf("===>[UnionERROR]Verify new public key Signature failed, From Node[%d]\n", msg.From)
	}

	// unmarshal message
	submit := &message.Submit{}
	if err := json.Unmarshal(msg.Payload, submit); err != nil {
		return fmt.Errorf("======>[Union]Invalid[%s] Union message[%s]", err, msg)
	}

	fmt.Println("submit", submit)
	fmt.Printf("Union Message from Node[%d],length is %d\n", msg.From, len(submit.CollectTC))
	for _, value := range submit.CollectTC {
		key := string(util.Digest(value))
		if _, ok := s.TimeCommitment[key]; !ok {
			fmt.Println("new key is", key)
			fmt.Println("new value is", value)
			s.TimeCommitment[key] = value
		}
	}

	return
}

// handle different kinds of consensus messages
func (s *StateEngine) procConsensusMsg(msg *message.ConMessage) (err error) {
	fmt.Printf("\n======>[procConsensusMsg]Consesus message type:[%s] from Node[%d]\n", msg.Typ, msg.From)

	switch msg.Typ {
	case message.MTSubmit:
		return s.unionTC(msg)

		// case message.MTApprove:
		// 	prepare := &message.Prepare{}
		// 	if err := json.Unmarshal(msg.Payload, prepare); err != nil {
		// 		return fmt.Errorf("======>[procConsensusMsg]invalid[%s] Prepare message[%s]", err, msg)
		// 	}
		// 	return s.prePrepare2Prepare(prepare, msg)
		// }

	}
	return
}
