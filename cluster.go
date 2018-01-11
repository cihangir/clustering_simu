package main

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/looplab/fsm"
)

func NewCluster() *Cluster {
	c := &Cluster{}
	clusterFSM := fsm.NewFSM(
		"idle",
		fsm.Events{
			fsm.EventDesc{Name: "node leaves", Src: []string{"idle", "not complete", "critical", "complete"}, Dst: "node leaving"},
			fsm.EventDesc{Name: "node joins", Src: []string{"idle", "not complete", "critical", "complete"}, Dst: "node joining"},
			fsm.EventDesc{Name: "critical", Src: []string{"idle", "not complete", "critical", "missing"}, Dst: "requesting slave node"},
			fsm.EventDesc{Name: "missing", Src: []string{"idle", "not complete", "critical", "missing"}, Dst: "recovering"},
			fsm.EventDesc{Name: "unknown", Src: []string{"idle", "not complete", "critical", "missing"}, Dst: "recovering"},
			fsm.EventDesc{Name: "complete", Src: []string{"idle", "not complete", "critical", "missing"}, Dst: "complete"}, // event names and states are different...
		},
		fsm.Callbacks{
			"node leaving":             c.NodeLeaving,
			"node joining":             c.NodeJoining,
			"requesting node":          c.RequestingMember,
			"requesting slave node":    c.RequestingSlave,
			"requesting health checks": c.RequestingHealthChecks,
			"recovering":               c.RecoveringCluster,
			"missing":                  c.RecoveringCluster,

			"critical":    c.OnCriticalState,
			"complete":    c.Complete,
			"after_event": c.LogTransition,
		},
	)

	c.fsm = clusterFSM
	return c
}

type Cluster struct {
	fsm     *fsm.FSM
	network *Network
	nodes   []*Node
	master  *Node
	inQueue []*Node

	clusterSize  int
	criticalSize int
}

func (c *Cluster) NodeLeaving(e *fsm.Event) {
	if len(c.nodes) == 0 {
		e.FSM.Event("missing")
		return
	}

	// we are critical state, should recover ASAP.
	if len(c.nodes) <= c.criticalSize {
		e.FSM.Event("critical")
		return
	}

	if len(c.nodes) == c.clusterSize {
		e.FSM.Event("complete")
		return
	}

	e.FSM.Event("unknown")
}

func (c *Cluster) OnCriticalState(e *fsm.Event) {
	// health checks are required for external validation
	e.FSM.Event("requesting health checks")

	// we gotta get a node - if possible. If not possible, we will get a slave
	// not from some other clusters
	e.FSM.Event("requesting node")

	e.FSM.Event("requesting slave node")
	// do other stuff
}

func (c *Cluster) OnNotCompleteState(e *fsm.Event) {
	// do other stuff
	e.FSM.Event("requesting node")
	// do other stuff
}

func (c *Cluster) RequestingMember(e *fsm.Event) {
	// we can continue to function without a node, but it is always better to
	// have a full quorum
	fmt.Println("requesting a new member-->")

	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (c *Cluster) RequestingSlave(e *fsm.Event) {
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (c *Cluster) RequestingHealthChecks(e *fsm.Event) {
	// we should request health checks from nearest possible clusters, if they
	// are not on critical state
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (c *Cluster) NodeJoining(e *fsm.Event) {
	// TODO: while a new node is joining, if we have a slave node that is
	// already subscribed to our events, direct the new instance to join to the
	// other cluster and request the slave node.

	// If the other cluster rejects the request, continue as is.

	// If the other cluster accepts the request, this cluster will have it as
	// full member, but it will continue operating as slave member to the
	// previous cluster.

	// Removing the new node from the slave position is up the previous cluster.
	// Most probably, the other cluster will remove the node from being slave
	// when the new node starts operating properly.
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (c *Cluster) RecoveringCluster(e *fsm.Event) {
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (c *Cluster) Complete(e *fsm.Event) {
	// TODO: notify other "listening members" that now we are on complete state,
	// there is nothing to worry about us anymore.
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (c *Cluster) LogTransition(e *fsm.Event) {
	fmt.Println("transitioned to " + e.FSM.Current())
}

// MoveToNewCluster forms a new cluster with a subset of the existing cluster
// members. This is useful where some of the nodes are slaves of another
// cluster and the corresponding cluster is in critical state.
func (c *Cluster) MoveToNewCluster(nodes ...*Node) error {
	// TODO: check if on-move are we going into critical state
	return nil
}
