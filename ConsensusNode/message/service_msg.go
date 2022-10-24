package message

import (
	"crypto/sha256"
	"fmt"
)

type MType int16

// number different kinds of message types
const (
	MTRequest MType = iota
	MTPrePrepare
	MTPrepare
	MTCommit
	MTViewChange
	MTNewView
	MTIdentity
)

// to be modified
const MaxFaultyNode = 2
const TotalNodeNum = 3*MaxFaultyNode + 1

// Hash message v, SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}

// Get Port(30000 + id)
func PortByID(id int64) int {
	return 30000 + int(id)
}

// MType.String()
func (mt MType) String() string {
	switch mt {
	case MTRequest:
		return "Request"
	case MTPrePrepare:
		return "PrePrepare"
	case MTPrepare:
		return "Prepare"
	case MTCommit:
		return "Commit"
	case MTViewChange:
		return "ViewChange"
	case MTNewView:
		return "NewView"
	case MTIdentity:
		return "PublicKey"
	}
	return "Unknown"
}

// message type from client
type ClientMessage struct {
	Sig       []byte `json:"sig"`
	TimeStamp int64  `json:"timestamp"`
	ClientID  string `json:"clientID"`
	Operation string `json:"operation"`
	PublicKey []byte `json:"pk"`
}

// request type in consensus
type Request struct {
	SeqID     int64  `json:"sequenceID"`
	TimeStamp int64  `json:"timestamp"`
	ClientID  string `json:"clientID"`
	Operation string `json:"operation"`
}

// request.String()
func (r *Request) String() string {
	return fmt.Sprintf("\n clientID:%s"+
		"\n time:%d"+
		"\n operation:%s",
		r.ClientID,
		r.TimeStamp,
		r.Operation)
}

// reply type in consensus
type Reply struct {
	SeqID     int64  `json:"sequenceID"`
	ViewID    int64  `json:"viewID"`
	Timestamp int64  `json:"timestamp"`
	ClientID  string `json:"clientID"`
	NodeID    int64  `json:"nodeID"`
	Result    string `json:"result"`
}
