package consensus

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/signature"
	tc "consensusNode/timedCommitment"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"time"
)

// backups send union message
func (s *StateEngine) sendUnionMsg() {
	time.Sleep(200 * time.Millisecond)
	// new submit message
	// send submit message to primary node
	for key, value := range s.TimeCommitment {
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
		time.Sleep(800 * time.Millisecond)
	}

	s.stage = Approve
	fmt.Printf("\n===>[Submit]Send Submit(Union) message success\n")
}

// primary union received TC
func (s *StateEngine) unionTC(msg *message.ConMessage) (err error) {
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

	// verify union TC
	key := string(util.Digest(submit.CollectTC))
	verifyResult := tc.VerifyTC(submit.TCProof[0], submit.TCProof[1], submit.TCProof[2], submit.TCProof[3], submit.CollectTC[1], submit.CollectTC[2], submit.CollectTC[3])
	if verifyResult {
		fmt.Println("===>[Union]pass all tests!")
	} else {
		fmt.Println("===>[Union]Failed to pass all tests!")
		return
	}

	// Add new TC
	if value, ok := s.TimeCommitment[key]; !ok {
		fmt.Println("===>[Union]new key is", key)
		fmt.Println("===>[Union]new TimeCommitmentC value is", value[0])
		s.TimeCommitment[key] = submit.CollectTC
		s.TimeCommitmentApprove[key] = false
	}

	// Count SubmitNum
	if _, ok := s.TimeCommitmentSubmit[msg.From]; !ok {
		fmt.Printf("===>[Union]node[%d]Starts sending submit messages\n", msg.From)
		s.Mutex.Lock()
		s.TimeCommitmentSubmit[msg.From] = 0
		s.Mutex.Unlock()
	}
	s.TimeCommitmentSubmit[msg.From]++
	if s.TimeCommitmentSubmit[msg.From] == submit.Length {
		s.SubmitNum++
	}

	return
}
