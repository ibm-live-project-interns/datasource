package simulator

import (
	"context"
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
	var wg sync.WaitGroup

	for _, d := range m.devices {
		wg.Add(1)
		go func(dev Device) {
			defer wg.Done()
			dev.Run(ctx)
		}(d)
	}

	wg.Wait()
}
