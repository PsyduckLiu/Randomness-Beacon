package message

import (
	"crypto/sha256"
	"fmt"
)

type MType int16

// number different kinds of message types
const (
	MTVRFVerify MType = iota
	MTCollect
	MTSubmit
	MTApprove
	MTViewChange
	MTNewView
	MTIdentity
)

// MType.String()
func (mt MType) String() string {
	switch mt {
	case MTVRFVerify:
		return "VRFVerify"
	case MTCollect:
		return "Collect"
	case MTSubmit:
		return "Submit"
	case MTApprove:
		return "Approve"
	case MTViewChange:
		return "ViewChange"
	case MTNewView:
		return "NewView"
	case MTIdentity:
		return "Identity"
	}
	return "Unknown"
}

// message type for entropy node
type EntropyVRFMessage struct {
	PublicKey [32]byte `json:"pk"`
	VRFResult [80]byte `json:"vrfresult"`
	ClientID  int64    `json:"clientID"`
	Msg       []byte   `json:"timecommitmentmsg"`
}

type EntropyTCMessage struct {
	ClientID               int64  `json:"clientID"`
	TimeCommitmentC        string `json:"timecommitmentC"`
	TimeCommitmentH        string `json:"timecommitmentH"`
	TimeCommitmentrKSubOne string `json:"timecommitmentrKSubOne"`
	TimeCommitmentrK       string `json:"timecommitmentrK"`
}

// Hash message v, SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}
