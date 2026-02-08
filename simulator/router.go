package simulator

import (
	"context"
	"log"
	"time"
)

// Router simulates a router device that generates SNMP trap telemetry.
//
// NOTE: Currently a placeholder — Run() logs a message every 5 seconds but
// does not produce actual SNMP traps. To complete the implementation, integrate
// with the pkg/snmptrap package to call snmptrap.RandomTrap("router", r.Name)
// and send the trap via snmptrap.SendTrap.
type Router struct {
	Name string
}

// Run starts the router simulation loop, logging every 5 seconds until the
// context is cancelled.
//
// TODO: Replace the log-only stub with actual SNMP trap generation using
// pkg/snmptrap.RandomTrap and pkg/snmptrap.SendTrap.
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
