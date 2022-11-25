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
func (s *StateEngine) approveTC(msg *message.ConMessage) (err error) {
	fmt.Printf("\n===>[Approve]Current Approve Message from Node[%d]\n", msg.From)
	fmt.Println("===>[Approve]my TC length is", len(s.TimeCommitment))
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from approveTC]Parse elliptic curve error:%s", err))
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
		panic(fmt.Errorf("===>[ERROR from approveTC]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	approve := &message.Approve{}
	if err := json.Unmarshal(msg.Payload, approve); err != nil {
		panic(fmt.Errorf("===>[ERROR from approveTC]Invalid[%s] Approve message[%s]", err, msg))
	}
	fmt.Printf("===>[Approve]Approve Message from Node[%d],length is %d\n", msg.From, approve.Length)

	// check union tc set
	if approve.Length == 0 && len(s.TimeCommitment) != 0 {
		return fmt.Errorf("===>[ERROR from approveTC]are you sure????")
	}

	key := string(util.Digest(approve.UnionTC))
	if _, ok := s.TimeCommitment[key]; !ok {
		fmt.Println("===>[Approve]new key is", key)
		s.TimeCommitment[key] = approve.UnionTC
	}
	s.TimeCommitmentApprove[key] = true
	s.ApproveNum++
	fmt.Printf("===>[Approve]Current length is %d\n", s.ApproveNum)

	if s.ApproveNum == approve.Length {
		// check whether local TC is a subset of Approve.UnionTC
		for _, flag := range s.TimeCommitmentApprove {
			if !flag {
				s.stage = Error
				return fmt.Errorf("===>[ERROR from approveTC]where is my tc????")
			}
		}
		s.ConfirmNum++
		s.stage = Confirm
		s.sendConfirmMsg()
	}

	return
}

// backups broadcast confirm message
func (s *StateEngine) sendConfirmMsg() {
	time.Sleep(500 * time.Millisecond)

	confirm := &message.Confirm{
		Length: len(s.TimeCommitment),
	}

	sk := s.P2pWire.GetMySecretkey()
	cMsg := message.CreateConMsg(message.MTConfirm, confirm, sk, s.NodeID)

	if err := s.P2pWire.BroadCast(cMsg); err != nil {
		panic(fmt.Errorf("===>[ERROR from sendConfirmMsg]send message error:%s", err))
	}

	fmt.Printf("===>[Approve]Send confirm message success\n")
}
