package message

import (
	"fmt"
)

type MType int16

// number different kinds of message types
const (
	MTCollect MType = iota
	MTSubmit
	MTApprove
	MTOutput
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
	case MTOutput:
		return "Output"
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
