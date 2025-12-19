package main

import (
	"fmt"
	"net"
)

func main() {
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
