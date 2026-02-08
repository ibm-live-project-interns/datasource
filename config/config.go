// Package config provides YAML-based configuration loading for the datasource
// simulator. It defines device and simulator configuration structures that map
// to the config/sample.yml reference file.
//
// NOTE: This package is available but not yet integrated into the simulator/
// runtime. The simulator.Manager currently hardcodes its device list. Future
// iterations should call LoadConfig to read devices from sample.yml and pass
// them to the Manager.
package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// DeviceConfig defines a single simulated device's properties.
type DeviceConfig struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	TrapInterval   int    `yaml:"trap_interval"`
	SyslogInterval int    `yaml:"syslog_interval"`
}

// SimulatorConfig is the top-level configuration loaded from YAML.
type SimulatorConfig struct {
	Devices []DeviceConfig `yaml:"devices"`
}

// LoadConfig reads and parses a YAML configuration file at the given path.
func LoadConfig(path string) (*SimulatorConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg SimulatorConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
