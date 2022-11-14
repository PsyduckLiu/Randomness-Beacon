package message

import (
	"crypto/ecdsa"
	"encoding/json"
	"entropyNode/signature"
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
