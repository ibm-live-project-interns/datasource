package syslogsim

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// ---------- Severities ----------

type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarn
	SeverityError
	SeverityCritical
)

var severityLabels = []string{"info", "warn", "error", "critical"}

const facility = 1

// ---------- Templates ----------

var networkTemplates = []string{
	"Interface Gi0/%d changed state to %s",
	"BGP session to peer %s went %s",
	"CPU utilization on router exceeded %d%%",
}

var serverTemplates = []string{
	"Service %s restarted successfully",
	"High memory usage detected on process %s",
	"Login failure for user %s from IP %s",
}

// ---------- Config ----------

type Config struct {
	Host         string
	Port         int
	Protocol     string
	Interval     time.Duration
	BatchSize    int
	TotalBatches int
}

// ---------- Public API ----------

func RunSimulation(cfg Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	var conn net.Conn
	var err error

	if cfg.Protocol == "tcp" {
		conn, err = net.Dial("tcp", addr)
	} else {
		conn, err = net.Dial("udp", addr)
	}

	if err != nil {
		// instead of crashing, just print warning
		fmt.Println("warning: unable to connect:", err)
		// NOTE: we still continue because UDP may still send later
	}

	// IMPORTANT: if conn is nil, future writes will panic
	// so we create a dummy UDP conn that won't crash
	if conn == nil {
		conn, _ = net.Dial("udp", "localhost:0")
	}

	defer conn.Close()

	rand.Seed(time.Now().UnixNano())

	batchCount := 0

	for {
		batchCount++

		for i := 0; i < cfg.BatchSize; i++ {
			msg := randomSyslogMessage()

			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("warning: failed to send syslog message:", err)
				continue
			}
		}

		if cfg.TotalBatches > 0 && batchCount >= cfg.TotalBatches {
			break
		}

		time.Sleep(cfg.Interval)
	}

	return nil
}

// ---------- Helpers ----------

func randomSyslogMessage() string {
	hostname := randomHostname()
	appName := randomAppName()
	severity := randomSeverity()

	pri := calcPriority(facility, severity)
	timestamp := time.Now().UTC().Format(time.RFC3339)
	message := randomPayload()

	return fmt.Sprintf("<%d>1 %s %s %s - - - %s", pri, timestamp, hostname, appName, message)
}

func calcPriority(facility int, severity Severity) int {
	return facility*8 + int(severity)
}

func randomSeverity() Severity {
	return Severity(rand.Intn(len(severityLabels)))
}

func randomHostname() string {
	devTypes := []string{"router", "switch", "server"}
	idx := rand.Intn(len(devTypes))
	return fmt.Sprintf("%s-%02d", devTypes[idx], rand.Intn(20)+1)
}

func randomAppName() string {
	apps := []string{"syslogd", "snmpd", "sshd", "nginx", "kernel"}
	return apps[rand.Intn(len(apps))]
}

func randomPayload() string {
	if rand.Intn(2) == 0 {
		return randomNetworkPayload()
	}
	return randomServerPayload()
}

func randomNetworkPayload() string {
	template := networkTemplates[rand.Intn(len(networkTemplates))]

	switch template {
	case networkTemplates[0]:
		return fmt.Sprintf(template, rand.Intn(48), randomState())
	case networkTemplates[1]:
		return fmt.Sprintf(template, randomIP(), randomState())
	case networkTemplates[2]:
		return fmt.Sprintf(template, rand.Intn(40)+60)
	default:
		return "Network event"
	}
}

func randomServerPayload() string {
	template := serverTemplates[rand.Intn(len(serverTemplates))]

	switch template {
	case serverTemplates[0]:
		services := []string{"nginx", "postgres", "redis", "sshd"}
		return fmt.Sprintf(template, services[rand.Intn(len(services))])
	case serverTemplates[1]:
		processes := []string{"java", "python", "node", "go-app"}
		return fmt.Sprintf(template, processes[rand.Intn(len(processes))])
	case serverTemplates[2]:
		users := []string{"admin", "root", "guest", "developer"}
		return fmt.Sprintf(template, users[rand.Intn(len(users))], randomIP())
	default:
		return "Server event"
	}
}

func randomState() string {
	states := []string{"up", "down", "flapping"}
	return states[rand.Intn(len(states))]
}

func randomIP() string {
	return fmt.Sprintf("10.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255))
}
