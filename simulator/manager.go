package simulator

import (
	"context"
	"log"
	"sync"
)

// Manager controls multiple device simulators
type Manager struct {
	devices []Device
}

func NewManager() *Manager {
	return &Manager{
		devices: []Device{
			&Router{Name: "router-1"},
			&Switch{Name: "switch-1"},
		},
	}
}

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
