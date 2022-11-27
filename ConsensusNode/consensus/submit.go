package consensus

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/signature"
	tc "consensusNode/timedCommitment"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime"
	"time"
)

// backups send union message
func (s *StateEngine) sendUnionMsg() {

	// new submit message
	// send submit message to primary node
	for key, value := range s.TimeCommitment {
		// rand2.Seed(time.Now().UnixNano())
		// delay := 400 + rand2.Intn(100)
		// time.Sleep(time.Duration(delay) * time.Millisecond)
		// time.Sleep(2 * time.Second)
		time.Sleep(500 * time.Millisecond)
		tc := value
		tcProof := s.TimeCommitmentProof[key]
		submit := &message.Submit{
			Length:    len(s.TimeCommitment),
			CollectTC: tc,
			TCProof:   tcProof,
		}

		sk := s.P2pWire.GetMySecretkey()
		sMsg := message.CreateConMsg(message.MTSubmit, submit, sk, s.NodeID)
		conn := s.P2pWire.GetPrimaryConn(s.PrimaryID)
		if err := s.P2pWire.SendUniqueNode(conn, sMsg); err != nil {
			panic(fmt.Errorf("===>[ERROR from sendUnionMsg]send message error:%s", err))
		}
		// time.Sleep(500 * time.Millisecond)
		// time.Sleep(1500 * time.Millisecond)
	}

	s.stage = Approve
	fmt.Printf("\n===>[Submit]Send Submit(Union) message success\n")
}

// primary union received TC
func (s *StateEngine) unionTC(msg *message.ConMessage) (err error) {
	start := time.Now()
	fmt.Printf("\n===>[Union]Current Submit(Union) Message from Node[%d]\n", msg.From)

	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from unionTC]Parse elliptic curve error:%s", err))
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
		panic(fmt.Errorf("===>[ERROR from unionTC]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	submit := &message.Submit{}
	if err := json.Unmarshal(msg.Payload, submit); err != nil {
		panic(fmt.Errorf("===>[ERROR from unionTC]Invalid[%s] Union message[%s]", err, msg))
	}
	fmt.Printf("===>[Union]Submit(Union) Message from Node[%d],length is %d\n", msg.From, submit.Length)

	// Add new TC
	key := string(util.Digest(submit.CollectTC))
	if value, ok := s.TimeCommitment[key]; !ok {
		fmt.Println("===>[Union]new key is", key)
		fmt.Println("===>[Union]new TimeCommitmentC value is", value[0])
		s.TimeCommitment[key] = submit.CollectTC
		s.TimeCommitmentApprove[key] = false
	}

	// Count SubmitNum
	// s.Mutex.Lock()
	if _, ok := s.TimeCommitmentSubmit[msg.From]; !ok {
		fmt.Printf("===>[Union]node[%d]Starts sending submit messages\n", msg.From)
		s.TimeCommitmentSubmit[msg.From] = 0
	}
	s.TimeCommitmentSubmit[msg.From]++
	if s.TimeCommitmentSubmit[msg.From] == submit.Length {
		s.SubmitNum++
	}

	fmt.Println(s.SubmitNum)
	if s.SubmitNum == 2*util.MaxFaultyNode+1 {
		currentTime1 := time.Now()
		fmt.Println("===>[Union]SubmitNum is", s.SubmitNum)
		fmt.Println("===>[Union]From start to submit finished,passed time1 is", currentTime1.Sub(startTime).Seconds())
		// s.Mutex.Unlock()
		for key := range s.TimeCommitment {
			fmt.Println(len(s.TimeCommitment))
			fmt.Println(len(s.TimeCommitmentProof))
			verifyResult, verifyTime := tc.VerifyTC(s.TimeCommitmentProof[key][0], s.TimeCommitmentProof[key][1], s.TimeCommitmentProof[key][2], s.TimeCommitmentProof[key][3], s.TimeCommitment[key][1], s.TimeCommitment[key][2], s.TimeCommitment[key][3])
			if verifyResult {
				s.VerifyTime += verifyTime
				fmt.Println("===>[Union]pass all tests!")
				fmt.Println("===>[Union]Current verify total time is", s.VerifyTime)
			} else {
				fmt.Println("===>[Union]Failed to pass all tests!")
				delete(s.TimeCommitment, key)
				delete(s.TimeCommitmentProof, key)
				delete(s.TimeCommitmentApprove, key)
			}
		}

		currentTime := time.Now()
		fmt.Println("===>[Union]From start to submit finished,passed time is", currentTime.Sub(startTime).Seconds())

		// new approve message
		// send approve message to backup nodes
		approve := &message.Approve{}
		result, _ := rand.Int(rand.Reader, big.NewInt(10000))
		if result.Cmp(big.NewInt(0)) == 0 {
			// send wrong approve message
			fmt.Println("===>[Union]I'm crazy!!!!!")
			sk := s.P2pWire.GetMySecretkey()
			aMsg := message.CreateConMsg(message.MTApprove, approve, sk, s.NodeID)
			if err := s.P2pWire.BroadCast(aMsg); err != nil {
				panic(fmt.Errorf("===>[ERROR from StartConsensus]Broadcast failed:%s", err))
			}
		} else {
			// time.Sleep(1 * time.Second)

			// send right approve message
			var tc [4]string
			for _, value := range s.TimeCommitment {
				// rand2.Seed(time.Now().UnixNano())
				// delay := 400 + rand2.Intn(100)
				// time.Sleep(time.Duration(delay) * time.Millisecond)
				time.Sleep(500 * time.Millisecond)
				// time.Sleep(1 * time.Second)
				tc = value
				approve = &message.Approve{
					Length:  len(s.TimeCommitment),
					UnionTC: tc,
				}
				sk := s.P2pWire.GetMySecretkey()
				aMsg := message.CreateConMsg(message.MTApprove, approve, sk, s.NodeID)
				if err := s.P2pWire.BroadCast(aMsg); err != nil {
					panic(fmt.Errorf("===>[ERROR from StartConsensus]Broadcast failed:%s", err))
				}
				// time.Sleep(150 * time.Millisecond)
				// time.Sleep(300 * time.Millisecond)
			}
		}

		s.stage = Confirm
		s.ConfirmNum++
		fmt.Printf("===>[Union]Send approve message success\n")
	}
	// else {
	// 	s.Mutex.Unlock()
	// }

	end := time.Now()
	fmt.Println("===>[Union]passed time", end.Sub(start).Seconds())
	fmt.Println("===>[Union]the number of goroutines: ", runtime.NumGoroutine())

	return
}
