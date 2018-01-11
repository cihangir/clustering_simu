package main

import (
	"fmt"
)

func main() {
	n := NewNode()
	err := n.fsm.Event("setup", "hadi", "bakam")
	if err != nil {
		fmt.Println("err while transition-->", err)
	}

	fmt.Println("node :" + n.fsm.Current())

	c := NewCluster()

	fmt.Println(c.fsm.Current())

	err = c.fsm.Event("complete", "anber", "dostum")
	if err != nil {
		fmt.Println("err while transition-->", err)
	}

	fmt.Println("cluster :" + c.fsm.Current())
}
