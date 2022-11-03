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

func NewVCCache() *VCCache {
	return &VCCache{
		vcMsg: make(message.VMessage),
		nvMsg: make(map[int64]*message.NewView),
	}
}

func (vcc *VCCache) pushVC(vc *message.ViewChange) {
	vcc.vcMsg[vc.NodeID] = vc
}

func (vcc *VCCache) hasNewViewYet(vid int64) bool {
	if _, ok := vcc.nvMsg[vid]; ok {
		return true
	}
	return false
}

func (vcc *VCCache) addNewView(nv *message.NewView) {
	vcc.nvMsg[nv.NewViewID] = nv
}

// invoked by state.go when timeout
func (s *StateEngine) ViewChange() {
	// fmt.Printf("======>[ViewChange] (%d, %d).....\n", s.CurViewID, s.lastCP.Seq)
	fmt.Printf("======>[ViewChange] Current view is(%d).....\n", s.CurViewID)
	// s.nodeStatus = ViewChanging
	// s.CollectTimer.tack()

	vc := &message.ViewChange{
		NewViewID: s.CurViewID + 1,
		NodeID:    s.NodeID,
	}

	nextPrimaryID := vc.NewViewID % util.TotalNodeNum
	if s.NodeID == nextPrimaryID {
		s.sCache.pushVC(vc) //[vc.NodeID] = vc
	}

	sk := s.P2pWire.GetMySecretkey()
	consMsg := message.CreateConMsg(message.MTViewChange, vc, sk, s.NodeID)
	// locker := &sync.RWMutex{}
	if err := s.P2pWire.BroadCast(consMsg); err != nil {
		fmt.Println(err)
		return
	}
	s.CurViewID++
	s.msgLogs = make(map[int64]*NormalLog)
}

// invoked by state.go when received a viewchage message
func (s *StateEngine) procViewChange(msg *message.ConMessage) error {
	fmt.Printf("======>[ViewChange]Current Approve Message from Node[%d]\n", msg.From)

	nodeConfig := config.GetConsensusNode()
	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ViewChangeERROR]Key message parse err:%s", err))
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
	fmt.Printf("======>[ViewChange]Verify success\n")

	// unmarshal message
	vc := &message.ViewChange{}
	if err := json.Unmarshal(msg.Payload, vc); err != nil {
		panic(fmt.Errorf("======>[Approve]Invalid[%s] Approve message[%s]", err, msg))
	}

	nextPrimaryID := vc.NewViewID % util.TotalNodeNum
	if s.NodeID != nextPrimaryID {
		fmt.Printf("I'm Node[%d] not the new[%d] primary node\n", s.NodeID, nextPrimaryID)
		return nil
	}

	s.sCache.pushVC(vc)
	if len(s.sCache.vcMsg) < 2*util.MaxFaultyNode {
		return nil
	}
	if s.sCache.hasNewViewYet(vc.NewViewID) {
		fmt.Printf("view change[%d] is in processing......\n", vc.NewViewID)
		return nil
	}

	return s.createNewViewMsg(vc.NewViewID)
}

func (s *StateEngine) createNewViewMsg(newVID int64) error {
	s.CurViewID = newVID
	nv := &message.NewView{
		NewViewID: s.CurViewID,
		VMsg:      s.sCache.vcMsg,
	}

	s.sCache.addNewView(nv)
	s.CurSequence = 0
	s.PrimaryID = s.CurViewID % util.TotalNodeNum
	fmt.Printf("======>[ViewChange] New primary is me[%d].....\n", s.PrimaryID)

	sk := s.P2pWire.GetMySecretkey()
	msg := message.CreateConMsg(message.MTNewView, nv, sk, s.NodeID)
	if err := s.P2pWire.BroadCast(msg); err != nil {
		return err
	}

	s.cleanLogandRequest()
	fmt.Printf("======>[ViewChange]Send NewView message success\n")

	go s.reSendApproveMsg()
	return nil
}

// invoked by state.go when received a newview message
func (s *StateEngine) didChangeView(msg *message.ConMessage) error {
	fmt.Printf("======>[NewView]Current Approve Message from Node[%d]\n", msg.From)

	nodeConfig := config.GetConsensusNode()
	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[NewViewERROR]Key message parse err:%s", err))
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
		panic(fmt.Errorf("===>[NewVieweERROR]Verify new public key Signature failed, From Node[%d]", msg.From))
	}
	fmt.Printf("======>[NewView]Verify success\n")

	// unmarshal message
	nv := &message.NewView{}
	if err := json.Unmarshal(msg.Payload, nv); err != nil {
		panic(fmt.Errorf("======>[NewView]Invalid[%s] Approve message[%s]", err, msg))
	}

	s.CurViewID = nv.NewViewID
	s.sCache.vcMsg = nv.VMsg
	s.sCache.addNewView(nv)
	s.CurSequence = 0
	s.PrimaryID = s.CurViewID % util.TotalNodeNum
	s.stage = Approve
	fmt.Printf("======>[NewView] New primary is(%d).....\n", s.PrimaryID)

	s.cleanLogandRequest()
	return nil
}

func (s *StateEngine) cleanLogandRequest() {
	s.msgLogs = make(map[int64]*NormalLog)
}

func (s *StateEngine) reSendApproveMsg() {
	time.Sleep(5 * time.Second)
	// new approve message
	// send approve message to backup nodes
	var tc []string
	for _, value := range s.TimeCommitment {
		tc = append(tc, value)
	}

	approve := &message.Approve{
		UnionTC: tc,
	}

	sk := s.P2pWire.GetMySecretkey()
	aMsg := message.CreateConMsg(message.MTApprove, approve, sk, s.NodeID)
	if err := s.P2pWire.BroadCast(aMsg); err != nil {
		panic(err)
	}
	s.stage = Confirm

	s.ConfirmNum++
	fmt.Printf("======>[Union]Send approve message success\n")
}
