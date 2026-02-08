package sysloglistener

import (
	"fmt"
	"net"
)

// StartUDPListener starts a UDP server on the given address
// and prints incoming syslog messages.
func StartUDPListener(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := make([]byte, 2048)
	fmt.Printf("Listening on UDP %s ...\n", address)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		fmt.Printf("From %s: %s\n", remoteAddr, string(buf[:n]))
	}
}
