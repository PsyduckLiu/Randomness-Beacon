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
		"\n--------------------",
		cm.Typ.String(),
		cm.Sig,
		string(rune(cm.From)),
		len(cm.Payload))
}

// create consensus message
func CreateConMsg(t MType, msg interface{}, sk *ecdsa.PrivateKey, id int64) *ConMessage {
	data, err := json.Marshal(msg)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from CreateConMsg]Generate consensus message failed:%s", err))
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

// Collect type in Consensus
type Collect struct {
	Length    int
	CollectTC [4]string
	TCProof   [4]string
}

// Propose type in Consensus
type Propose struct {
	UnionTC [4]string
	Length  int
}

// Prepare type in Consensus
type Prepare struct {
	// PrepareTC [4]string
	Length int
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
