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
	configPath := "config/sample.yml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	log.Printf("Loading config from %s", configPath)

	// Validate required environment variables
	// Note: We might still keep INGESTOR_CORE_URL if you use it for health checks, 
	// but KAFKA_BROKER is now implicitly used by the client.
	requiredEnvVars := []string{"INGESTOR_CORE_URL"}
	if err := config.ValidateRequiredEnvVars(requiredEnvVars); err != nil {
		log.Fatal("Environment validation failed:", err)
	}

	// 1. Initialize Ingestor HTTP Client (Optional: Keep for Health Checks if needed)
	ingestorURL := config.GetEnvRequired("INGESTOR_CORE_URL")
	ingestorClient := client.NewIngestorClient(ingestorURL)

	fmt.Printf("🔍 Checking Ingestor Core health at %s...\n", ingestorURL)
	if err := ingestorClient.HealthCheck(); err != nil {
		log.Printf("⚠️  Warning: Ingestor Core HTTP health check failed: %v\n", err)
	} else {
		fmt.Println("✅ Ingestor Core HTTP is reachable")
	}

	// 2. Initialize Kafka Producer (NEW)
	fmt.Println("🔌 Initializing Kafka Producer...")
	kafkaClient, err := client.NewKafkaProducer()
	if err != nil {
		log.Fatalf("❌ Failed to start Kafka producer: %v", err)
	}
	defer kafkaClient.Close()
	fmt.Println("✅ Kafka Producer Ready")

	// --- Process Syslog events ---
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

	fmt.Println("\n📤 Sending Syslog events to Kafka...")
	for i, raw := range rawSyslogs {
		event, err := mapper.MapSyslog(raw)
		if err != nil {
			log.Printf("❌ Syslog %d mapping failed: %v\n", i+1, err)
			continue
		}

		// Changed to use Kafka Client
		if err := kafkaClient.SendEventAsync(event); err != nil {
			log.Printf("❌ Syslog %d Kafka send failed: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ Syslog %d queued successfully\n", i+1)
		}
	}

	// --- Process SNMP events ---
	rawSNMPs := [][]byte{
		[]byte(`{
			"source": "router1",
			"oid": "1.3.6.1.2.1.1.3",
			"value": "123456",
			"severity": "CRITICAL",
			"timestamp": "2025-01-05T10:47:00Z"
		}`),
	}

	fmt.Println("\n📤 Sending SNMP events to Kafka...")
	for i, raw := range rawSNMPs {
		event, err := mapper.MapSNMP(raw)
		if err != nil {
			log.Printf("❌ SNMP %d mapping failed: %v\n", i+1, err)
			continue
		}

		if err := kafkaClient.SendEventAsync(event); err != nil {
			log.Printf("❌ SNMP %d Kafka send failed: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ SNMP %d queued successfully\n", i+1)
		}
	}

	// --- Process Metadata events ---
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

	fmt.Println("\n📤 Sending Metadata events to Kafka...")
	for i, raw := range rawMetadata {
		event, err := mapper.MapMetadata(raw)
		if err != nil {
			log.Printf("❌ Metadata %d mapping failed: %v\n", i+1, err)
			continue
		}

		if err := kafkaClient.SendEventAsync(event); err != nil {
			log.Printf("❌ Metadata %d Kafka send failed: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ Metadata %d queued successfully\n", i+1)
		}
	}

	// Flush Kafka messages before exiting (critical for short-lived programs)
	fmt.Println("\n⏳ Flushing Kafka messages...")
	kafkaClient.Close() // This calls Flush internally based on the previous snippet provided

	fmt.Println("🎉 Datasource completed processing all events")
	os.Exit(0)
}