package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	count := 15

	network := NewNetWork()
	arr := make([]*Node, count)

	for i := 0; i < count; i++ {
		n := NewNode()
		if err := network.AddNode(n); err != nil {
			fmt.Println("err-->", err)
		}
		arr[i] = n
		fmt.Println("n-->", n)
		time.Sleep(time.Second)
	}
	go func() {
		for {
			addClusterAddNewNodeLatency()
			n := NewNode()
			if err := network.AddNode(n); err != nil {
				fmt.Println("err-->", err)
			}
			arr = append(arr, n)
		}

	}()
	for {
		addClusterRemoveNodeLatency()
		if err := arr[rand.Intn(len(arr)-1)].Close(); err != nil {
			fmt.Println("err-->", err)
		}
	}
	// select {}
}
