package snmptrap

// TrapTemplate defines a reusable SNMP trap pattern with a standard OID,
// human-readable message, and severity level.
type TrapTemplate struct {
	OID      string // Standard SNMP OID (e.g. "1.3.6.1.6.3.1.1.5.3" for linkDown)
	Message  string // Human-readable description of the trap condition
	Severity string // One of: info, warning, error, critical
}

// RouterTraps contains SNMP trap templates for router devices.
// OIDs reference standard MIB-II and Cisco enterprise OIDs.
var RouterTraps = []TrapTemplate{
	{OID: "1.3.6.1.6.3.1.1.5.3", Message: "Interface down", Severity: "critical"},
	{OID: "1.3.6.1.6.3.1.1.5.4", Message: "Interface up", Severity: "info"},
	{OID: "1.3.6.1.4.1.9.2.1.57", Message: "High CPU utilization", Severity: "warning"},
}

// SwitchTraps contains SNMP trap templates for switch devices.
var SwitchTraps = []TrapTemplate{
	{OID: "1.3.6.1.4.1.9.9.13.3.1.3", Message: "Port security violation", Severity: "error"},
	{OID: "1.3.6.1.4.1.9.9.46.2.1.1", Message: "Spanning tree topology change", Severity: "warning"},
}

// FirewallTraps contains SNMP trap templates for firewall devices.
var FirewallTraps = []TrapTemplate{
	{OID: "1.3.6.1.4.1.9.9.147.1.2", Message: "Firewall authentication failure", Severity: "critical"},
}
