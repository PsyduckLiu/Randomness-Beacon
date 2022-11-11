package consensus

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/p2pnetwork"
	"consensusNode/signature"
	tc "consensusNode/timedCommitment"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/protocol"
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
	Collect Stage = iota
	Submit
	Approve
	Confirm
	Output
	Error
	ViewChange
)

// stage.string()
func (s Stage) String() string {
	switch s {
	case Collect:
		return "Collect"
	case Submit:
		return "Submit"
	case Approve:
		return "Approve"
	case Confirm:
		return "Confirm"
	case Output:
		return "Output"
	case Error:
		return "Error"
	case ViewChange:
		return "ViewChange"
	}

	return "Unknown"
}

// timer
const StateTimerOut = 6 * time.Second
const MaxStateMsgNO = 150

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
func (rt *RequestTimer) tick(time time.Duration) {
	if rt.IsOk {
		return
	}
	rt.Reset(time)
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

type TestingHashable struct {
	data []byte
}

func (s TestingHashable) ToBeHashed() (protocol.HashID, []byte) {
	return "test", s.data
}

type StateEngine struct {
	NodeID      int64 `json:"nodeID"`
	CurViewID   int64 `json:"viewID"`
	CurSequence int64 `json:"curSeq"`
	PrimaryID   int64 `json:"primaryID"`
	nodeStatus  EngineStatus
	SrvHub      *net.TCPListener

	GlobalTimer  *RequestTimer
	CollectTimer *RequestTimer
	SubmitTimer  *RequestTimer

	P2pWire  p2pnetwork.P2pNetwork
	MsgChan  <-chan *message.ConMessage
	nodeChan chan<- *message.RequestRecord

	SubmitNum             int64
	ConfirmNum            int64
	OutputNum             int64
	ApproveNum            int
	stage                 Stage
	Result                *big.Int
	Mutex                 sync.Mutex
	MiniSeq               int64 `json:"miniSeq"`
	MaxSeq                int64 `json:"maxSeq"`
	entropyNode           map[int64]bool
	msgLogs               map[int64]*NormalLog
	TimeCommitment        map[string][4]string
	TimeCommitmentSubmit  map[int64]int
	TimeCommitmentApprove map[string]bool
	sCache                *VCCache
}

func InitConsensus(id int64, cChan chan<- *message.RequestRecord) *StateEngine {
	fmt.Printf("===>Service is Listening at[%d]\n", util.PortByID(id))

	ch := make(chan *message.ConMessage, MaxStateMsgNO)
	p2p := p2pnetwork.NewSimpleP2pLib(id, ch)
	se := &StateEngine{
		NodeID:       id,
		CurViewID:    0,
		CurSequence:  0,
		SubmitNum:    0,
		OutputNum:    0,
		ConfirmNum:   0,
		ApproveNum:   0,
		stage:        Collect,
		Result:       big.NewInt(0),
		Mutex:        sync.Mutex{},
		MiniSeq:      0,
		MaxSeq:       64,
		GlobalTimer:  newRequestTimer(),
		CollectTimer: newRequestTimer(),
		SubmitTimer:  newRequestTimer(),
		P2pWire:      p2p,
		MsgChan:      ch,
		nodeChan:     cChan,
		entropyNode:  make(map[int64]bool),
		msgLogs:      make(map[int64]*NormalLog),
		sCache:       NewVCCache(),
		// SrvHub:         new(net.TCPListener),
		TimeCommitment:        make(map[string][4]string),
		TimeCommitmentSubmit:  make(map[int64]int),
		TimeCommitmentApprove: make(map[string]bool),
	}
	se.PrimaryID = se.CurViewID % util.TotalNodeNum

	if se.PrimaryID == se.NodeID {
		go se.WriteRandomOutput()
	}

	return se
}

func (s *StateEngine) WriteRandomOutput() {
	time.Sleep(30 * time.Second)
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
		case <-s.GlobalTimer.C:
			s.GlobalTimer.tack()
			fmt.Printf("======>[Node%d]Time is out and view change starts\n", s.NodeID)
			s.nodeStatus = ViewChanging
			s.stage = ViewChange
			s.ViewChange()
		case <-s.CollectTimer.C:
			s.CollectTimer.tack()
			s.stage = Submit
			if s.NodeID != s.PrimaryID {
				if s.PrimaryID == 0 && s.NodeID%2 == 1 {
					fmt.Println("Rotten")
				} else {
					go s.sendUnionMsg()
				}
			}

			fmt.Println(time.Now())
			fmt.Printf("======>[Node%d]Stop Receive messages\n", s.NodeID)
			for key, value := range s.TimeCommitment {
				fmt.Println("key is", key)
				fmt.Println("value is", value)
			}
		case <-s.SubmitTimer.C:
			fmt.Println("submit timer out,submit number is", s.SubmitNum)
			if s.SubmitNum >= 2*util.MaxFaultyNode+1 {
				s.SubmitTimer.tack()

				// new approve message
				// send approve message to backup nodes
				approve := &message.Approve{}
				result, _ := rand.Int(rand.Reader, big.NewInt(1000))
				// result := big.NewInt(0)
				fmt.Println(result)
				if result.Cmp(big.NewInt(0)) == 0 {
					fmt.Println("I'm crazy!!!!!")
					sk := s.P2pWire.GetMySecretkey()
					aMsg := message.CreateConMsg(message.MTApprove, approve, sk, s.NodeID)
					if err := s.P2pWire.BroadCast(aMsg); err != nil {
						panic(err)
					}
				} else {
					var tc [4]string
					for _, value := range s.TimeCommitment {
						tc = value
						approve = &message.Approve{
							Length:  len(s.TimeCommitment),
							UnionTC: tc,
						}
						sk := s.P2pWire.GetMySecretkey()
						aMsg := message.CreateConMsg(message.MTApprove, approve, sk, s.NodeID)
						if err := s.P2pWire.BroadCast(aMsg); err != nil {
							panic(err)
						}
						time.Sleep(150 * time.Millisecond)
					}
				}

				s.stage = Confirm
				s.ConfirmNum++
				fmt.Printf("======>[Union]Send approve message success\n")
			}
		case conMsg := <-s.MsgChan:
			switch conMsg.Typ {
			case message.MTSubmit,
				message.MTApprove,
				message.MTConfirm,
				message.MTOutput:
				if s.nodeStatus != Serving {
					fmt.Println("======>[ERROR]node is not in service status now......")
					continue
				}
				if err := s.procConsensusMsg(conMsg); err != nil {
					fmt.Println(err)
				}
			case message.MTViewChange,
				message.MTNewView:
				if err := s.procManageMsg(conMsg); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

// watch config
// when previousout changes, start a new round
func (s *StateEngine) WatchConfig(id int64, sig chan interface{}) {
	previousOutput := string(config.GetPreviousOutput())
	fmt.Println("init output", previousOutput)

	myViper := viper.New()
	// set config file
	myViper.SetConfigFile("../output.yml")
	myViper.WatchConfig()
	myViper.OnConfigChange(func(e fsnotify.Event) {
		// lock file
		f, err := os.Open("../lock")
		if err != nil {
			panic(err)
		}
		// if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH); err != nil {
			log.Println("add share lock in no block failed", err)
		}

		fmt.Println(time.Now())
		fmt.Println("Config Change")

		// config.ReadConfig()
		if err := myViper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		newOutput := string(config.GetPreviousOutput())
		if previousOutput != newOutput && newOutput != "" && s.stage == Collect {
			fmt.Println("output change", newOutput)

			s.stage = Collect
			s.OutputNum = 0
			s.ConfirmNum = 0
			s.SubmitNum = 0
			s.ApproveNum = 0
			s.TimeCommitment = make(map[string][4]string)
			s.TimeCommitmentSubmit = make(map[int64]int)
			s.TimeCommitmentApprove = make(map[string]bool)

			if s.SrvHub == nil {
				locAddr := net.TCPAddr{
					Port: util.EntropyPortByID(id),
				}
				srvHub, err := net.ListenTCP("tcp4", &locAddr)
				if err != nil {
					panic(err)
				}
				s.SrvHub = srvHub
			}
			s.SrvHub.SetDeadline(time.Now().Add(5 * time.Second))

			go s.WaitTC(sig)
			previousOutput = newOutput
			s.GlobalTimer.tick(35 * time.Second)
			s.CollectTimer.tick(5 * time.Second)
			if s.NodeID == s.PrimaryID {
				s.SubmitTimer.tick(20 * time.Second)
			}
		}

		// unlock file
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
			log.Println("unlock share lock failed", err)
		}
	})
}

// wait for request from client in UDP channel
func (s *StateEngine) WaitTC(sig chan interface{}) {
	defer func() {
		if r := recover(); r != nil {
			sig <- r
		}
	}()

	buf := make([]byte, 8192)
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
		// fmt.Println(msgFromEntropyNode)
		if msgFromEntropyNode.Typ != message.MTCollect && msgFromEntropyNode.Typ != message.MTVRFVerify {
			fmt.Printf("===>[ERROR]Not collect message:%s\n", err)
			continue
		}

		// handle vrf verify message
		if msgFromEntropyNode.Typ == message.MTVRFVerify {
			fmt.Printf("===>[WaitTC]new vrf verify message from Node[%d]\n", msgFromEntropyNode.From)

			entropyVRFMsg := &message.EntropyVRFMessage{}
			if err := json.Unmarshal(msgFromEntropyNode.Payload, entropyVRFMsg); err != nil {
				fmt.Printf("===>[ERROR] Invalid[%s] Entropy message[%s]", err, msgFromEntropyNode)
				continue
			}

			// verify the VRF result
			msg := TestingHashable{
				data: entropyVRFMsg.Msg,
			}
			var pk crypto.VrfPubkey = entropyVRFMsg.PublicKey
			verify, output := pk.Verify(entropyVRFMsg.VRFResult, msg)
			if !verify {
				fmt.Printf("===>[ERROR]Verify new VRF result failed, From Entropy Node[%d]\n", entropyVRFMsg.ClientID)
				continue
			}
			fmt.Println("VRF Output is", output)
			fmt.Printf("======>[VRFresult]Verify success\n")

			// verify the VRF and difficuly
			difficulty := config.GetDifficulty()
			vrfResultBinary := util.BytesToBinaryString(entropyVRFMsg.VRFResult)
			vrfResultTail, err := strconv.Atoi(vrfResultBinary[len(vrfResultBinary)-1:])
			if err != nil {
				fmt.Printf("===>[ERROR]Get vrfResultTail, From Entropy Node[%d]\n", entropyVRFMsg.ClientID)
				continue
			}
			if vrfResultTail != difficulty {
				fmt.Println("Cheater!!!!")
				fmt.Printf("===>[ERROR]Verify difficulty failed, From Entropy Node[%d]\n", entropyVRFMsg.ClientID)
				continue
			}

			// s.Mutex.Lock()
			// s.entropyNode[entropyVRFMsg.ClientID] = true
			// s.Mutex.Unlock()
		}

		// handle collect message
		if msgFromEntropyNode.Typ == message.MTCollect {
			fmt.Println("===>[WaitTC]new collect message")

			entropyTCMsg := &message.EntropyTCMessage{}
			if err := json.Unmarshal(msgFromEntropyNode.Payload, entropyTCMsg); err != nil {
				fmt.Printf("===>[ERROR] Invalid[%s] Entropy TC message[%s]", err, msgFromEntropyNode)
				continue
			}

			mArray := config.GetMArray()
			g := config.GetG()
			N := config.GetN()
			a1, _ := new(big.Int).SetString(entropyTCMsg.TimeCommitmentA1, 10)
			a2, _ := new(big.Int).SetString(entropyTCMsg.TimeCommitmentA2, 10)
			a3, _ := new(big.Int).SetString(entropyTCMsg.TimeCommitmentA3, 10)
			z, _ := new(big.Int).SetString(entropyTCMsg.TimeCommitmentZ, 10)
			h, _ := new(big.Int).SetString(entropyTCMsg.TimeCommitmentH, 10)
			rKSubOne, _ := new(big.Int).SetString(entropyTCMsg.TimeCommitmentrKSubOne, 10)
			rK, _ := new(big.Int).SetString(entropyTCMsg.TimeCommitmentrK, 10)

			nHash := new(big.Int).SetBytes(util.Digest(N))
			gHash := new(big.Int).SetBytes(util.Digest((g)))
			mSubOneHash := new(big.Int).SetBytes(util.Digest(mArray[len(mArray)-3]))
			mHash := new(big.Int).SetBytes(util.Digest(mArray[len(mArray)-2]))
			a1Hash := new(big.Int).SetBytes(util.Digest(a1))
			a2Hash := new(big.Int).SetBytes(util.Digest(a2))
			a3Hash := new(big.Int).SetBytes(util.Digest(a3))

			e := big.NewInt(0)
			e.Xor(e, gHash)
			e.Xor(e, nHash)
			e.Xor(e, mSubOneHash)
			e.Xor(e, mHash)
			e.Xor(e, a1Hash)
			e.Xor(e, a2Hash)
			e.Xor(e, a3Hash)

			result1 := new(big.Int).Set(g)
			result1.Exp(result1, z, N)
			result2 := new(big.Int).Set(h)
			result2.Exp(result2, e, N)
			result1.Mul(result1, result2)
			result1.Mod(result1, N)

			result3 := new(big.Int).Set(mArray[len(mArray)-3])
			result3.Exp(result3, z, N)
			result4 := new(big.Int).Set(rKSubOne)
			result4.Exp(result4, e, N)
			result3.Mul(result3, result4)
			result3.Mod(result3, N)

			result5 := new(big.Int).Set(mArray[len(mArray)-2])
			result5.Exp(result5, z, N)
			result6 := new(big.Int).Set(rK)
			result6.Exp(result6, e, N)
			result5.Mul(result5, result6)
			result5.Mod(result5, N)

			if a1.Cmp(result1) != 0 {
				fmt.Println("test1 error")
				continue
			}
			if a2.Cmp(result3) != 0 {
				fmt.Println("test2 error")
				continue
			}
			if a3.Cmp(result5) != 0 {
				fmt.Println("test3 error")
				continue
			}
			fmt.Println("[WaitTC]pass all tests!")

			// value, ok := s.entropyNode[entropyTCMsg.ClientID]
			// if !value || !ok {
			// 	fmt.Printf("===>[ERROR]Not verified")
			// 	continue
			// }
			s.Mutex.Lock()
			s.entropyNode[entropyTCMsg.ClientID] = true
			// s.Mutex.Unlock()

			// s.Mutex.Lock()
			timedCommitment := [4]string{entropyTCMsg.TimeCommitmentC, entropyTCMsg.TimeCommitmentH, entropyTCMsg.TimeCommitmentrKSubOne, entropyTCMsg.TimeCommitmentrK}
			s.TimeCommitment[string(util.Digest(timedCommitment))] = timedCommitment
			s.TimeCommitmentApprove[string(util.Digest(timedCommitment))] = false
			s.Mutex.Unlock()

			fmt.Printf("===>Entropy message from Node[%d]\n", entropyTCMsg.ClientID)
		}
	}
}

