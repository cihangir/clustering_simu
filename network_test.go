package main

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/kr/pretty"
)

func TestNetwork_AddNode(t *testing.T) {
	count := 5
	maxNetworkJoinLatency = 0
	maxNetworkEventLatency = 0

	network := NewNetWork()
	arr := make([]*Node, count)

	for i := 0; i < count; i++ {
		n := NewNode()
		if err := network.AddNode(n); err != nil {
			t.Errorf("Network.AddNode() error = %v", err)
		}
		arr[i] = n
	}
	time.Sleep(time.Second)
	// calculate the cluster count
	clusterNames := make(map[string]int)
	for _, item := range arr {
		clusterNames[item.mainCluster.id]++
	}

	// check if all the clusters has the same number of nodes attached
	for _, item := range arr {
		if clusterNames[item.mainCluster.id] != len(item.mainCluster.nodes) {

			t.Errorf("clusterNames[item.mainCluster.id] != len(item.mainCluster.nodes) [%d] != [%d]", clusterNames[item.mainCluster.id], len(item.mainCluster.nodes))
			fmt.Printf("item %# v \n ", pretty.Formatter(item))
		}
	}

	if float64(len(clusterNames)) > math.Ceil(float64(count/arr[0].mainCluster.clusterSize)) {
		t.Errorf("len(clusterNames) [%d] > count/arr[0].mainCluster.clusterSize [%d]", len(clusterNames), count/arr[0].mainCluster.clusterSize)
	}
}
