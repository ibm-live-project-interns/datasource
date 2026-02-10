package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ibm-live-project-interns/datasource/pkg/metadatasim"
)

func main() {
	output := flag.String("output", "./data/devices-metadata.json", "Path to metadata output file")
	deviceCount := flag.Int("devices", 10, "Number of devices to generate")
	updates := flag.Int("updates", 0, "Number of metadata update cycles (0 = none)")
	updateInterval := flag.Duration("update-interval", 30*time.Second, "Time between metadata updates")

	flag.Parse()

	cfg := metadatasim.Config{
		OutputPath:     *output,
		DeviceCount:    *deviceCount,
		Updates:        *updates,
		UpdateInterval: *updateInterval,
	}

	fmt.Printf("Starting metadata publisher: devices=%d, output=%s, updates=%d, interval=%s\n",
		cfg.DeviceCount, cfg.OutputPath, cfg.Updates, cfg.UpdateInterval)

	if err := metadatasim.Run(cfg); err != nil {
		log.Fatalf("metadata publisher failed: %v", err)
	}
}
