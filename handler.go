package hub

import (
	"encoding/gob"
	"net"
	"github.com/as/hub/wire"
)

// handle handles each client connection
func (h *Hub) handle(conn net.Conn) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	p := wire.Packet{}
	err := dec.Decode(&p)
	if !h.errok(err) {
		return
	}

	user, err := h.usergen(&p)
	if !h.errok(err) {
		return
	}

	h.enter <- user

	go func() {
		for {
			p := wire.Packet{}
			err := dec.Decode(&p)
			if !h.errok(err) {
				close(user.disconnect)
				return
			}
			h.event <- p
		}
	}()

	for {
		var err error
		select {
		case <-user.disconnect:
			return
		case p := <-user.rc:
			n := wire.Note{Ch: wire.Reply, Packet: p}
			err = enc.Encode(n)
		case p := <-user.out:
			n := wire.Note{Ch: wire.Broadcast, Packet: p}
			err = enc.Encode(n)
		}
		if !h.errok(err) {
			return
		}
	}
}
