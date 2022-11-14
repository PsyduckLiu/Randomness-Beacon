package message

import (
	"consensusNode/signature"
	"crypto/ecdsa"
)

// create new public key message of backups
func CreateIdentityMsg(t MType, id int64, localAddr string, sk *ecdsa.PrivateKey) *ConMessage {
	// sign message.Payload
	sig := signature.GenerateSig([]byte(localAddr), sk)
	identityMsg := &ConMessage{
		Typ:     t,
		Sig:     sig,
		From:    id,
		Payload: []byte(localAddr),
	}

	return identityMsg
}
