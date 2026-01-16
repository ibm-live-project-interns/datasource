package simulator

import "context"

// Device represents a simulated network device
type Device interface {
	Run(ctx context.Context)
}


