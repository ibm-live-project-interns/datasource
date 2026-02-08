package snmptrap

import (
	"encoding/json"
	"net"
)

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
