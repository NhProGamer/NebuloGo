package config

import (
	"NebuloGo/utils"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	HOST string `yaml:"host"`
	PORT int    `yaml:"port"`
}

func LoadConfig() *Config {
	exist, err := utils.DoesExistFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if !exist {
		err := writeConfig(Config{
			HOST: "127.0.0.1",
			PORT: 8080,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Read YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	// Unmarshal YAML data into Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}

	return &config
}

func writeConfig(config Config) error {

	// Marshal struct into YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Write YAML data to file
	err = os.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
