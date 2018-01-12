package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// all consts are in ms resolution.
const (
	maxNetworkJoinLatency  = 100
	maxNetworkEventLatency = 400 // around the world latency
)

func addLatency(max int32) {
	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(max)))
}

func addNetworkJoinLatency() {
	addLatency(maxNetworkJoinLatency)
}

func addNetworkEventLatency() {
	addLatency(maxNetworkEventLatency)
}
