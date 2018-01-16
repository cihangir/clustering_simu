package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/kr/pretty"
	"github.com/looplab/fsm"
)

func NewNode() *Node {
	n := &Node{
		id:          uuid.Must(uuid.NewV4()).String(),
		mainCluster: NewCluster(),
	}

	n.fsm = newNodeFSM(n)

	return n
}

type Node struct {
	id  string
	fsm *fsm.FSM

	// joinChan only listens for join events, this is used to simulate a random
	// node, joins to a random node as peer. All the nodes will listen from this
	// chan.
	joinChan    <-chan *Node
	networkChan <-chan *NetworkEvent

	mainCluster *Cluster

	// if a node becomes a slave of another cluster, we will store their data as
	// well
	slaveCluster *Cluster
	data         Database

	network *Network // TODO: remove network
}

func (n *Node) SetRandomJoinHub(hub <-chan *Node) {
	n.joinChan = hub
	go n.startAcceptingPeers()
}

func (n *Node) SetNetworkEventHub(neh <-chan *NetworkEvent) {
	n.networkChan = neh
	go func() {
		for msg := range n.networkChan {
			fmt.Printf("network chan. node id: %q msg: %# v \n ", n.id, pretty.Formatter(msg))
		}
	}()
}

func (n *Node) SendEventToNetwork(event *NetworkEvent) error {
	return n.network.SendMessageToNetwork(event)
}

func (n *Node) startAcceptingPeers() {
	for node := range n.joinChan {
		////////////////////////////////////
		// Set other peer's (node) cluster as this ones (n), because they join to us.
		////////////////////////////////////

		node.mainCluster = n.mainCluster
		if err := n.mainCluster.addNode(node); err != nil {
			fmt.Println("err while adding to main cluster-->", err)
		}
	}
}
