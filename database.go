package main

type State uint

const (
	Missing State = iota + 1
	Critical
	NotComplete
	Complete

	HasQueue
)

type Event uint

const (
	StateTransition Event = iota + 1
	NewSlotAdded
	SlotRemoved
	NeedSlot
)

type Database struct {
	data map[string]map[string]string // map[slot_id]{key:val}
}
