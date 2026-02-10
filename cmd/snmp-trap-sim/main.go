package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/ibm-live-project-interns/datasource/pkg/snmptrap"
)

func main() {
	addr := flag.String("addr", "localhost:5162", "UDP address")
	device := flag.String("device", "router", "router|switch|firewall")
	freq := flag.Int("freq", 3, "seconds between traps")
	file := flag.String("file", "data/snmp-traps.json", "file to store traps")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Starting SNMP trap simulator (%s) → %s\n", *device, *addr)

	for {
		trap := snmptrap.RandomTrap(*device, "device-01")

		// Send over UDP
		err := snmptrap.SendTrap(*addr, trap)
		if err != nil {
			fmt.Println("failed to send trap:", err)
		} else {
			fmt.Println("sent trap:", trap.OID, "-", trap.Message)
		}

		// Save to JSON file
		err = snmptrap.SaveTrapToFile(*file, trap)
		if err != nil {
			fmt.Println("failed to save trap:", err)
		}

		time.Sleep(time.Duration(*freq) * time.Second)
	}
}
