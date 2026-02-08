package simulator

import (
	"context"
	"log"
	"time"
)

// Switch simulates a network switch device that generates syslog telemetry.
//
// NOTE: Currently a placeholder — Run() logs a message every 7 seconds but
// does not produce actual syslog events. To complete the implementation,
// integrate with the pkg/syslogsim package to generate and transmit RFC 5424
// syslog messages.
type Switch struct {
	Name string
}

// Run starts the switch simulation loop, logging every 7 seconds until the
// context is cancelled.
//
// TODO: Replace the log-only stub with actual syslog event generation using
// pkg/syslogsim.RunSimulation or a dedicated syslog sender.
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
