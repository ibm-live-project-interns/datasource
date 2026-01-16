package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type DeviceConfig struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	TrapInterval   int    `yaml:"trap_interval"`
	SyslogInterval int    `yaml:"syslog_interval"`
}

type SimulatorConfig struct {
	Devices []DeviceConfig `yaml:"devices"`
}

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
