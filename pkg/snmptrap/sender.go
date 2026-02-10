package snmptrap

import (
	"encoding/json"
	"net"
)

// SendTrap sends a JSON-encoded trap to the specified UDP address.
//
// NOTE: A new UDP connection is created for each call. For high-throughput
// scenarios, consider reusing a single connection or implementing a connection
// pool. Additionally, there is no send timeout configured — production code
// should set a write deadline via conn.SetWriteDeadline.
func SendTrap(addr string, trap Trap) error {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(trap)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
