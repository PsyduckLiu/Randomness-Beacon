package message

import (
	"crypto/sha256"
	"fmt"
)

type MType int16

// number different kinds of message types
const (
	MTCollect MType = iota
	MTSubmit
	MTApprove
	MTViewChange
	MTNewView
	MTIdentity
)

// MType.String()
func (mt MType) String() string {
	switch mt {
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
type EntropyMessage struct {
	PublicKey      []byte `json:"pk"`
	VRFResult      []byte `json:"vrfresult"`
	TimeStamp      int64  `json:"timestamp"`
	ClientID       int64  `json:"clientID"`
	TimeCommitment string `json:"timecommitment"`
}

// Hash message v, SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}
