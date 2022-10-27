package node

import (
	"consensusNode/consensus"
	"consensusNode/message"
	"consensusNode/service"
	"fmt"
)

const MaxMsgNum = 100

// [signal]:  a channel connects [node] with [consensus] and [service], deliver the exit message
// [srvChan]: a channel connects [node] with [service], deliver service message(request)
// [conChan]: a channel connects [node] with [consensus], deliver {message.RequestRecord} message to clients
// [directReplyChan]: a channel connects [node] with [consensus], deliver {message.Reply} message to clients
type Node struct {
	NodeID int64
	signal chan interface{}
	// srvChan         chan interface{}
	conChan         <-chan *message.RequestRecord
	directReplyChan <-chan *message.Reply
	// waitQueue       []*message.Request
	consensus *consensus.StateEngine
	service   *service.Service
}

// initialize a new node
func NewNode(id int64) *Node {
	// srvChan := make(chan interface{}, MaxMsgNum)
	conChan := make(chan *message.RequestRecord, MaxMsgNum)
	// rChan := make(chan *message.Reply, MaxMsgNum)

	c := consensus.InitConsensus(id, conChan)
	// sr := service.InitService(util.PortByID(id), srvChan)

	n := &Node{
		NodeID:    id,
		consensus: c,
		// service:   sr,
		// srvChan: srvChan,
		// waitQueue:       make([]*message.Request, 0),
		signal:  make(chan interface{}),
		conChan: conChan,
		// directReplyChan: rChan,
	}

	return n
}

// run a node
func (n *Node) Run() {
	fmt.Printf("===>Consensus node[%d] start......\n", n.NodeID)

	// go config.WatchConfig()
	go n.consensus.StartConsensus(n.signal)
	// go n.consensus.WaitTC(n.signal)
	go n.consensus.WatchConfig(n.NodeID, n.signal)
	// go n.service.WaitRequest(n.signal, n.consensus)
	// go n.Dispatch()

	s := <-n.signal
	fmt.Printf("===>[EXIT]Node[%d] exit because of:%s", n.NodeID, s)
}
