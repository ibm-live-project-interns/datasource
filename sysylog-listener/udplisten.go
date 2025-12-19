package main

import (
	"fmt"
	"net"
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":5140")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 2048)

	fmt.Println("Listening on UDP :5140 ...")
	for {
		n, remoteAddr, _ := conn.ReadFromUDP(buf)
		fmt.Printf("From %s: %s\n", remoteAddr, string(buf[:n]))
	}
}
