package main

import (
	"fmt"
	"log"

	"github.com/aishwaryagilhotra/datasource/db"
	"github.com/aishwaryagilhotra/datasource/mapper"
)

func main() {
	fmt.Println("📡 Datasource starting...")

	// Initialize DB
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("DB init failed:", err)
	}
	defer database.Close()

	rawSyslogs := [][]byte{
		[]byte(`{
			"host": "server1",
			"severity": "ERROR",
			"message": "Disk space low",
			"timestamp": "2025-01-05T10:45:12Z"
		}`),
		[]byte(`{
			"host": "server2",
			"severity": "WARN",
			"message": "High memory usage",
			"timestamp": "2025-01-05T10:46:12Z"
		}`),
	}

	for _, raw := range rawSyslogs {
		event, err := mapper.MapSyslog(raw)
		if err != nil {
			log.Println("Syslog mapping failed:", err)
			continue
		}

		if err := db.InsertEvent(database, event); err != nil {
			log.Println("Syslog insert failed:", err)
		}
	}

	rawSNMPs := [][]byte{
		[]byte(`{
			"source": "router1",
			"oid": "1.3.6.1.2.1.1.3",
			"value": "123456",
			"severity": "CRITICAL",
			"timestamp": "2025-01-05T10:47:00Z"
		}`),
	}

	for _, raw := range rawSNMPs {
		event, err := mapper.MapSNMP(raw)
		if err != nil {
			log.Println("SNMP mapping failed:", err)
			continue
		}

		if err := db.InsertEvent(database, event); err != nil {
			log.Println("SNMP insert failed:", err)
		}
	}

	rawMetadata := [][]byte{
		[]byte(`{
			"entity": "service-auth",
			"data": {
				"version": "1.2.3",
				"region": "ap-south"
			},
			"timestamp": "2025-01-05T10:48:00Z"
		}`),
	}

	for _, raw := range rawMetadata {
		event, err := mapper.MapMetadata(raw)
		if err != nil {
			log.Println("Metadata mapping failed:", err)
			continue
		}

		if err := db.InsertEvent(database, event); err != nil {
			log.Println("Metadata insert failed:", err)
		}
	}

	fmt.Println("Syslog + SNMP + Metadata events inserted successfully")
}
