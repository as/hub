package hub

import (
	"errors"
	"log"
)

var (
	ErrBadUserPacket = errors.New("bad user packet")
)

func (h *Hub) errok(err error) (ok bool) {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (h *Hub) ck(where string, err error) (ok bool) {
	if err == nil {
		return true
	}
	log.Printf("%s: %s\n", where, err)
	return false
}