// backups send union message
func (s *StateEngine) sendUnionMsg() {
	// new submit message
	// send submit message to primary node
	var tc [4]string
	for _, value := range s.TimeCommitment {
		tc = value
		submit := &message.Submit{
			Length:    len(s.TimeCommitment),
			CollectTC: tc,
		}
		sk := s.P2pWire.GetMySecretkey()
		sMsg := message.CreateConMsg(message.MTSubmit, submit, sk, s.NodeID)
		conn := s.P2pWire.GetPrimaryConn(s.PrimaryID)
		if err := s.P2pWire.SendUniqueNode(conn, sMsg); err != nil {
			panic(err)
		}
		time.Sleep(150 * time.Millisecond)
	}

	s.stage = Approve
}

// primary union TC
func (s *StateEngine) unionTC(msg *message.ConMessage) (err error) {
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
	verify := signature.VerifySig(msg.Payload, msg.Sig, newPublicKey)
	if !verify {
		panic(fmt.Errorf("===>[UnionERROR]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	submit := &message.Submit{}
	if err := json.Unmarshal(msg.Payload, submit); err != nil {
		panic(fmt.Errorf("======>[Union]Invalid[%s] Union message[%s]", err, msg))
	}

	fmt.Println("submit", submit)
	fmt.Printf("Union Message from Node[%d],length is %d\n", msg.From, submit.Length)

	key := string(util.Digest(submit.CollectTC))
	if value, ok := s.TimeCommitment[key]; !ok {
		fmt.Println("new key is", key)
		fmt.Println("new TimeCommitmentC value is", value[0])
		s.TimeCommitment[key] = submit.CollectTC
		s.TimeCommitmentApprove[key] = false
	}

	if _, ok := s.TimeCommitmentSubmit[msg.From]; !ok {
		fmt.Printf("node[%d] is sending submit messages\n", msg.From)
		s.TimeCommitmentSubmit[msg.From] = 0
	}
	s.TimeCommitmentSubmit[msg.From]++
	if s.TimeCommitmentSubmit[msg.From] == submit.Length {
		s.SubmitNum++
	}

	return
}

// backups check union tc sent by primary
func (s *StateEngine) approveTC(msg *message.ConMessage) (err error) {
	fmt.Printf("======>[Approve]Current Approve Message from Node[%d]\n", msg.From)
	fmt.Println("my length is", len(s.TimeCommitment))

	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ApproveERROR]Key message parse err:%s", err))
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
	verify := signature.VerifySig(msg.Payload, msg.Sig, newPublicKey)
	if !verify {
		panic(fmt.Errorf("===>[ApproveERROR]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	approve := &message.Approve{}
	if err := json.Unmarshal(msg.Payload, approve); err != nil {
		panic(fmt.Errorf("======>[Approve]Invalid[%s] Approve message[%s]", err, msg))
	}

	fmt.Println("approve", approve)
	fmt.Printf("Approve Message from Node[%d],length is %d\n", msg.From, approve.Length)
	if approve.Length == 0 {
		return fmt.Errorf("===>[ApproveERROR]are you sure????")
	}
	for _, value := range approve.UnionTC {
		fmt.Println("[UnionTC] value is", value)
	}

	// check union tc set
	key := string(util.Digest(approve.UnionTC))
	if _, ok := s.TimeCommitment[key]; !ok {
		fmt.Println("new key is", key)
		s.TimeCommitment[key] = approve.UnionTC
	}
	s.TimeCommitmentApprove[key] = true

	s.ApproveNum++
	if s.ApproveNum == approve.Length {
		for _, flag := range s.TimeCommitmentApprove {
			if !flag {
				s.stage = Error
				return fmt.Errorf("===>[ApproveERROR]where is my tc????")
			}
		}
		s.ConfirmNum++
		s.sendConfirmMsg()
	}

	// go s.handleTC()
	return
}

// backups broadcast confirm message
func (s *StateEngine) sendConfirmMsg() {
	confirm := &message.Confirm{
		Length: len(s.TimeCommitment),
	}

	sk := s.P2pWire.GetMySecretkey()
	cMsg := message.CreateConMsg(message.MTConfirm, confirm, sk, s.NodeID)

	if err := s.P2pWire.BroadCast(cMsg); err != nil {
		panic(err)
	}

	s.stage = Confirm
	fmt.Printf("======>[Approve]Send confirm message success\n")
}

// confirm TC
func (s *StateEngine) confirmTC(msg *message.ConMessage) (err error) {
	if s.ConfirmNum >= 2*util.MaxFaultyNode+1 {
		fmt.Println("late!")
		return
	}

	fmt.Printf("======>[Confirm]Current Confirm Message from Node[%d]\n", msg.From)
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ConfirmERROR]Key message parse err:%s", err))
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
	verify := signature.VerifySig(msg.Payload, msg.Sig, newPublicKey)
	if !verify {
		panic(fmt.Errorf("===>[ConfirmERROR]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	confirm := &message.Confirm{}
	if err := json.Unmarshal(msg.Payload, confirm); err != nil {
		panic(fmt.Errorf("======>[Confirm]Invalid[%s] Confirm message[%s]", err, msg))
	}

	// if len(confirm.ConfirmTC) != len(s.TimeCommitment) {
	// 	fmt.Printf("======>[Confirm]Invalid tc set\n")
	// 	return
	// }

	// confirmTC := make(map[string][4]string)
	// for _, value := range confirm.ConfirmTC {
	// 	key := string(util.Digest(value))
	// 	confirmTC[key] = value
	// }
	// for key := range s.TimeCommitment {
	// 	if _, ok := confirmTC[key]; !ok {
	// 		panic(fmt.Errorf("======>[Confirm]Invalid tc set"))
	// 	}
	// }

	if confirm.Length != len(s.TimeCommitment) {
		fmt.Printf("======>[Confirm]Invalid tc set\n")
		return
	}

	s.ConfirmNum++
	if s.ConfirmNum == 2*util.MaxFaultyNode+1 {
		fmt.Println("[Confirm]Confirm success")
		s.stage = Output
		go s.handleTC()
	}

	return
}

// resolve tc and compute F
func (s *StateEngine) handleTC() (err error) {
	var resolvedTC []*big.Int
	for _, value := range s.TimeCommitment {
		// TODO:resolve tc
		timeParameter := 10
		c := new(big.Int)
		c.UnmarshalJSON([]byte(value[0]))
		h := new(big.Int)
		h.UnmarshalJSON([]byte(value[1]))
		rKSubOne := new(big.Int)
		rKSubOne.UnmarshalJSON([]byte(value[2]))
		rK := new(big.Int)
		rK.UnmarshalJSON([]byte(value[3]))

		result := tc.ForcedOpen(c, h, rKSubOne, rK, timeParameter)
		resolvedTC = append(resolvedTC, result)
	}

	// Xor all resolved tc
	result := big.NewInt(0)
	for _, value := range resolvedTC {
		result.Xor(result, value)
	}

	s.Result = result
	fmt.Println("after xor all resolved tc,result is:", s.Result)

	if s.NodeID != s.PrimaryID {
		sk := s.P2pWire.GetMySecretkey()
		oMsg := message.CreateConMsg(message.MTOutput, s.Result.String(), sk, s.NodeID)
		conn := s.P2pWire.GetPrimaryConn(s.PrimaryID)
		time.Sleep(200 * time.Millisecond)
		if err := s.P2pWire.SendUniqueNode(conn, oMsg); err != nil {
			panic(err)
		}
		s.stage = Collect
		s.GlobalTimer.tack()
	}

	return
}

// output tc and compute F
func (s *StateEngine) outputTC(msg *message.ConMessage) (err error) {
	if s.OutputNum >= 2*util.MaxFaultyNode {
		fmt.Println("late!")
		return
	}

	fmt.Printf("======>[Output]Current Output Message from Node[%d]\n", msg.From)
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[OutputERROR]Key message parse err:%s", err))
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
	verify := signature.VerifySig(msg.Payload, msg.Sig, newPublicKey)
	if !verify {
		panic(fmt.Errorf("===>[OutputERROR]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	result := new(string)
	if err := json.Unmarshal(msg.Payload, result); err != nil {
		panic(fmt.Errorf("======>[Output]Invalid[%s] Output message[%s]", err, msg))
	}

	// new output message
	// send output message to primary node
	if *result == s.Result.String() {
		s.OutputNum++
		// if s.OutputNum == util.TotalNodeNum-1 {
		if s.OutputNum == 2*util.MaxFaultyNode {
			fmt.Println("[Output]new output is", s.Result)
			s.GlobalTimer.tack()
			util.WriteResult(s.Result.String())
			s.stage = Collect
			time.Sleep(5 * time.Second)
			config.WriteOutput(s.Result.String())
		}
	}

	return
}

// handle different kinds of consensus messages
func (s *StateEngine) procConsensusMsg(msg *message.ConMessage) (err error) {
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("\n======>[procConsensusMsg]Consesus message type:[%s] from Node[%d]\n", msg.Typ, msg.From)
	fmt.Println(s.stage)

	switch msg.Typ {
	case message.MTSubmit:
		if s.stage == Submit {
			// return s.unionTC(msg)
			go s.unionTC(msg)
		}
	case message.MTApprove:
		if s.stage == Approve {
			// return s.approveTC(msg)
			go s.approveTC(msg)
		}
	case message.MTConfirm:
		if s.stage == Confirm {
			// return s.confirmTC(msg)
			go s.confirmTC(msg)
		}
	case message.MTOutput:
		if s.stage == Output {
			// return s.outputTC(msg)
			go s.outputTC(msg)
		}

	}

	return
}

// handle different kinds of manage messages
func (s *StateEngine) procManageMsg(msg *message.ConMessage) (err error) {
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("\n======>[procConsensusMsg]Manage message type:[%s] from Node[%d]\n", msg.Typ, msg.From)
	fmt.Println(s.stage)

	switch msg.Typ {
	case message.MTViewChange:
		// if s.stage == ViewChange {
		go s.procViewChange(msg)
		// }
		// vc := &message.ViewChange{}
		// if err := json.Unmarshal(msg.Payload, vc); err != nil {
		// 	return fmt.Errorf("======>[procConsensusMsg] invalid[%s]ViewChange message[%s]", err, msg)
		// }
		// return s.procViewChange(vc, msg)

	case message.MTNewView:
		go s.didChangeView(msg)
		// vc := &message.NewView{}
		// if err := json.Unmarshal(msg.Payload, vc); err != nil {
		// 	return fmt.Errorf("======>[procConsensusMsg] invalid[%s] didiViewChange message[%s]", err, msg)
		// }
		// return s.didChangeView(vc, msg)
	}

	return nil
}
