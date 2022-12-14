package consensus

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/p2pnetwork"
	"consensusNode/util"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/algorand/go-algorand/crypto"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const MaxStateMsgNO = 150

type StateEngine struct {
	NodeID     int64 `json:"nodeID"`
	CurViewID  int64 `json:"viewID"`
	PrimaryID  int64 `json:"primaryID"`
	SubmitNum  int64
	ConfirmNum int64
	OutputNum  int64
	ApproveNum int
	Result     *big.Int
	SrvHub     *net.TCPListener

	sCache     *VCCache
	stage      Stage
	Mutex      sync.Mutex
	P2pWire    p2pnetwork.P2pNetwork
	MsgChan    <-chan *message.ConMessage
	nodeStatus EngineStatus
	quit       chan bool

	GlobalTimer  *RequestTimer
	CollectTimer *RequestTimer
	SubmitTimer  *RequestTimer
	VerifyTime   float64

	entropyNode           map[int64]bool
	TimeCommitment        map[string][4]string
	TimeCommitmentSubmit  map[int64]int
	TimeCommitmentApprove map[string]bool
	TimeCommitmentProof   map[string][4]string
}

func InitConsensus(id int64) *StateEngine {
	fmt.Printf("\n===>[InitConsensus]Service is Listening at[%d]\n", util.PortByID(id))

	ch := make(chan *message.ConMessage, MaxStateMsgNO)
	p2p := p2pnetwork.NewSimpleP2pLib(id, ch)
	se := &StateEngine{
		NodeID:     id,
		CurViewID:  0,
		PrimaryID:  0,
		SubmitNum:  0,
		OutputNum:  0,
		ConfirmNum: 0,
		ApproveNum: 0,
		Result:     big.NewInt(0),

		sCache:  NewVCCache(),
		stage:   Collect,
		Mutex:   sync.Mutex{},
		P2pWire: p2p,
		MsgChan: ch,
		quit:    make(chan bool),

		GlobalTimer:  newRequestTimer(),
		CollectTimer: newRequestTimer(),
		SubmitTimer:  newRequestTimer(),
		VerifyTime:   0.0,

		entropyNode:           make(map[int64]bool),
		TimeCommitment:        make(map[string][4]string),
		TimeCommitmentSubmit:  make(map[int64]int),
		TimeCommitmentApprove: make(map[string]bool),
		TimeCommitmentProof:   make(map[string][4]string),
	}
	se.PrimaryID = se.CurViewID % util.TotalNodeNum

	if se.PrimaryID == se.NodeID {
		go se.WriteRandomOutput()
	}

	return se
}

// To start randomness beacon, primary writes a random output into output.yml
func (s *StateEngine) WriteRandomOutput() {
	time.Sleep(60 * time.Second)
	fmt.Println("\n===>[WriteRandomOutput]start wirte config")

	// generate random init input
	message := []byte("hello world")
	randomNum := util.Digest(message)

	config.WriteOutput(string(randomNum))
}

// receive and handle consensus message
func (s *StateEngine) StartConsensus(sig chan interface{}) {
	s.nodeStatus = Serving

	for {
		select {
		// Global timer out, starts viewchange
		case <-s.GlobalTimer.C:
			s.GlobalTimer.tack()
			fmt.Printf("\n===[Node%d]Time is out and view change starts\n", s.NodeID)

			s.nodeStatus = ViewChanging
			s.stage = ViewChange
			s.ViewChange()

		// Collect timer out
		// backups stop receiving VRF and TC messages
		// backups start sending union messages to the primary node
		case <-s.CollectTimer.C:
			s.CollectTimer.tack()
			s.stage = Submit
			if s.NodeID != s.PrimaryID {
				// if s.PrimaryID == 0 && s.NodeID%2 == 1 {
				// 	fmt.Println("Rotten")
				// } else {
				// 	go s.sendUnionMsg()
				// }
				go s.sendUnionMsg()
			}

			fmt.Printf("\n===>[Node%d]Stop Receive messages\n", s.NodeID)
			for key, value := range s.TimeCommitment {
				fmt.Println("===>[Collect]key is", key)
				fmt.Println("===>[Collect]value is", value)
			}

		// handle consensus message received from other consensus nodes
		case conMsg := <-s.MsgChan:
			switch conMsg.Typ {
			case message.MTSubmit,
				message.MTApprove,
				message.MTConfirm,
				message.MTOutput:
				if s.nodeStatus != Serving {
					fmt.Println("===>[ERROR from StartConsensus]node is not in service status now......")
					continue
				}
				if err := s.procConsensusMsg(conMsg); err != nil {
					panic(fmt.Errorf("===>[ERROR from StartConsensus]procConsensusMsg failed:%s", err))
				}
			case message.MTViewChange,
				message.MTNewView:
				if err := s.procManageMsg(conMsg); err != nil {
					panic(fmt.Errorf("===>[ERROR from StartConsensus]procManageMsg failed:%s", err))
				}
			}
		}
	}
}

