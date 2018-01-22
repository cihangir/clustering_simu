package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/kr/pretty"
	"github.com/looplab/fsm"
)

func NewNode() *Node {
	n := &Node{
		id:                  uuid.Must(uuid.NewV4()).String(),
		mainCluster:         NewCluster(),
		latestHealthChecks:  make(map[string]time.Time),
		healthCheckTimeout:  time.Second * 10,
		healthCheckInterval: time.Second,
		closeChan:           make(chan struct{}),
	}

	n.fsm = newNodeFSM(n)
	go n.sendHealthChecks()
	go n.checkMemberships()

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

	healthMux           sync.RWMutex
	latestHealthChecks  map[string]time.Time
	healthCheckTimeout  time.Duration
	healthCheckInterval time.Duration

	closeMu   sync.Mutex
	closed    bool
	closeChan chan struct{}
}

func (n *Node) Close() error {
	n.closeMu.Lock()
	defer n.closeMu.Unlock()
	if n.closed {
		return errors.New("already closed")
	}
	n.closed = true
	close(n.closeChan)
	return nil
}

func (n *Node) HandleClusterMessages(msg *ClusterMessage) error {
	switch msg.Type {
	case "healthcheck":
		return n.handleHealthCheck(msg)
	}
	return nil
}

func (n *Node) handleHealthCheck(msg *ClusterMessage) error {
	if msg.NodeID == "" {
		return fmt.Errorf("not id is not set. msg: %+v", msg)
	}
	// fmt.Println("msg.NodeID-->", msg.NodeID, "msg.Type", msg.Type)
	n.healthMux.Lock()
	n.latestHealthChecks[msg.NodeID] = time.Now().UTC()
	n.healthMux.Unlock()
	return nil
}

func (n *Node) sendHealthChecks() {
	tc := time.NewTicker(n.healthCheckInterval)
	defer tc.Stop()
	for {
		select {
		case <-tc.C:
			n.mainCluster.SendMessageToCluster(&ClusterMessage{
				Type:   "healthcheck",
				NodeID: n.id,
				Data:   "hede hodo",
			})
		case <-n.closeChan:
			fmt.Println("stopping sending healthchecks-->")
			return
		}
	}
}

func (n *Node) checkMemberships() {
	tc := time.NewTicker(n.healthCheckInterval)
	defer tc.Stop()
	for {
		select {
		case <-tc.C:
			for _, node := range n.mainCluster.nodes {
				// TODO update locks when we actually start to delete the node
				n.healthMux.RLock()
				h, ok := n.latestHealthChecks[node.id]
				n.healthMux.RUnlock()

				if !ok {
					fmt.Printf("node belongs to cluster: %+v, but not sending health checks: %+v", n.mainCluster, node)
					return
				}

				if h.Add(n.healthCheckTimeout).Before(time.Now().UTC()) {
					fmt.Println("node is dead, needs clean up-->", node)
				}
			}

			n.healthMux.RLock()
			for id, t := range n.latestHealthChecks {
				if t.Add(n.healthCheckTimeout).Before(time.Now().UTC()) {
					fmt.Println("node id is dead, needs clean up-->", id)
				}
			}
			n.healthMux.RUnlock()

		case <-n.closeChan:
			fmt.Println("stopping checking members-->")
			return
		}
	}
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
