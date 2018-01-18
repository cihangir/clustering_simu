package main

import (
	"math"
	"testing"
	"time"
)

func TestNetwork_AddNode_same_clusterCount(t *testing.T) {
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

	time.Sleep(time.Second * 2)

	// calculate the cluster count
	clusterNames := make(map[string]int)
	for _, item := range arr {
		clusterNames[item.mainCluster.id]++
	}

	// check if all the clusters has the same number of nodes attached
	for _, item := range arr {
		if clusterNames[item.mainCluster.id] != len(item.mainCluster.nodes) {
			t.Errorf("clusterNames[item.mainCluster.id] != len(item.mainCluster.nodes) [%d] != [%d]", clusterNames[item.mainCluster.id], len(item.mainCluster.nodes))
		}
	}

	if float64(len(clusterNames)) > math.Ceil(float64(count/arr[0].mainCluster.clusterSize)) {
		t.Errorf("len(clusterNames) [%d] > count/arr[0].mainCluster.clusterSize [%d]", len(clusterNames), count/arr[0].mainCluster.clusterSize)
	}
}

func TestNetwork_AddNode_bigger_than_clusterCount_smaller_than_reclustering(t *testing.T) {
	maxNetworkJoinLatency = 0
	maxNetworkEventLatency = 0

	network := NewNetWork()
	count := network.entryCluster.clusterSize + network.entryCluster.criticalSize - 1

	arr := make([]*Node, count)

	for i := 0; i < count; i++ {
		n := NewNode()
		if err := network.AddNode(n); err != nil {
			t.Errorf("Network.AddNode() error = %v", err)
		}
		time.Sleep(time.Millisecond * 50)
		arr[i] = n
	}

	time.Sleep(time.Second * 1)

	// calculate the cluster count
	clusterNames := make(map[string]int)
	for _, item := range arr {
		clusterNames[item.mainCluster.id]++
	}

	// check if all the clusters has the same number of nodes attached
	for _, item := range arr {
		if clusterNames[item.mainCluster.id] != count {
			t.Errorf("clusterNames[item.mainCluster.id] != count [%d] != [%d]", clusterNames[item.mainCluster.id], count)
		}
	}

	if float64(len(clusterNames)) != math.Ceil(float64(count/arr[0].mainCluster.clusterSize)) {
		t.Errorf("len(clusterNames) [%d] > count/arr[0].mainCluster.clusterSize [%d]", len(clusterNames), count/arr[0].mainCluster.clusterSize)
	}
}

func TestNetwork_AddNode_bigger_than_clusterCount_bigger_than_reclustering(t *testing.T) {
	maxNetworkJoinLatency = 0
	maxNetworkEventLatency = 0

	network := NewNetWork()
	count := network.entryCluster.clusterSize + network.entryCluster.criticalSize

	arr := make([]*Node, count)

	for i := 0; i < count; i++ {
		n := NewNode()
		if err := network.AddNode(n); err != nil {
			t.Errorf("Network.AddNode() error = %v", err)
		}
		time.Sleep(time.Millisecond * 50)
		arr[i] = n
	}

	time.Sleep(time.Second * 1)

	// calculate the cluster count
	clusterNames := make(map[string]int)
	for _, item := range arr {
		clusterNames[item.mainCluster.id]++
	}

	// check if all the clusters has the same number of nodes attached
	for _, item := range arr {
		if clusterNames[item.mainCluster.id] < item.mainCluster.criticalSize {
			t.Errorf("clusterNames[item.mainCluster.id] != item.mainCluster.clusterSize [%d] != [%d]", clusterNames[item.mainCluster.id], item.mainCluster.clusterSize)
		}
	}

	if float64(len(clusterNames)) < math.Ceil(float64(count/arr[0].mainCluster.clusterSize)) {
		t.Errorf("len(clusterNames) [%d] > count/arr[0].mainCluster.clusterSize [%d]", len(clusterNames), count/arr[0].mainCluster.clusterSize)
	}
}