// consensus node continues watching on output.yml
// when [previousOutput] in output.yml changes, consensus node starts a new round
func (s *StateEngine) WatchConfig(id int64, sig chan interface{}) {
	// get the earliest output written in output.yml
	previousOutput := string(config.GetPreviousOutput())
	fmt.Println("\n===>[Watching]The earliest output is", previousOutput)

	// watch config file
	myViper := viper.New()
	myViper.SetConfigFile("../Configuration/output.yml")
	myViper.WatchConfig()
	myViper.OnConfigChange(func(e fsnotify.Event) {
		// lock file
		f, err := os.Open("../lock")
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from WatchConfig]Open lock failed:%s", err))
		}
		// share lock, concurrently read
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH); err != nil {
			panic(fmt.Errorf("===>[ERROR from WatchConfig]Add share lock failed:%s", err))
		}
		fmt.Println("===>[Watching]Configuration Changed", time.Now())

		// when new output comes, entropy node starts calculating VRF and sending TC
		newOutput := string(config.GetPreviousOutput())
		if previousOutput != newOutput && newOutput != "" && s.stage == Collect {
			startTime = time.Now()
			fmt.Println("\n===>[Watching]Start time is", startTime)
			fmt.Println("\n===>[Watching]Output changed,new output is", newOutput)
			previousOutput = newOutput

			// initialize a series of variales
			s.stage = Collect
			s.OutputNum = 0
			s.ConfirmNum = 0
			s.SubmitNum = 0
			s.ApproveNum = 0
			s.VerifyTime = 0.0
			s.quit = make(chan bool)
			s.entropyNode = make(map[int64]bool)
			s.TimeCommitment = make(map[string][4]string)
			s.TimeCommitmentSubmit = make(map[int64]int)
			s.TimeCommitmentApprove = make(map[string]bool)
			s.TimeCommitmentProof = make(map[string][4]string)

			// listen the srvHub and set a deadline
			if s.SrvHub == nil {
				locAddr := net.TCPAddr{
					Port: util.EntropyPortByID(id),
				}
				srvHub, err := net.ListenTCP("tcp4", &locAddr)
				if err != nil {
					panic(fmt.Errorf("===>[ERROR from WatchConfig]Listen TCP port failed:%s", err))
				}
				s.SrvHub = srvHub
			}
			s.SrvHub.SetDeadline(time.Now().Add(10 * time.Second))

			// wait for TC messages from entropy nodes
			go s.WaitTC(sig, s.quit)

			// start 3 timers for a new round
			s.GlobalTimer.tick(10 * time.Minute)
			s.CollectTimer.tick(10 * time.Second)
			// if s.NodeID == s.PrimaryID {
			// 	s.SubmitTimer.tick(1 * time.Second)
			// }
		}

		// unlock file
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
			panic(fmt.Errorf("===>[ERROR from WatchConfig]Unlock share lock failed:%s", err))
		}
	})
}

