package main

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/looplab/fsm"
)

type State uint

const (
	Missing State = iota + 1
	Critical
	NotComplete
	Complete

	HasQueue
)

type Event uint

const (
	StateTransition Event = iota + 1
	NewSlotAdded
	SlotRemoved
	NeedSlot
)

type Database struct {
	data map[string]map[string]string // map[slot_id]{key:val}
}

func NewNode() *Node {
	n := &Node{}
	nodeFSM := fsm.NewFSM(
		"created",
		fsm.Events{
			fsm.EventDesc{Name: "setup", Src: []string{"created"}, Dst: "looking for a cluster"},
			fsm.EventDesc{Name: "finds cluster", Src: []string{"complete"}, Dst: "joining"}, // we might not be able join to a cluster if the cluster is full.
			fsm.EventDesc{Name: "joins", Src: []string{"missing"}, Dst: "starts syncing"},
			fsm.EventDesc{Name: "sync fails", Src: []string{"missing"}, Dst: "starts syncing"},
			fsm.EventDesc{Name: "sync succeeds", Src: []string{"missing"}, Dst: "member"},

			fsm.EventDesc{Name: "finds slave cluster", Src: []string{"complete"}, Dst: "joining as slave"},
			fsm.EventDesc{Name: "joins as slave", Src: []string{"missing"}, Dst: "starts slave syncing"},
			fsm.EventDesc{Name: "slave sync fails", Src: []string{"missing"}, Dst: "starts slave syncing"},
			fsm.EventDesc{Name: "slave sync succeeds", Src: []string{"missing"}, Dst: "slave"},
			fsm.EventDesc{Name: "promotes to member", Src: []string{"missing"}, Dst: "member"}, // event names and states are different...
		},
		fsm.Callbacks{
			"setup":                 n.SetupNode,
			"looking for a cluster": n.LookForCluster,
			"sync fails":            n.SyncFails,
			"sync succeeds":         n.SyncSucceeds,
			"finds slave cluster":   n.FindsSlaveCluster,
			"joins as slave":        n.JoinsAsSlave,
			"slave sync fails":      n.SlaveSyncFails,
			"slave sync succeeds":   n.SlaveSyncSucceeds,
			"promotes to member":    n.PromotesToMember,
			"after_event":           n.LogTransition,
		},
	)
	n.fsm = nodeFSM
	return n
}

type Node struct {
	fsm *fsm.FSM

	parent *Cluster

	// if a node becomes a slave of another cluster, we will store their data as
	// well
	slaveOf *Cluster
	data    Database
}

func (n *Node) SetupNode(e *fsm.Event) {
	fmt.Println("SetupNode is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (n *Node) LookForCluster(e *fsm.Event) {
	fmt.Println("LookForCluster is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}
func (n *Node) LogTransition(e *fsm.Event) {
	fmt.Println("transitioned to " + e.FSM.Current())
}

func (n *Node) SyncFails(e *fsm.Event) {
	fmt.Println("SyncFails is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (n *Node) SyncSucceeds(e *fsm.Event) {
	fmt.Println("SyncSucceeds is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (n *Node) FindsSlaveCluster(e *fsm.Event) {
	fmt.Println("FindsSlaveCluster is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (n *Node) JoinsAsSlave(e *fsm.Event) {
	fmt.Println("JoinsasSlave is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (n *Node) SlaveSyncFails(e *fsm.Event) {
	fmt.Println("SlaveSyncFails is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (n *Node) SlaveSyncSucceeds(e *fsm.Event) {
	fmt.Println("SlaveSyncSucceeds is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

func (n *Node) PromotesToMember(e *fsm.Event) {
	fmt.Println("PromotesToMember is called with:")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}
