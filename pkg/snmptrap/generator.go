package snmptrap

import (
	"math/rand"
	"time"
)

type Trap struct {
	Version   string            `json:"version"`
	Community string            `json:"community"`
	OID       string            `json:"oid"`
	Source    string            `json:"source"`
	Message   string            `json:"message"`
	Severity  string            `json:"severity"`
	Timestamp time.Time         `json:"timestamp"`
	Variables map[string]string `json:"variables"`
}

func RandomTrap(deviceType string, source string) Trap {
	var templates []TrapTemplate

	switch deviceType {
	case "router":
		templates = RouterTraps
	case "switch":
		templates = SwitchTraps
	case "firewall":
		templates = FirewallTraps
	default:
		templates = RouterTraps
	}

	t := templates[rand.Intn(len(templates))]

	return Trap{
		Version:   "v2c",
		Community: "public",
		OID:       t.OID,
		Source:    source,
		Message:   t.Message,
		Severity:  t.Severity,
		Timestamp: time.Now().UTC(),
		Variables: map[string]string{
			"ifIndex": "1",
			"value":   "threshold-crossed",
		},
	}
}
