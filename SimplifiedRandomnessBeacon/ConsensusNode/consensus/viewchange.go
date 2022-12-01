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

// view change cache
// vcMsg = view change message
// VMessage map[int64]*ViewChange
// nvMsg = new view message
type VCCache struct {
	vcMsg message.VMessage
	nvMsg map[int64]*message.NewView
}

// initialize a new VCCache
func NewVCCache() *VCCache {
	return &VCCache{
		vcMsg: make(message.VMessage),
		nvMsg: make(map[int64]*message.NewView),
	}
}

// push new vc in VCCache
func (vcc *VCCache) pushVC(vc *message.ViewChange) {
	vcc.vcMsg[vc.NodeID] = vc
}

// check whether hasNewViewYet
func (vcc *VCCache) hasNewViewYet(vid int64) bool {
	if _, ok := vcc.nvMsg[vid]; ok {
		return true
	}
	return false
}

// add newview into VCCache
func (vcc *VCCache) addNewView(nv *message.NewView) {
	vcc.nvMsg[nv.NewViewID] = nv
}

// invoked by state.go when timeout
func (s *StateEngine) ViewChange() {
	fmt.Printf("\n===>[ViewChange]Current view is(%d).....\n", s.CurViewID)

	// reset configuration
	s.CollectTimer.tack()
	s.CurViewID++
	s.CollectNum = 0
	s.sCache.nvMsg = make(map[int64]*message.NewView)

	// handle viewChange message
	vc := &message.ViewChange{
		NewViewID: s.CurViewID,
		NodeID:    s.NodeID,
	}
	nextPrimaryID := vc.NewViewID % util.TotalNodeNum
	for key, value := range s.sCache.vcMsg {
		if value.NewViewID%util.TotalNodeNum != nextPrimaryID {
			delete(s.sCache.vcMsg, key)
		}
	}
	fmt.Println("===>[ViewChange]Current viewchange length is", len(s.sCache.vcMsg))
	if s.NodeID == nextPrimaryID {
		s.sCache.pushVC(vc)
	}

	// broadcast viewChange message
	sk := s.P2pWire.GetMySecretkey()
	consMsg := message.CreateConMsg(message.MTViewChange, vc, sk, s.NodeID)
	if err := s.P2pWire.BroadCast(consMsg); err != nil {
		fmt.Println(err)
		return
	}
}

// invoked by state.go when received a viewchage message
func (s *StateEngine) procViewChange(msg *message.ConMessage) error {
	fmt.Printf("===>[ViewChange]Current ViewChange Message from Node[%d]\n", msg.From)
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from ViewChange]Parse elliptic curve error:%s", err))
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
		panic(fmt.Errorf("===>[ERROR from ViewChange]Verify new public key Signature failed, From Node[%d]", msg.From))
	}
	fmt.Printf("===>[ViewChange]Verify success\n")

	// unmarshal message
	vc := &message.ViewChange{}
	if err := json.Unmarshal(msg.Payload, vc); err != nil {
		panic(fmt.Errorf("===>[ERROR from ViewChange]Invalid[%s] ViewChange message[%s]", err, msg))
	}

	nextPrimaryID := vc.NewViewID % util.TotalNodeNum
	if s.NodeID != nextPrimaryID {
		fmt.Printf("===>[ViewChange]I'm Node[%d] not the new[%d] primary node\n", s.NodeID, nextPrimaryID)
		return nil
	}

	s.sCache.pushVC(vc)
	if len(s.sCache.vcMsg) < 2*util.MaxFaultyNode+1 {
		return nil
	}
	if s.sCache.hasNewViewYet(vc.NewViewID) {
		fmt.Printf("===>[ViewChange]view change[%d] is in processing......\n", vc.NewViewID)
		return nil
	}

	return s.createNewViewMsg(vc.NewViewID)
}

func (s *StateEngine) createNewViewMsg(newVID int64) error {
	// handle newView message
	s.CurViewID = newVID
	nv := &message.NewView{
		NewViewID: s.CurViewID,
		VMsg:      s.sCache.vcMsg,
	}

	s.sCache.addNewView(nv)
	s.PrimaryID = s.CurViewID % util.TotalNodeNum
	fmt.Printf("\n===>[NewView] New primary is me[%d].....\n", s.PrimaryID)

	sk := s.P2pWire.GetMySecretkey()
	msg := message.CreateConMsg(message.MTNewView, nv, sk, s.NodeID)
	if err := s.P2pWire.BroadCast(msg); err != nil {
		return err
	}
	fmt.Printf("===>[NewView]Send NewView message success\n")

	s.cleanLogandRequest()
	s.stage = Collect
	s.nodeStatus = Serving
	s.CollectNum = 0
	s.TimeCommitmentCollect = make(map[int64]int)
	s.TimeCommitmentPropose = make(map[string]bool)
	// s.CollectTimer.tick(1 * time.Second)
	s.GlobalTimer.tick(10 * time.Minute)

	return nil
}

// invoked by state.go when received a newview message
func (s *StateEngine) didChangeView(msg *message.ConMessage) error {
	fmt.Printf("===>[NewView]Current Propose Message from Node[%d]\n", msg.From)
	nodeConfig := config.GetConsensusNode()

	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from didChangeView]Parse elliptic curve error:%s", err))
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
		panic(fmt.Errorf("===>[ERROR from didChangeView]Verify new public key Signature failed, From Node[%d]", msg.From))
	}
	fmt.Printf("===>[NewView]Verify success\n")

	// unmarshal message
	nv := &message.NewView{}
	if err := json.Unmarshal(msg.Payload, nv); err != nil {
		panic(fmt.Errorf("===>[ERROR from didChangeView]Invalid[%s] Propose message[%s]", err, msg))
	}

	s.GlobalTimer.tick(10 * time.Minute)
	s.CurViewID = nv.NewViewID
	s.sCache.vcMsg = nv.VMsg
	s.sCache.addNewView(nv)
	s.PrimaryID = s.CurViewID % util.TotalNodeNum
	s.nodeStatus = Serving
	s.stage = Propose
	fmt.Printf("===>[NewView]New primary is(%d).....\n", s.PrimaryID)

	s.cleanLogandRequest()
	// go s.reSendCollectMsg()
	s.reSendCollectMsg()
	return nil
}

func (s *StateEngine) cleanLogandRequest() {
	s.sCache = NewVCCache()
}

func (s *StateEngine) reSendCollectMsg() {
	time.Sleep(5 * time.Second)

	// new Collect message
	// send Collect message to primary node
	var tc [4]string
	for _, value := range s.TimeCommitment {
		time.Sleep(500 * time.Millisecond)

		tc = value
		Collect := &message.Collect{
			Length:    len(s.TimeCommitment),
			CollectTC: tc,
		}
		sk := s.P2pWire.GetMySecretkey()
		sMsg := message.CreateConMsg(message.MTCollect, Collect, sk, s.NodeID)
		conn := s.P2pWire.GetPrimaryConn(s.PrimaryID)
		if err := s.P2pWire.SendUniqueNode(conn, sMsg); err != nil {
			panic(fmt.Errorf("===>[ERROR from reSendCollectMsg]Send message error:%s", err))
		}

	}

	s.stage = Propose
	fmt.Printf("\n===>[Collect]Send Collect message success\n")
}
