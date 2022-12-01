package message

type MType int16

// number different kinds of message types
const (
	MTVRFVerify MType = iota
	MTCommitFromEntropy
	MTCollect
	MTPropose
	MTPrepare
	MTReveal
	MTViewChange
	MTNewView
	MTIdentity
)

// MType.String()
func (mt MType) String() string {
	switch mt {
	case MTVRFVerify:
		return "VRFVerify"
	case MTCommitFromEntropy:
		return "CommitFromEntropyt"
	case MTCollect:
		return "Collect"
	case MTPropose:
		return "Propose"
	case MTPrepare:
		return "Prepare"
	case MTReveal:
		return "Reveal"
	case MTViewChange:
		return "ViewChange"
	case MTNewView:
		return "NewView"
	case MTIdentity:
		return "Identity"
	}
	return "Unknown"
}

// VRF message type for entropy node
type EntropyVRFMessage struct {
	PublicKey [32]byte `json:"pk"`
	VRFResult [80]byte `json:"vrfresult"`
	ClientID  int64    `json:"clientID"`
	Msg       []byte   `json:"timecommitmentmsg"`
}

// TC message type for entropy node
type EntropyTCMessage struct {
	ClientID               int64  `json:"clientID"`
	TimeCommitmentC        string `json:"timecommitmentC"`
	TimeCommitmentH        string `json:"timecommitmentH"`
	TimeCommitmentrKSubOne string `json:"timecommitmentrKSubOne"`
	TimeCommitmentrK       string `json:"timecommitmentrK"`
	TimeCommitmentA1       string `json:"timecommitmentA1"`
	TimeCommitmentA2       string `json:"timecommitmentA2"`
	TimeCommitmentA3       string `json:"timecommitmentrKSubOnA3"`
	TimeCommitmentZ        string `json:"timecommitmentrZ"`
}
