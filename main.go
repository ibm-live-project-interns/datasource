package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ibm-live-project-interns/datasource/client"
	"github.com/ibm-live-project-interns/datasource/mapper"
	"github.com/ibm-live-project-interns/ingestor/shared/config"
)

func main() {
	fmt.Println("📡 Datasource starting...")

	// Validate required environment variables
	requiredEnvVars := []string{"INGESTOR_CORE_URL"}
	if err := config.ValidateRequiredEnvVars(requiredEnvVars); err != nil {
		log.Fatal("Environment validation failed:", err)
	}

	// Initialize Ingestor Client
	ingestorURL := config.GetEnvRequired("INGESTOR_CORE_URL")
	ingestorClient := client.NewIngestorClient(ingestorURL)

	// Health check
	fmt.Printf("🔍 Checking Ingestor Core health at %s...\n", ingestorURL)
	if err := ingestorClient.HealthCheck(); err != nil {
		log.Printf("⚠️  Warning: Ingestor Core health check failed: %v\n", err)
		log.Println("Continuing anyway, events will be retried...")
	} else {
		fmt.Println("✅ Ingestor Core is healthy")
	}

	// Process Syslog events
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

	fmt.Println("\n📤 Sending Syslog events...")
	for i, raw := range rawSyslogs {
		event, err := mapper.MapSyslog(raw)
		if err != nil {
			log.Printf("❌ Syslog %d mapping failed: %v\n", i+1, err)
			continue
		}

		if err := ingestorClient.SendEvent(event); err != nil {
			log.Printf("❌ Syslog %d send failed: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ Syslog %d sent successfully\n", i+1)
		}
	}

	// Process SNMP events
	rawSNMPs := [][]byte{
		[]byte(`{
			"source": "router1",
			"oid": "1.3.6.1.2.1.1.3",
			"value": "123456",
			"severity": "CRITICAL",
			"timestamp": "2025-01-05T10:47:00Z"
		}`),
	}

	fmt.Println("\n📤 Sending SNMP events...")
	for i, raw := range rawSNMPs {
		event, err := mapper.MapSNMP(raw)
		if err != nil {
			log.Printf("❌ SNMP %d mapping failed: %v\n", i+1, err)
			continue
		}

		if err := ingestorClient.SendEvent(event); err != nil {
			log.Printf("❌ SNMP %d send failed: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ SNMP %d sent successfully\n", i+1)
		}
	}

	// Process Metadata events
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

	fmt.Println("\n📤 Sending Metadata events...")
	for i, raw := range rawMetadata {
		event, err := mapper.MapMetadata(raw)
		if err != nil {
			log.Printf("❌ Metadata %d mapping failed: %v\n", i+1, err)
			continue
		}

		if err := ingestorClient.SendEvent(event); err != nil {
			log.Printf("❌ Metadata %d send failed: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ Metadata %d sent successfully\n", i+1)
		}
	}

	fmt.Println("\n🎉 Datasource completed processing all events")
	os.Exit(0)
}
