package service

import (
	"consensusNode/config"
	"consensusNode/consensus"
	"consensusNode/message"
	"consensusNode/signature"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// [SrvHub]: contians UDP connection with client
// [nodeChan]: a channel connects [service] with [node], deliver service message(request), corresponding to [srvChan] in [node]
type Service struct {
	SrvHub   *net.UDPConn
	nodeChan chan interface{}
}

// initialize node's service
func InitService(port int, msgChan chan interface{}) *Service {
	locAddr := net.UDPAddr{
		Port: port,
	}
	srv, err := net.ListenUDP("udp4", &locAddr)
	if err != nil {
		return nil
	}
	fmt.Printf("===>Service is Listening at[%d]\n", port)

	s := &Service{
		SrvHub:   srv,
		nodeChan: msgChan,
	}

	return s
}

// wait for request from client in UDP channel
func (s *Service) WaitRequest(sig chan interface{}, stateMachine *consensus.StateEngine) {
	defer func() {
		if r := recover(); r != nil {
			sig <- r
		}
	}()

	buf := make([]byte, 2048)
	for {
		n, rAddr, err := s.SrvHub.ReadFromUDP(buf)
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
		verify = signature.VerifySig(previousOutput, entropyMsg.VRFResult, entropyPK)
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

		// NewTimeCommitment(tc string)
		fmt.Printf("===>Entropy message from[%s], Entropy Node id is[%d],time commitment is:%s\n", rAddr.String(), entropyMsg.ClientID, entropyMsg.TimeCommitment)

		// process the request message
		// go s.process(requestInPBFT)
	}
}

// process the request message
// send the message by s.nodeChan(node.srvChan), then it will invokeconsensus.inspireConsensus()
func (s *Service) process(op *message.Request) {
	/*
		TODO:: Check operation
		1. if clientID is authorized
		2. if operation is valid
	*/
	s.nodeChan <- op
}

// After commit phase, execute the operation in the request message
func (s *Service) Execute(v, n, seq int64, o *message.Request) (reply *message.Reply, err error) {
	fmt.Printf("===>Service is executing opertion[%s]......\n", o.Operation)
	r := &message.Reply{
		SeqID:     seq,
		ViewID:    v,
		Timestamp: o.TimeStamp,
		ClientID:  o.ClientID,
		NodeID:    n,
		Result:    "success",
	}

	bs, _ := json.Marshal(r)
	cAddr := net.UDPAddr{
		Port: 8088,
	}
	no, err := s.SrvHub.WriteToUDP(bs, &cAddr)
	if err != nil {
		fmt.Printf("===>[ERROR]Reply client failed:%s\n", err)
		return nil, err
	}
	fmt.Printf("===>Reply Success! Seq is [%d], Result is [%s], Length is [%d]\n", seq, r.Result, no)

	return r, nil
}

// the same request asked twice, directly reply the result
func (s *Service) DirectReply(r *message.Reply) error {
	bs, _ := json.Marshal(r)
	cAddr := net.UDPAddr{
		Port: 8088,
	}
	no, err := s.SrvHub.WriteToUDP(bs, &cAddr)
	if err != nil {
		fmt.Printf("===>[ERROR]Reply client failed:%s\n", err)
		return err
	}
	fmt.Printf("===>Reply Success! Seq is [%d], Result is [%s], Length is [%d]\n", r.SeqID, r.Result, no)

	return nil
}
