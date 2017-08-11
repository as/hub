package hub

import "net"

// ListenAndServe listens on addr for wire messages and handles
// incoming connections based on the client's advertised user
// id. This mechanism is insecure and should be updated if used
// over an untrusted network.
func (h *Hub) ListenAndServe(addr string) error {
	fd, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := fd.Accept()
			if !h.errok(err) {
				continue
			}
			go h.handle(conn)
		}
	}()
	return nil
}
