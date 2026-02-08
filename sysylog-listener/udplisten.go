// Package sysloglistener provides a simple UDP listener for receiving and
// displaying syslog messages during development and testing.
//
// NOTE: The directory name "sysylog-listener" contains a typo (double 'y').
// This is preserved to maintain git history attribution. A future cleanup
// could rename it to "syslog-listener" with a git mv.
//
// TODO: Add graceful shutdown via context.Context, configurable buffer size,
// and structured logging instead of raw fmt.Printf output.
package sysloglistener

import (
	"fmt"
	"net"
)

// StartUDPListener starts a UDP server on the given address and prints
// incoming syslog messages to stdout. It blocks forever once started.
//
// NOTE: This function does not support graceful shutdown — it will block
// indefinitely on ReadFromUDP. Consider accepting a context.Context and
// using conn.SetReadDeadline for interruptible reads.
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
