package consensus

import "consensusNode/message"

type NormalLog struct {
	ClientID   string                    `json:"clientID"`
	Stage      Stage                     `json:"Stage"`
	PrePrepare *message.PrePrepare       `json:"PrePrepare"`
	Prepare    message.PrepareMsg        `json:"Prepare"`
	Commit     map[int64]*message.Commit `json:"Commit"`
}

// func NewNormalLog() *NormalLog {
// 	nl := &NormalLog{
// 		Stage:      Idle,
// 		PrePrepare: nil,
// 		Prepare:    make(message.PrepareMsg),
// 		Commit:     make(map[int64]*message.Commit),
// 	}
// 	return nl
// }
