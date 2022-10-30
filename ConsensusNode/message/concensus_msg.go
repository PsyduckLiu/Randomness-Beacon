package message

import (
	"consensusNode/signature"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
)

// consensus message type in Consensus
type ConMessage struct {
	Typ     MType  `json:"type"`
	Sig     []byte `json:"sig"`
	From    int64  `json:"from"`
	Payload []byte `json:"payload"`
}

// consensus message.String()
func (cm *ConMessage) String() string {
	return fmt.Sprintf("\n======Consensus Messagetype======"+
		"\ntype:%s"+
		"\nsig:%s"+
		"\nFrom:%s"+
		"\npayload:%d"+
		"\n<------------------>",
		cm.Typ.String(),
		cm.Sig,
		string(rune(cm.From)),
		len(cm.Payload))
}

// create consensus message
func CreateConMsg(t MType, msg interface{}, sk *ecdsa.PrivateKey, id int64) *ConMessage {
	data, e := json.Marshal(msg)
	if e != nil {
		return nil
	}

	// sign message.Payload
	sig := signature.GenerateSig(data, sk)
	consMsg := &ConMessage{
		Typ:     t,
		Sig:     sig,
		From:    id,
		Payload: data,
	}

	return consMsg
}

// RequestRecord type in Consensus
type Submit struct {
	CollectTC []string
}

// RequestRecord type in Consensus
type Approve struct {
	UnionTC []string
}

// RequestRecord type in Consensus
type RequestRecord struct {
	*PrePrepare
	*Request
}

// PrePrepare type in Consensus
type PrePrepare struct {
	// TimeStamp  int64  `json:"timestamp"`
	// ClientID   string `json:"clientID"`
	// Operation  string `json:"operation"`
	ViewID     int64  `json:"viewID"`
	SequenceID int64  `json:"sequenceID"`
	Digest     []byte `json:"digest"`
}

// Prepare array
// index is node ID
type PrepareMsg map[int64]*Prepare

// Prepare type in Consensus
type Prepare struct {
	ViewID     int64  `json:"viewID"`
	SequenceID int64  `json:"sequenceID"`
	Digest     []byte `json:"digest"`
	NodeID     int64  `json:"nodeID"`
}

// Commit type in Consensus
type Commit struct {
	ViewID     int64  `json:"viewID"`
	SequenceID int64  `json:"sequenceID"`
	Digest     []byte `json:"digest"`
	NodeID     int64  `json:"nodeID"`
}

// ViewChange array
type VMessage map[int64]*ViewChange

// ViewChange type in Consensus
type ViewChange struct {
	NewViewID int64 `json:"newViewID"`
	NodeID    int64 `json:"nodeID"`
}

// NewView type in Consensus
type NewView struct {
	NewViewID int64    `json:"newViewID"`
	VMsg      VMessage `json:"vMSG"`
}
