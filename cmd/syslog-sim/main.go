package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"datasource/pkg/syslogsim"
)

func main() {
	host := flag.String("host", "localhost", "Ingestor host")
	port := flag.Int("port", 5140, "Ingestor port")
	protocol := flag.String("protocol", "udp", `"udp" or "tcp"`)
	interval := flag.Duration("interval", 2*time.Second, "Interval between batches")
	batchSize := flag.Int("batch", 5, "Number of messages per batch")
	totalBatches := flag.Int("batches", 3, "Total batches to send (0 = infinite)")

	flag.Parse()

	cfg := syslogsim.Config{
		Host:         *host,
		Port:         *port,
		Protocol:     *protocol,
		Interval:     *interval,
		BatchSize:    *batchSize,
		TotalBatches: *totalBatches,
	}

	fmt.Printf("Starting syslog simulation to %s:%d over %s\n",
		cfg.Host, cfg.Port, cfg.Protocol)

	if err := syslogsim.RunSimulation(cfg); err != nil {
		log.Fatalf("simulation failed: %v", err)
	}
}
