package simulator

import (
	"context"
	"log"
	"sync"
)

// Manager orchestrates multiple concurrent Device simulators.
// It launches each device in its own goroutine and waits for all to complete.
//
// NOTE: The device list is currently hardcoded in NewManager. Future iterations
// should accept a configuration (e.g. from YAML or CLI flags) to define which
// devices to simulate and their parameters.
type Manager struct {
	devices []Device
}

// NewManager creates a Manager with a default set of simulated devices.
//
// TODO: Accept a device configuration parameter instead of hardcoding the
// device list. This would allow dynamic device counts and types.
func NewManager() *Manager {
	return &Manager{
		devices: []Device{
			&Router{Name: "router-1"},
			&Switch{Name: "switch-1"},
		},
	}
}

// Start launches all device simulators concurrently and blocks until all
// have stopped (either via context cancellation or completion).
//
// NOTE: Panics if no devices are configured. Consider returning an error
// instead to allow callers to handle this gracefully.
func (m *Manager) Start(ctx context.Context) {
	if len(m.devices) == 0 {
		panic("no devices configured for simulation")
	}

	log.Printf("Starting simulator with %d devices", len(m.devices))

	var wg sync.WaitGroup

	for _, d := range m.devices {
		wg.Add(1)
		log.Printf("Launching device simulator: %T", d)

		go func(dev Device) {
			defer wg.Done()
			dev.Run(ctx)
		}(d)
	}

	wg.Wait()
	log.Println("All device simulators stopped")
}
