package hub

import (
	"fmt"

	"github.com/as/frame"
	"github.com/as/hub/wire"
)

type user struct {
	Id         int
	Q0, Q1, s  int64
	col        *frame.Color
	out        chan wire.Packet
	rc         chan wire.Packet
	disconnect chan bool
}

// usergen allocates a new user from an initilization packet. The mechanism for
// doing so is currently insecure and unsuitable for use on public networks.
func (h *Hub) usergen(p *wire.Packet) (*user, error) {
	if p.Kind != 'E' {
		return nil, ErrBadUserPacket
	}
	return &user{
		Id:         p.Id,
		col:        nil,
		out:        make(chan wire.Packet),
		rc:         make(chan wire.Packet),
		disconnect: make(chan bool),
	}, nil
}

func (h *Hub) who(u *user, e wire.Packet) {
	p := wire.Packet{}
	for _, u := range h.Id {
		p.Data = wire.Data{Q0: u.Q0, Q1: u.Q1, N: u.Id}
		u.rc <- p
	}
	p.N = -1
	u.rc <- p
}

//
// Server side requests with operations on the underlying text.Buffer
//

func (h *Hub) userDot(u *user) wire.Data {
	Q0, Q1 := h.clamp(u.Q0, u.Q1)
	return wire.Data{Q0: Q0, Q1: Q1}
}
func (h *Hub) userInsert(u *user, p []byte) wire.Data {
	n := h.Buffer.Insert(p, u.Q0)
	Q0, Q1 := u.Q0, u.Q1
	h.coherence(u.Id, 1, Q0, Q1)

	return wire.Data{N: int(n)}
}

func (h *Hub) userReadAt(u *user, p []byte, at int64) wire.Data {
	q := h.Buffer.Bytes()
	if at > int64(len(q)) || at < 0 {
		return wire.Data{N: 0, Err: "read out of bounds"}
	}
	q = q[at:]
	n := copy(p, q[:min(int64(len(q)), int64(len(p)))])
	return wire.Data{N: n}
}

func (h *Hub) userDelete(u *user, q0, q1 int64) wire.Data {
	n := h.Buffer.Delete(u.Q0, u.Q1)
	u.Q0, u.Q1 = q0, q0
	h.coherence(u.Id, -1, u.Q0, u.Q1)
	return wire.Data{N: int(n), Q0: u.Q0, Q1: u.Q1}
}

func (h *Hub) userSelect(u *user, Q0, Q1 int64) wire.Data {
	u.Q0, u.Q1 = Q0, Q1
	return wire.Data{Q0: Q0, Q1: Q1}
}

func (h *Hub) userMark(u *user, Q0 int64) wire.Data {
	return h.userSelect(u, Q0, Q0)
}

func (h *Hub) userOk(Id int) (u *user, err error) {
	u, ok := h.Id[Id]
	if !ok {
		return nil, fmt.Errorf("hub: user #%d not connected\n", Id)
	}
	if u == nil {
		return nil, fmt.Errorf("hub: user nil")
	}
	return u, nil
}

/*
func (h *Hub) NL(u *user, q0, q1 int64, n int, p []byte) wire.Data{
	var reverse bool
	if q1 < q0{
		reverse = true
		q0,q1 = q1,q0
	}
	if q0 < 0{
		return wire.Data{Err: "q0 < 0"}
	}
	if q1 >= len(h.Buffer.Len()){
		return wire.Data{Err: "q0 >= 0"}
	}
	if reverse{
		r = rev.NewReader(h.Buffer.Bytes())
	} else {
		r = bytes.NewReader(h.Buffer.Bytes())
	}
	bytes.Index(s, p []byte)
}
*/
