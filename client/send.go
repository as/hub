package client

import (
	"encoding/gob"
//	"log"
	"net"

	"github.com/as/event"
	"github.com/as/frame"
	"github.com/as/hub/wire"
	"github.com/as/text"
)

//var debug = //log.Printf

func DialEvent(Id int, fr *frame.Frame, eventsink text.Sender, network, address string) (c *User) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil
	}
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	enc.Encode(wire.Packet{Id: Id, Kind: 'E'})
	c = &User{
		Id:        Id,
		askin:     make(chan request),
		replyc:    make(chan wire.Packet),
		broadcast: make(chan wire.Packet),
		rc:        make(chan wire.Packet),
		ask:       make(chan wire.Packet),
		knows:     make(map[int]*userinfo),
		sem:       make(chan bool, 1),
		replyfns:  make(map[int]replyfn),
		Sender:    eventsink,
		fr: fr,	// TODO: Remove this
	}
	c.knows[Id] = &userinfo{
		Id:  Id,
		Q0:  0,
		Q1:  0,
		col: &frame.Acme,
	}
	go c.asker()
	go func() {
		var Note wire.Note
		for {
			err = dec.Decode(&Note)
			if err != nil {
				break
			}
			switch Note.Ch {
			case wire.Broadcast:
				//log.Printf("Recv broadcast: %#v\n", Note)
				c.broadcast <- Note.Packet
			case wire.Reply:
				//log.Printf("Recv reply: %#v\n", Note)
				c.replyc <- Note.Packet
			}

		}
	}()
	go func() {
		for e := range c.ask {
			err = enc.Encode(e)
			if err != nil {
				//log.Printf("ask: encode error: %s\n", err)
			}
		}
	}()
	go func() {
		for {
			select {
			case e := <-c.broadcast:
				switch e.Kind {
				case 'i':
					////log.Printf("broadcast action: frameInsert")
					//c.frameInsert(e.P, e.Q0)
					c.Send(event.Insert{ID: e.Id, Q0: e.Q0, Q1: e.Q1, P: e.P})
				case 'd':
					////log.Printf("broadcast action: frameDelete")
					//c.frameDelete(e.Q0, e.Q1)
					c.Send(event.Delete{ID: e.Id, Q0: e.Q0, Q1: e.Q1})
				case 's':
					////log.Printf("broadcast action: frameSelect")
					//c.frameSelect(e.Id, e.Q0, e.Q1)
					c.Send(event.Select{ID: e.Id, Q0: e.Q0, Q1: e.Q1})
				default:
					//log.Printf("broadcast action: unknown: %s", e)
				}
				//log.Printf("broadcast action resolved")
			}
		}
	}()
	return c

}
