package consensus

import "github.com/algorand/go-algorand/protocol"

type Stage int

// number different kinds of stage types
const (
	Collect Stage = iota
	Submit
	Approve
	Confirm
	Output
	Error
	ViewChange
)

// stage.string()
func (s Stage) String() string {
	switch s {
	case Collect:
		return "Collect"
	case Submit:
		return "Submit"
	case Approve:
		return "Approve"
	case Confirm:
		return "Confirm"
	case Output:
		return "Output"
	case Error:
		return "Error"
	case ViewChange:
		return "ViewChange"
	}

	return "Unknown"
}

type EngineStatus int8

// number different kinds of EngineStatus types
const (
	Serving EngineStatus = iota
	ViewChanging
)

// EngineStatus.string()
func (es EngineStatus) String() string {
	switch es {
	case Serving:
		return "Server consensus......"
	case ViewChanging:
		return "Changing views......"
	}
	return "Unknown"
}

type MessageHashable struct {
	Data []byte
}

func (s MessageHashable) ToBeHashed() (protocol.HashID, []byte) {
	return "msg", s.Data
}