// wait for TC messages from entropy nodes
func (s *StateEngine) WaitTC(sig chan interface{}, quit chan bool) {
	defer func() {
		if r := recover(); r != nil {
			sig <- r
		}
	}()

	buf := make([]byte, 8192)
	for {
		select {
		case <-quit:
			fmt.Println("Quit WaitTC")
			return

		default:
			conn, err := s.SrvHub.AcceptTCP()
			if err != nil {
				// fmt.Printf("===>[ERROR]TCP:%s\n", err)
				continue
			}
			n, err := conn.Read(buf)
			if err != nil {
				panic(fmt.Errorf("===>[ERROR from WaitTC]Service received failed:%s", err))
			}

			// get message from entropy node
			msgFromEntropyNode := &message.ConMessage{}
			if err := json.Unmarshal(buf[:n], msgFromEntropyNode); err != nil {
				fmt.Printf("===>[ERROR from WaitTC]Message parse failed:%s", err)
				continue
			}
			if msgFromEntropyNode.Typ != message.MTCollect && msgFromEntropyNode.Typ != message.MTVRFVerify {
				fmt.Printf("===>[ERROR from WaitTC]Not vrf Verify message or collect message:\n")
				continue
			}

			// handle vrf verify message
			if msgFromEntropyNode.Typ == message.MTVRFVerify {
				fmt.Printf("===>[WaitTC]new vrf verify message from Node[%d]\n", msgFromEntropyNode.From)

				entropyVRFMsg := &message.EntropyVRFMessage{}
				if err := json.Unmarshal(msgFromEntropyNode.Payload, entropyVRFMsg); err != nil {
					fmt.Printf("===>[ERROR from WaitTC]Invalid[%s] Entropy VRF message[%s]", err, msgFromEntropyNode)
					continue
				}

				// verify the VRF result
				msg := MessageHashable{
					Data: entropyVRFMsg.Msg,
				}
				var pk crypto.VrfPubkey = entropyVRFMsg.PublicKey
				verify, _ := pk.Verify(entropyVRFMsg.VRFResult, msg)
				if !verify {
					fmt.Printf("===>[ERROR from WaitTC]Verify new VRF result failed, From Entropy Node[%d]\n", entropyVRFMsg.ClientID)
					continue
				}
				// fmt.Println("===>[WaitTC]VRF Output is", output)
				fmt.Println("===>[WaitTC]VRF Verify success!!!")

				// verify the VRF and difficuly
				difficulty := config.GetDifficulty()
				vrfResultBinary := util.BytesToBinaryString(entropyVRFMsg.VRFResult)
				vrfResultTail, err := strconv.Atoi(vrfResultBinary[len(vrfResultBinary)-1:])
				if err != nil {
					fmt.Printf("===>[ERROR from WaitTC]Failed to get VRF result's last bit, From Entropy Node[%d]\n", entropyVRFMsg.ClientID)
					continue
				}
				if vrfResultTail != difficulty {
					fmt.Println("===>[WaitTC]Not pass but who cares!!!!")
					// fmt.Println("===>[WaitTC]Cheater!!!!")
					// fmt.Printf("===>[ERROR from WaitTC]Verify difficulty failed, From Entropy Node[%d]\n", entropyVRFMsg.ClientID)
					// continue
				}

				// register a new entropy node
				s.Mutex.Lock()
				s.entropyNode[entropyVRFMsg.ClientID] = false
				s.Mutex.Unlock()
			}

			// handle collect message
			if msgFromEntropyNode.Typ == message.MTCollect {
				fmt.Printf("===>[WaitTC]new collect message from Node[%d]\n", msgFromEntropyNode.From)

				entropyTCMsg := &message.EntropyTCMessage{}
				if err := json.Unmarshal(msgFromEntropyNode.Payload, entropyTCMsg); err != nil {
					fmt.Printf("===>[ERROR from WaitTC]Invalid[%s] Entropy TC message[%s]", err, msgFromEntropyNode)
					continue
				}

				timedCommitment := [4]string{entropyTCMsg.TimeCommitmentC, entropyTCMsg.TimeCommitmentH, entropyTCMsg.TimeCommitmentrKSubOne, entropyTCMsg.TimeCommitmentrK}
				timedCommitmentProof := [4]string{entropyTCMsg.TimeCommitmentA1, entropyTCMsg.TimeCommitmentA2, entropyTCMsg.TimeCommitmentA3, entropyTCMsg.TimeCommitmentZ}

				s.Mutex.Lock()
				// add new TC element
				_, ok := s.entropyNode[entropyTCMsg.ClientID]
				if !ok {
					fmt.Printf("===>[ERROR from WaitTC]Not registered")
					continue
				}

				key := string(util.Digest(timedCommitment))
				s.entropyNode[entropyTCMsg.ClientID] = true
				s.TimeCommitment[key] = timedCommitment
				s.TimeCommitmentProof[key] = timedCommitmentProof
				s.TimeCommitmentApprove[key] = false
				s.Mutex.Unlock()

				fmt.Printf("===>[WaitTC]Legal Entropy message from Node[%d]\n", entropyTCMsg.ClientID)
			}
		}
	}
}

// handle different kinds of consensus messages
func (s *StateEngine) procConsensusMsg(msg *message.ConMessage) (err error) {
	// time.Sleep(100 * time.Millisecond)

	fmt.Printf("\n===>[procConsensusMsg]Consesus message type:[%s] from Node[%d]\n", msg.Typ, msg.From)
	fmt.Println("===>[procConsensusMsg]Stage is", s.stage)

	switch msg.Typ {
	case message.MTSubmit:
		if s.stage == Submit {
			// go s.unionTC(msg)
			s.unionTC(msg)
		}
	case message.MTApprove:
		if s.stage == Approve {
			// go s.approveTC(msg)
			s.approveTC(msg)
		}
	case message.MTConfirm:
		if s.stage == Confirm {
			// go s.confirmTC(msg)
			s.confirmTC(msg)
		}
	case message.MTOutput:
		if s.stage == Output {
			// go s.outputTC(msg)
			s.outputTC(msg)
		}

	}

	return
}

// handle different kinds of manage messages
func (s *StateEngine) procManageMsg(msg *message.ConMessage) (err error) {
	// time.Sleep(100 * time.Millisecond)

	fmt.Printf("\n===>[procConsensusMsg]Manage message type:[%s] from Node[%d]\n", msg.Typ, msg.From)
	fmt.Println("===>[procConsensusMsg]Stage is", s.stage)

	switch msg.Typ {
	case message.MTViewChange:
		// go s.procViewChange(msg)
		s.procViewChange(msg)

	case message.MTNewView:
		// go s.didChangeView(msg)
		s.didChangeView(msg)
	}

	return nil
}
