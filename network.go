package main

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/looplab/fsm"
)

type Network struct {
}

func (n *Network) Missing(e *fsm.Event) {
	fmt.Println("there is a missing slot, we gotto fix it..")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

// FindSlot finds the slot within the network, if slot is not available, an
// error will be returned
func (n *Network) FindSlot(slot int) (*Cluster, error) {
	return nil, nil
}

// RecoverSlot tries to re-create a slot within a cluster. First, it tries to
// create it from any possible slave nodes. If there isn't any slave nodes, it
// might try to recreate it from a persistent storage.
func (n *Network) RecoverSlot(slot int) (*Cluster, error) {
	return nil, nil
}
