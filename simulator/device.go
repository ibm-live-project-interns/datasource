// Package simulator provides a framework for simulating network devices that
// generate telemetry data (SNMP traps, syslog events, metadata updates).
//
// The package defines a Device interface and concrete implementations for
// routers and switches, managed by a central Manager that runs them
// concurrently. Each device type is intended to generate its specific
// telemetry type (Router → SNMP traps, Switch → syslog events).
//
// NOTE: The current Router and Switch implementations are placeholder stubs
// that log messages without generating actual telemetry. Future iterations
// should integrate with the pkg/snmptrap and pkg/syslogsim packages to
// produce real simulated events.
package simulator

import "context"

// Device represents a simulated network device capable of generating telemetry.
// Implementations should respect context cancellation for graceful shutdown.
//
// TODO: Consider adding an error return to Run() so the Manager can detect
// and handle device failures. Also consider adding Name() and Type() methods
// for better observability and logging.
type Device interface {
	Run(ctx context.Context)
}


