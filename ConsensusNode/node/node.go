package node

import (
	"consensusNode/consensus"
	"fmt"
)

// [signal]:  a channel connects [node] with [consensus], deliver the exit message
// [conChan]: a channel connects [node] with [consensus]
type Node struct {
	NodeID    int64
	signal    chan interface{}
	consensus *consensus.StateEngine
}

// initialize a new node
func NewNode(id int64) *Node {
	c := consensus.InitConsensus(id)
	n := &Node{
		NodeID:    id,
		consensus: c,
		signal:    make(chan interface{}),
	}

	return n
}

// run a node
func (n *Node) Run() {
	fmt.Printf("===>[NewNode]Consensus node[%d] start......\n", n.NodeID)

	go n.consensus.StartConsensus(n.signal)
	go n.consensus.WatchConfig(n.NodeID, n.signal)

	s := <-n.signal
	fmt.Printf("===>[Node EXIT]Node[%d] exit because of:%s\n", n.NodeID, s)
}
