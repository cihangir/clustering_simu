package main

import (
	"fmt"
	"sync"

	"github.com/kr/pretty"
	"github.com/looplab/fsm"
	uuid "github.com/satori/go.uuid"
)

func NewNetWork() *Network {
	return &Network{
		entryCluster: NewCluster(),
		joinChan:     make(chan *Node),
	}
}

type Network struct {
	entryCluster     *Cluster
	joinChan         chan *Node
	networkEventHubs []chan *NetworkEvent
}

func (n *Network) Missing(e *fsm.Event) {
	fmt.Println("there is a missing slot, we gotto fix it..")
	fmt.Printf("e %# v \n ", pretty.Formatter(e))
}

type NetworkEvent struct {
	id   string
	Name string
}

type Peer interface {
	SetNetworkEventHub(<-chan *NetworkEvent)
	SetRandomJoinHub(<-chan *Node)
}

func (n *Network) AddNode(peer *Node) error {
	peer.SetRandomJoinHub(n.joinChan)

	// TODO: change this to a gossip based system.
	neh := make(chan *NetworkEvent)
	peer.SetNetworkEventHub(neh)
	n.networkEventHubs = append(n.networkEventHubs, neh)

	peer.network = n
	go func() {
		/////////////////////////////////
		addNetworkJoinLatency()
		/////////////////////////////////
		n.joinChan <- peer
	}()

	return nil
}

func (n *Network) SendMessageToNetwork(msg *NetworkEvent) error {
	if msg.id == "" {
		msg.id = uuid.Must(uuid.NewV4()).String()
	}
	const concurrency = 3
	workChan := make(chan chan *NetworkEvent)
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			for nehChan := range workChan {
				go func(nehChan chan *NetworkEvent) {
					/////////////////////////////////
					addNetworkEventLatency()
					/////////////////////////////////
					select {
					case nehChan <- msg:
					default:
						// TODO: delete the chans that cant receive msgs
						fmt.Println("cant send msg-->")
					}
				}(nehChan)
			}
			wg.Done()
		}()
	}
	go func() {
		for _, nehChan := range n.networkEventHubs {
			workChan <- nehChan
		}
		close(workChan)
	}()

	wg.Wait()
	return nil
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
