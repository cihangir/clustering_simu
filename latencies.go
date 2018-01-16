package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// all consts are in ms resolution.
var (
	maxNetworkJoinLatency  = 100
	maxNetworkEventLatency = 400 // around the world latency
)

func addLatency(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(int32(max))))
}

func addNetworkJoinLatency() {
	if maxNetworkJoinLatency != 0 {
		addLatency(maxNetworkJoinLatency)
	}
}

func addNetworkEventLatency() {
	if maxNetworkEventLatency != 0 {
		addLatency(maxNetworkEventLatency)
	}
}
