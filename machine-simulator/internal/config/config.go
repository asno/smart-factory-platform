package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Simulation struct {
		TickIntervalSeconds int `yaml:"tick_interval_seconds"`
	} `yaml:"simulation"`

	Machines struct {
		Conveyors []struct {
			Name string `yaml:"name"`
		} `yaml:"conveyors"`

		Ovens []struct {
			Name string `yaml:"name"`
		} `yaml:"ovens"`

		Pumps []struct {
			Name string `yaml:"name"`
		} `yaml:"pumps"`
	} `yaml:"machines"`
}

func Load(path string) (*Config, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}