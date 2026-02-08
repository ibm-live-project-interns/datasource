package snmptrap

type TrapTemplate struct {
	OID      string
	Message  string
	Severity string
}

var RouterTraps = []TrapTemplate{
	{"1.3.6.1.6.3.1.1.5.3", "Interface down", "critical"},
	{"1.3.6.1.6.3.1.1.5.4", "Interface up", "info"},
	{"1.3.6.1.4.1.9.2.1.57", "High CPU utilization", "warning"},
}

var SwitchTraps = []TrapTemplate{
	{"1.3.6.1.4.1.9.9.13.3.1.3", "Port security violation", "error"},
	{"1.3.6.1.4.1.9.9.46.2.1.1", "Spanning tree topology change", "warning"},
}

var FirewallTraps = []TrapTemplate{
	{"1.3.6.1.4.1.9.9.147.1.2", "Firewall authentication failure", "critical"},
}
