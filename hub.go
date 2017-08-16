// Package hub provides a concurrent demultiplexer for text editing operations
// from multiple clients. To connect to a hub and make concurrent edits, see hub/client
package hub

import (
	"github.com/as/hub/wire"
	"github.com/as/text"
)

// Hub is a server for multi-editing text files
type Hub struct {
	text.Buffer
	live         map[*user]bool
	Id           map[int]*user
	enter, leave chan *user
	event        chan wire.Packet
	broadcast    chan wire.Packet
	kill         chan bool
	rst          chan bool
	teardown     chan bool
	running      bool
}

// NewHub creates a new hub from the given text.Buffer. It returns
// a Hub type capable of coalescing text.Editor operations from
// multiple users concurrenly.
func NewHub(b text.Buffer) *Hub {
	return &Hub{
		Buffer:    b,
		live:      make(map[*user]bool),
		Id:        make(map[int]*user),
		enter:     make(chan *user),
		leave:     make(chan *user),
		event:     make(chan wire.Packet),
		broadcast: make(chan wire.Packet),
		kill:      make(chan bool),
		rst:       make(chan bool),
	}
}

// coherence calls the text.Coherence function on all connected users after
// an insertion or deletion occurs by one user. This extends, retracts, or
// shifts selections according to where the insertion (sign=1) or deletion
// (sign = -1) was made.
func (h *Hub) coherence(Id int, sign int, r0, r1 int64) {
	for k, u := range h.Id {
		if k == Id {
			continue
		}
		q0, q1 := h.clamp(u.Q0, u.Q1)
		u.Q0, u.Q1 = text.Coherence(sign, r0, r1, q0, q1)
	}
}
