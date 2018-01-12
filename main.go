package main

import (
	"fmt"
	"time"
)

func main() {
	network := NewNetWork()
	// node := NewNode(network)

	if err := network.AddNode(NewNode()); err != nil {
		fmt.Println("err-->", err)
	}
	if err := network.AddNode(NewNode()); err != nil {
		fmt.Println("err-->", err)
	}
	if err := network.AddNode(NewNode()); err != nil {
		fmt.Println("err-->", err)
	}
	if err := network.AddNode(NewNode()); err != nil {
		fmt.Println("err-->", err)
	}
	network.SendMessageToNetwork(&NetworkEvent{
		Name: "test",
	})
	time.Sleep(time.Second)
}
