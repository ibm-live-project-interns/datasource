package metadatasim

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Device represents metadata about a network device.
type Device struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	IP        string `json:"ip"`
	Vendor    string `json:"vendor"`
	Model     string `json:"model"`
	OS        string `json:"os"`
	Location  string `json:"location"`
	UpdatedAt string `json:"updated_at"`
}

// Config controls how metadata is generated and written.
type Config struct {
	OutputPath     string        // where to write the shared metadata file
	DeviceCount    int           // how many devices to generate
	Updates        int           // how many times to update metadata (0 = no updates)
	UpdateInterval time.Duration // delay between updates
}

// Run generates sample metadata and writes it to a common file.
// Optionally performs a few update cycles to simulate changes.
func Run(cfg Config) error {
	rand.Seed(time.Now().UnixNano())

	// Ensure the output directory exists.
	dir := filepath.Dir(cfg.OutputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Initial metadata generation.
	devices := generateDevices(cfg.DeviceCount)
	if err := writeDevices(cfg.OutputPath, devices); err != nil {
		return err
	}

	fmt.Printf("Initial metadata for %d devices written to %s\n", cfg.DeviceCount, cfg.OutputPath)

	// Optional update cycles.
	for i := 0; i < cfg.Updates; i++ {
		time.Sleep(cfg.UpdateInterval)

		updateRandomDevice(devices)
		if err := writeDevices(cfg.OutputPath, devices); err != nil {
			return err
		}
		fmt.Printf("Metadata update %d written to %s\n", i+1, cfg.OutputPath)
	}

	return nil
}

// generateDevices creates a slice of fake devices.
func generateDevices(count int) []Device {
	vendors := []string{"Cisco", "Juniper", "Arista", "Huawei"}
	models := []string{"ISR-4000", "MX480", "7050X3", "NE40E"}
	oses := []string{"IOS-XE 17.3", "JUNOS 21.1", "EOS 4.28", "VRP 8.200"}
	locations := []string{
		"DC1-Rack1",
		"DC1-Rack2",
		"DC2-Rack5",
		"Branch-Mumbai",
		"Branch-Delhi",
	}

	devices := make([]Device, 0, count)

	for i := 0; i < count; i++ {
		vendorIdx := rand.Intn(len(vendors))
		now := time.Now().UTC().Format(time.RFC3339)

		d := Device{
			ID:        fmt.Sprintf("dev-%03d", i+1),
			Hostname:  fmt.Sprintf("device-%03d", i+1),
			IP:        randomIP(),
			Vendor:    vendors[vendorIdx],
			Model:     models[vendorIdx],
			OS:        oses[vendorIdx],
			Location:  locations[rand.Intn(len(locations))],
			UpdatedAt: now,
		}

		devices = append(devices, d)
	}

	return devices
}

// writeDevices writes the devices slice as pretty JSON to the given path.
func writeDevices(path string, devices []Device) error {
	data, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// updateRandomDevice simulates a metadata change on a random device.
func updateRandomDevice(devices []Device) {
	if len(devices) == 0 {
		return
	}

	idx := rand.Intn(len(devices))
	dev := &devices[idx]

	// Randomly change OS or Location to simulate drift.
	switch rand.Intn(2) {
	case 0:
		dev.OS = dev.OS + " (patched)"
	default:
		dev.Location = dev.Location + "-Alt"
	}

	dev.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

func randomIP() string {
	return fmt.Sprintf("10.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256))
}
