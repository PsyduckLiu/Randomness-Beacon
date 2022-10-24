package message

import (
	"crypto/ecdsa"
	"crypto/x509"
	"entropyNode/signature"
	"fmt"
	"time"
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

// create new public key message of the client
func CreateClientKeyMsg(sk *ecdsa.PrivateKey) *ClientMessage {
	publicKey := sk.PublicKey
	marshalledKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	keyMsg := &ClientMessage{
		TimeStamp: time.Now().Unix(),
		ClientID:  "Client's address",
		Operation: "<READ TX FROM POOL>",
		PublicKey: marshalledKey,
	}

	// sign the whole ClientMessage without signature
	sig := signature.GenerateSig([]byte(fmt.Sprintf("%v", keyMsg)), sk)
	keyMsg.Sig = sig
	return keyMsg
}
