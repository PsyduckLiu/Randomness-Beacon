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

// output tc and compute F
func (s *StateEngine) outputTC(msg *message.ConMessage) (err error) {
	if s.OutputNum >= 2*util.MaxFaultyNode {
		fmt.Println("late!")
		return
	}

	fmt.Printf("\n===>[Output]Current Output Message from Node[%d]\n", msg.From)
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from outputTC]Parse elliptic curve error:%s", err))
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
		panic(fmt.Errorf("===>[ERROR from outputTC]Verify new public key Signature failed, From Node[%d]", msg.From))
	}

	// unmarshal message
	result := new(string)
	if err := json.Unmarshal(msg.Payload, result); err != nil {
		panic(fmt.Errorf("===>[ERROR from outputTC]Invalid[%s] Output message[%s]", err, msg))
	}

	// new output message
	// send output message to primary node
	if *result == s.Result.String() {
		s.OutputNum++
		// if s.OutputNum == util.TotalNodeNum-1 {
		if s.OutputNum == 2*util.MaxFaultyNode {
			fmt.Println("===>[Output]new output is", s.Result)
			s.GlobalTimer.tack()
			util.WriteResult(s.Result.String())
			s.stage = Collect
			time.Sleep(5 * time.Second)
			config.WriteOutput(s.Result.String())
		}
	}

	return
}
