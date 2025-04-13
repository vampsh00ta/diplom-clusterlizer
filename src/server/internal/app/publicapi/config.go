package publicapi

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)

type (
	// Config -.
	Config struct {
		App   App   `yaml:"app"`
		PG    PG    `yaml:"postgres"`
		Front Front `yaml:"front"`
		Kafka Kafka `yaml:"broker"`
	}

	App struct {
		Name    string `yaml:"name" `
		Version string `yaml:"version"`
		Port    string `json:"port"`
	}

	Log struct {
		Level string ` yaml:"log_level"   env:"LOG_LEVEL"`
	}

	PG struct {
		PoolMax  int    ` yaml:"pool_max"`
		Username string ` yaml:"username"`
		Password string ` yaml:"password" `
		Host     string ` yaml:"host"`
		Port     string ` yaml:"port"`
		Name     string ` yaml:"name"`
	}
	Front struct {
		Static string `yaml:"static"`
	}
)

type (
	Kafka struct {
		Producer Producer `yaml:"producer"`
		Consumer Consumer `yaml:"consumer"`
	}
	Producer struct {
		DocumentSender KafkaBase `yaml:"document_sender"`
	}
	Consumer struct {
		DocumentSaver KafkaBase `yaml:"document_saver"`
	}

	KafkaBase struct {
		URL       string ` yaml:"url"`
		Topic     string `yaml:"topic"`
		MaxBytes  int    ` yaml:"max_bytes"`
		Partition int    `yaml:"partition"`
		Group     string `yaml:"group"`
	}
)

func NewDefaultConfig(configPath string) (*Config, error) {
	var err error
	cfg := &Config{}

	err = cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
