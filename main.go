package main

import (
	"fmt"
	"time"

	"github.com/kr/pretty"
)

func main() {
	count := 5

	network := NewNetWork()
	arr := make([]*Node, count)

	for i := 0; i < count; i++ {
		n := NewNode()
		if err := network.AddNode(n); err != nil {
			fmt.Println("err-->", err)
		}
		arr[i] = n
		time.Sleep(time.Second)
	}

	network.SendMessageToNetwork(&NetworkEvent{
		Name: "test",
	})
	time.Sleep(time.Second * 3)
	for i := 0; i < count; i++ {
		fmt.Printf("arr[%-3d].mainCluster.id: [%s] - node count: %-3d, queue count: %-3d\n", i, arr[i].mainCluster.id, len(arr[i].mainCluster.nodes), len(arr[i].mainCluster.inQueue))
	}

	fmt.Printf("count %# v \n ", pretty.Formatter(count))
}
