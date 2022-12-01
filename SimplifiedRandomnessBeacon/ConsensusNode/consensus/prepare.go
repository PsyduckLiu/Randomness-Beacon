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
	"math/big"
	"time"
)

// Prepare TC
func (s *StateEngine) PrepareTC(msg *message.ConMessage) (err error) {
	if s.PrepareNum >= 2*util.MaxFaultyNode+1 {
		fmt.Println("late!")
		return
	}

	fmt.Printf("\n===>[Prepare]Current Prepare Message from Node[%d]\n", msg.From)
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from PrepareTC]Parse elliptic curve error:%s", err))
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
		panic(fmt.Errorf("===>[ERROR from PrepareTC]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	Prepare := &message.Prepare{}
	if err := json.Unmarshal(msg.Payload, Prepare); err != nil {
		panic(fmt.Errorf("===>[ERROR from PrepareTC]Invalid[%s] Prepare message[%s]", err, msg))
	}

	if Prepare.Length != len(s.TimeCommitment) {
		return fmt.Errorf("===>[ERROR from PrepareTC]Invalid tc set")
	}

	// ToDo: verify each TC in cofirmTC
	s.PrepareNum++
	if s.PrepareNum == 2*util.MaxFaultyNode+1 {
		fmt.Println("===>[Prepare]Prepare success")
		s.stage = Reveal
		time.Sleep(1 * time.Second)
		// go s.handleTC()
		s.handleTC()
	}

	return
}

// resolve tc and compute F
func (s *StateEngine) handleTC() (err error) {
	// resolve TC set
	var resolvedTC []*big.Int
	for _, value := range s.TimeCommitment {
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
	fmt.Println("\n===>[handleTC]after xor all resolved tc,result is:", s.Result)

	if s.NodeID != s.PrimaryID {
		sk := s.P2pWire.GetMySecretkey()
		oMsg := message.CreateConMsg(message.MTReveal, s.Result.String(), sk, s.NodeID)
		conn := s.P2pWire.GetPrimaryConn(s.PrimaryID)

		// rand2.Seed(time.Now().UnixNano())
		// delay := 400 + rand2.Intn(100)
		// time.Sleep(time.Duration(delay) * time.Millisecond)
		time.Sleep(500 * time.Millisecond)
		// time.Sleep(1 * time.Second)
		if err := s.P2pWire.SendUniqueNode(conn, oMsg); err != nil {
			panic(fmt.Errorf("===>[ERROR from handleTC]send message error:%s", err))
		}
		s.stage = CommitFromEntropy
		s.GlobalTimer.tack()
	}

	return
}
