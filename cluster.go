package main

import (
	"fmt"

	"github.com/looplab/fsm"
	uuid "github.com/satori/go.uuid"
)

func NewCluster() *Cluster {
	c := &Cluster{
		id: uuid.Must(uuid.NewV4()).String(),
	}

	c.fsm = newClusterFSM(c)
	return c
}

type Cluster struct {
	id      string
	fsm     *fsm.FSM
	network *Network
	nodes   []*Node
	master  *Node
	inQueue []*Node

	clusterSize  int
	criticalSize int
}

// MoveToNewCluster forms a new cluster with a subset of the existing cluster
// members. This is useful where some of the nodes are slaves of another
// cluster and the corresponding cluster is in critical state.
func (c *Cluster) MoveToNewCluster(nodes ...*Node) error {
	// TODO: check if on-move are we going into critical state
	return nil
}

func (c *Cluster) addNode(n *Node) error {
	// if we have enough nodes, add them to the waiting line
	if len(c.nodes) >= c.clusterSize {
		c.inQueue = append(c.inQueue, n)
		if err := n.SendEventToNetwork(&NetworkEvent{
			Name: "added queue node",
		}); err != nil {
			fmt.Println("err while sending event to network-->", err)
		}
	} else {
		c.nodes = append(c.nodes, n)
		if len(c.nodes) == c.clusterSize {
			if err := c.fsm.Event("complete", n); err != nil {
				fmt.Println("c.fsm.Event(complete)-->", err)
			}
		}
	}

	return nil
}
