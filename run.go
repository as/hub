package hub

import (
	"log"

	"github.com/as/hub/wire"
)

func (h *Hub) Run() {
	if h.running {
		return
	}
	go h.broadcasts()
	go h.mux()
}

func (h *Hub) broadcasts() {
	for {
		select {
		case <-h.teardown:
			log.Printf("stopping server...\n")
			return
		case e := <-h.broadcast:
			for _, u := range h.Id {
				go func(u *user, e wire.Packet) {
					u.out <- e
				}(u, e)
			}
		}
	}
}

// ReplyTo sends a packet on the users allocated reply channel
func (h *Hub) ReplyTo(u *user, RcId int, p *wire.Packet) (err error) {
	log.Printf("sending reply to %d: %s\n", u.Id, p)
	p.RcId = RcId
	u.rc <- *p
	return nil
}

func (h *Hub) mux() {
	for {
		select {
		case u := <-h.enter:
			h.live[u] = true
			h.Id[u.Id] = u
			println("enter")
		case u := <-h.leave:
			delete(h.Id, u.Id)
			delete(h.live, u)
			close(u.out)
		case e := <-h.event:
			u, err := h.userOk(e.Id)
			if !h.errok(err) {
				continue
			}
			switch e.Kind {
			case 'i':
				r := h.userInsert(u, e.P)
				h.broadcast <- e
				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind, Data: r})
			case 'd':
				r := h.userDelete(u, e.Q0, e.Q1)
				h.broadcast <- e
				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind, Data: r})
			case 's':
				r := h.userSelect(u, e.Q0, e.Q1)
				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind, Data: r})
				h.broadcast <- e
			case 'w':
				h.who(u, e)
				continue
				//			case 'r':
				//				r := h.userReadAt(u, e.P, e.Q0)
				//				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind,Data: r})
			case 'R':
				r := h.userReadAt(u, e.P, e.Q0)
				r.P = e.P
				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind, Data: r})
			case '.':
				r := h.userDot(u)
				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind, Data: r})
			case 'b':
				r := wire.Data{N: int(h.Buffer.Len()), P: h.Buffer.Bytes()}
				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind, Data: r})
			case 'l':
				log.Printf("length packet")
				r := wire.Data{N: int(h.Buffer.Len())}
				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind, Data: r})
			case 'm':
				r := h.userMark(u, e.Q0)
				h.ReplyTo(u, e.RcId, &wire.Packet{Kind: e.Kind, Data: r})
			}
		case <-h.kill:
			close(h.enter)
			close(h.event)
			close(h.rst)
			// TODO: close the rest and drain the leave
			// channel for a few seconds, then tear down
			// forcefully
		}
	}
}
