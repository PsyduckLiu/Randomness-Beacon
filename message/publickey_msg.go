package message

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"randomnessBeacon/signature"
	"time"
)

// create new public key message of backups
func CreateKeyMsg(t MType, id int64, sk *ecdsa.PrivateKey) *ConMessage {
	publicKey := sk.PublicKey
	marshalledKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	// sign message.Payload(marshalledKey)
	sig := signature.GenerateSig(marshalledKey, sk)
	keyMsg := &ConMessage{
		Typ:     t,
		Sig:     sig,
		From:    id,
		Payload: marshalledKey,
	}

	return keyMsg
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
