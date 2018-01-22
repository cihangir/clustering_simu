package main

import (
	"fmt"
	"sync"

	"github.com/looplab/fsm"
	uuid "github.com/satori/go.uuid"
)

func NewCluster() *Cluster {
	c := &Cluster{
		clusterSize:      5,
		criticalSize:     3,
		id:               uuid.Must(uuid.NewV4()).String(),
		intraClusterConn: make(chan *ClusterMessage),
	}

	c.fsm = newClusterFSM(c)
	return c
}

type Cluster struct {
	id     string
	fsm    *fsm.FSM
	master *Node

	// protects the followings
	mu sync.Mutex

	// TODO: change arrays to map, we need uniqueness
	nodes   []*Node
	inQueue []*Node

	clusterSize  int
	criticalSize int

	intraClusterConn chan *ClusterMessage
}

type ClusterMessage struct {
	Type string
	Data string
}

// MoveToNewCluster forms a new cluster with a subset of the existing cluster
// members. This is useful where some of the nodes are slaves of another
// cluster and the corresponding cluster is in critical state.
func (c *Cluster) MoveToNewCluster(nodes ...*Node) error {
	// TODO: check if on-move are we going into critical state
	newCluster := NewCluster()
	for _, node := range nodes {
		if err := newCluster.addNode(node); err != nil {
			return err
		}
		if err := c.removeNode(node); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) addNode(n *Node) error {
	// if we have enough nodes, add them to the waiting line
	c.mu.Lock()

	if len(c.nodes) >= c.clusterSize {
		c.inQueue = append(c.inQueue, n)
		// this event is useful for signalling the network that we have more
		// nodes than required, if someone needs one, we can transfer some to
		// them
		if err := n.SendEventToNetwork(&NetworkEvent{
			Name: "added queue node",
		}); err != nil {
			fmt.Println("err while sending event to network-->", err)
			return err
		}
	} else {
		c.nodes = append(c.nodes, n)
		if len(c.nodes) == c.clusterSize {
			if err := c.fsm.Event("complete", n); err != nil {
				fmt.Println("c.fsm.Event(complete)-->", err)
				return err
			}
		}
	}
	n.mainCluster = c

	c.mu.Unlock()

	// if we have enough node count to form a cluster, no need to make them wait
	// in our queue
	if len(c.inQueue) >= c.criticalSize {
		return c.MoveToNewCluster(c.inQueue...)
	}

	return nil
}

func (c *Cluster) removeNode(n *Node) error {
	if n == nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	nodes := make([]*Node, 0, len(c.nodes))
	for _, nq := range c.nodes {
		if nq.id != n.id {
			nodes = append(nodes, nq)
		}
	}
	c.nodes = nodes

	inQueue := make([]*Node, 0, len(c.nodes))
	for _, nq := range c.inQueue {
		if nq.id != n.id {
			inQueue = append(inQueue, nq)
		}
	}
	c.inQueue = inQueue

	return nil
}

// SendMessageToCluster sends a message to the all members of this cluster
// TODO:
// change this sync code to accept a channel and send messages over that
func (c *Cluster) SendMessageToCluster(msg *ClusterMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// fmt.Printf("c.nodes %# v \n ", pretty.Formatter(c.nodes))
	for _, node := range c.nodes {
		go func(node *Node) {
			////////////////////////////////////////////////////////
			//// latency between cluster member's communication/////
			addClusterMessageLatency() /////////////////////////////
			////////////////////////////////////////////////////////
			////////////////////////////////////////////////////////
			if err := node.HandleClusterMessages(msg); err != nil {
				fmt.Println("err-->", err)
			}
		}(node)
	}
	for _, node := range c.inQueue {
		go func(node *Node) {
			////////////////////////////////////////////////////////
			//// latency between cluster member's communication/////
			addClusterMessageLatency() /////////////////////////////
			////////////////////////////////////////////////////////
			////////////////////////////////////////////////////////
			if err := node.HandleClusterMessages(msg); err != nil {
				fmt.Println("err-->", err)
			}
		}(node)
	}
	return nil
}
