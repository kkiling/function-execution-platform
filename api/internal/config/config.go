package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const (
	Namespace = "function_execution_platform_api"
)

type Config struct {
	MachineID              uint16                 `yaml:"machine_id"`
	Server                 Server                 `yaml:"server"`
	Mongo                  MongodbCfg             `yaml:"mongo"`
	Git                    GitConfig              `yaml:"git"`
	DefaultContainerParams DefaultContainerParams `yaml:"default_container_params"`
}

type DefaultContainerParams struct {
	MemoryLimitMb       int     `yaml:"memory_limit_mb" default:"128"`
	MemoryReservationMb int     `yaml:"memory_reservation_mb" default:"64"`
	DiskSizeMb          int     `yaml:"disk_size_mb" default:"64"`
	CPULimit            float32 `yaml:"cpu_limit" default:"0.5"`
	CPUReservation      float32 `yaml:"cpu_reservation" default:"0.2"`
	TimeoutSec          int     `yaml:"timeout_sec" default:"30"`
}

type GitConfig struct {
	PrivateKey string `yaml:"private_key"`
}

type Server struct {
	WebPort string `yaml:"web_port"`
	APIKey  string `yaml:"api_key"`
}

type MongodbCfg struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database" default:"function_execution_platform"`
}

func LoadConfig(configFile string) (*Config, error) {
	var config Config
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failure upload yaml file. err %w", err)
	}

	err = defaults.Set(&config)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &config, nil
}
