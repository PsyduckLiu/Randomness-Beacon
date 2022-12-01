package consensus

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/signature"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"time"
)

// backups check union tc sent by primary
func (s *StateEngine) ProposeTC(msg *message.ConMessage) (err error) {
	fmt.Printf("\n===>[Propose]Current Propose Message from Node[%d]\n", msg.From)
	fmt.Println("===>[Propose]my TC length is", len(s.TimeCommitment))
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from ProposeTC]Parse elliptic curve error:%s", err))
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
		panic(fmt.Errorf("===>[ERROR from ProposeTC]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	Propose := &message.Propose{}
	if err := json.Unmarshal(msg.Payload, Propose); err != nil {
		panic(fmt.Errorf("===>[ERROR from ProposeTC]Invalid[%s] Propose message[%s]", err, msg))
	}
	fmt.Printf("===>[Propose]Propose Message from Node[%d],length is %d\n", msg.From, Propose.Length)

	// check union tc set
	if Propose.Length == 0 && len(s.TimeCommitment) != 0 {
		return fmt.Errorf("===>[ERROR from ProposeTC]are you sure????")
	}

	key := string(util.Digest(Propose.UnionTC))
	if _, ok := s.TimeCommitment[key]; !ok {
		fmt.Println("===>[Propose]new key is", key)
		s.TimeCommitment[key] = Propose.UnionTC
	}
	s.TimeCommitmentPropose[key] = true
	s.ProposeNum++
	fmt.Printf("===>[Propose]Current length is %d\n", s.ProposeNum)

	if s.ProposeNum == Propose.Length {
		// check whether local TC is a subset of Propose.UnionTC
		for _, flag := range s.TimeCommitmentPropose {
			if !flag {
				s.stage = Error
				return fmt.Errorf("===>[ERROR from ProposeTC]where is my tc????")
			}
		}
		s.PrepareNum++
		s.stage = Prepare
		time.Sleep(1 * time.Second)
		s.sendPrepareMsg()
	}

	return
}

// backups broadcast Prepare message
func (s *StateEngine) sendPrepareMsg() {
	// rand2.Seed(time.Now().UnixNano())
	// delay := 400 + rand2.Intn(100)
	// time.Sleep(time.Duration(delay) * time.Millisecond)
	time.Sleep(500 * time.Millisecond)
	// time.Sleep(1 * time.Second)

	Prepare := &message.Prepare{
		Length: len(s.TimeCommitment),
	}

	sk := s.P2pWire.GetMySecretkey()
	cMsg := message.CreateConMsg(message.MTPrepare, Prepare, sk, s.NodeID)

	if err := s.P2pWire.BroadCast(cMsg); err != nil {
		panic(fmt.Errorf("===>[ERROR from sendPrepareMsg]send message error:%s", err))
	}

	fmt.Printf("===>[Propose]Send Prepare message success\n")
}
