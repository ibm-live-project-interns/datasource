package snmptrap

import (
	"fmt"
	"net"
)

// testListen is an unexported development/debugging utility that listens for
// incoming SNMP trap data on UDP port 5162 and prints it to stdout.
//
// NOTE: This function is intentionally unexported — it exists for local
// development testing only. It panics on bind errors and silently ignores
// read errors, which is acceptable for a throwaway debug tool but should not
// be used in production.
//
// TODO: Consider moving this into a cmd/snmp-trap-listener entry point with
// proper error handling and graceful shutdown support.
func testListen() {
	addr, _ := net.ResolveUDPAddr("udp", ":5162")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Listening for SNMP traps on UDP :5162")

	buf := make([]byte, 4096)

	for {
		n, remote, _ := conn.ReadFromUDP(buf)
		fmt.Printf("From %s:\n%s\n\n", remote, string(buf[:n]))
	}
}
