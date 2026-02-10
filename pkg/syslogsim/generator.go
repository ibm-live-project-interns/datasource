// Package syslogsim generates and transmits simulated syslog messages over
// UDP or TCP, following the RFC 5424 syslog message format.
//
// The simulator produces realistic network and server event messages with
// randomized hostnames, application names, severity levels, and payloads.
// Messages are both sent over the network and optionally persisted to a local
// JSON file for debugging and replay.
//
// TODO: Accept a context.Context in NewSimulator/Run for graceful shutdown
// support. Replace deprecated rand.Seed with rand.New(rand.NewSource(...))
// for Go 1.20+.
package syslogsim

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// Severity represents syslog severity levels per RFC 5424.
type Severity int

const (
	SeverityInfo     Severity = iota // Informational message
	SeverityWarn                     // Warning condition
	SeverityError                    // Error condition
	SeverityCritical                 // Critical condition requiring immediate action
)

// facility is the syslog facility code (1 = user-level messages).
const facility = 1

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

// Config defines the configuration for the syslog simulation.
type Config struct {
	Host         string
	Port         int
	Protocol     string
	Interval     time.Duration
	BatchSize    int
	TotalBatches int
	FilePath     string
}

// Simulator encapsulates syslog simulation logic and state.
type Simulator struct {
	cfg  Config
	conn net.Conn
}

// NewSimulator creates a new syslog Simulator using the provided configuration.
func NewSimulator(cfg Config) (*Simulator, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	var conn net.Conn
	var err error

	if cfg.Protocol == "tcp" {
		conn, err = net.Dial("tcp", addr)
	} else {
		conn, err = net.Dial("udp", addr)
	}

	if err != nil {
		fmt.Println("warning: unable to connect:", err)
		conn, _ = net.Dial("udp", "localhost:0")
	}

	rand.Seed(time.Now().UnixNano())

	return &Simulator{
		cfg:  cfg,
		conn: conn,
	}, nil
}

// Run starts the syslog traffic simulation.
// It generates syslog messages, sends them over the network,
// and optionally persists them to a file.
func (s *Simulator) Run() error {
	defer s.conn.Close()

	batchCount := 0

	for {
		batchCount++

		for i := 0; i < s.cfg.BatchSize; i++ {
			msg, pri := s.generateSyslog()

			_, err := s.conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("warning: failed to send syslog:", err)
			}

			if err := SaveSyslogToFile(s.cfg.FilePath, msg, pri); err != nil {
				fmt.Println("warning: failed to save syslog:", err)
			}
		}

		if s.cfg.TotalBatches > 0 && batchCount >= s.cfg.TotalBatches {
			break
		}

		time.Sleep(s.cfg.Interval)
	}

	return nil
}

// RunSimulation is a convenience wrapper to quickly start a syslog simulation.
func RunSimulation(cfg Config) error {
	sim, err := NewSimulator(cfg)
	if err != nil {
		return err
	}
	return sim.Run()
}

// ---------------- Helper Methods ----------------

func (s *Simulator) generateSyslog() (string, int) {
	hostname := randomHostname()
	appName := randomAppName()
	severity := randomSeverity()

	pri := calcPriority(facility, severity)
	timestamp := time.Now().UTC().Format(time.RFC3339)
	message := randomPayload()

	syslog := fmt.Sprintf(
		"<%d>1 %s %s %s - - - %s",
		pri, timestamp, hostname, appName, message,
	)

	return syslog, pri
}

func calcPriority(facility int, severity Severity) int {
	return facility*8 + int(severity)
}

func randomSeverity() Severity {
	return Severity(rand.Intn(4))
}

func randomHostname() string {
	devTypes := []string{"router", "switch", "server"}
	return fmt.Sprintf("%s-%02d", devTypes[rand.Intn(len(devTypes))], rand.Intn(20)+1)
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
	return fmt.Sprintf("10.%d.%d.%d",
		rand.Intn(255),
		rand.Intn(255),
		rand.Intn(255),
	)
}