package simulator

import (
	"context"
	"log"
	"time"
)

// Router simulates a router device
type Router struct {
	Name string
}

func (r *Router) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Router stopped:", r.Name)
			return
		case <-ticker.C:
			log.Println("Router sending SNMP trap:", r.Name)
		}
	}
}
