package config

import (
	"NebuloGo/utils"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Config structure to map the YAML fields
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
	Argon    ArgonConfig    `yaml:"argon"`
	Database DatabaseConfig `yaml:"database"`
	Debug    bool           `yaml:"debug"`
}

// ServerConfig for the server settings
type ServerConfig struct {
	Host           string   `yaml:"host"`
	Port           int      `yaml:"port"`
	ServerURL      string   `yaml:"server_url"`
	TrustedProxies []string `yaml:"trusted_proxies"`
}

// JWTConfig for the JWT secret
type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type ArgonConfig struct {
	Salt        string `yaml:"salt"`
	Parallelism uint8  `yaml:"parallelism"`
	Memory      uint32 `yaml:"memory"`
	Iterations  uint32 `yaml:"iterations"`
	HashLenght  uint32 `yaml:"hash_lenght"`
}

type DatabaseConfig struct {
	ServerURL    string `yaml:"mongodb_url"`
	DatabaseName string `yaml:"database_name"`
}

var Configuration *Config

func LoadConfig() {
	exist, err := utils.DoesExistFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if !exist {
		err := writeConfig(Config{
			Server: ServerConfig{
				Host:           "127.0.0.1",
				Port:           8080,
				ServerURL:      "http://127.0.0.1:8080",
				TrustedProxies: []string{},
			},
			JWT: JWTConfig{
				Secret: "ASuperSecretSecretlyHiddenThatNobodyKnows",
			},
			Argon: ArgonConfig{
				Salt:        "ASuperSaltForArgon2idHashFunction",
				Parallelism: 2,
				Memory:      64,
				Iterations:  2,
				HashLenght:  32,
			},
			Database: DatabaseConfig{
				ServerURL:    "mongodb://mongouser:mongopass@localhost:27017",
				DatabaseName: "nebulogo",
			},
			Debug: false,
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

	Configuration = &config
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
