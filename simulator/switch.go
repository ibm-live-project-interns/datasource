package simulator

import (
	"context"
	"log"
	"time"
)

// Switch simulates a switch device
type Switch struct {
	Name string
}

func (s *Switch) Run(ctx context.Context) {
	ticker := time.NewTicker(7 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Switch stopped:", s.Name)
			return
		case <-ticker.C:
			log.Println("Switch sending syslog:", s.Name)
		}
	}
}
