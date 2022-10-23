package config

import (
	"flag"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	DNS                string `toml:"dns" env:"DATABASE_URL" envDefault:"postgres://postgres:1110@localhost/test?sslmode=disable"`
	Port               string `toml:"port" env:"PORT" envDefault:"1234"`
	ServerURL          string `toml:"server_url" env:"SERVER_URL" envDefault:"http://localhost:1234/"`
	ReadTimeout        int    `toml:"read_timeout" env:"READ_TIMEOUT" envDefault:"30"`
	WriteTimeout       int    `toml:"write_timeout" env:"WRITE_TIMEOUT" envDefault:"30"`
	WriteHeaderTimeout int    `toml:"write_header_timeout" env:"WRITE_HEADER_TIMEOUT" envDefault:"30"`
	LogLevel           string `toml:"log_level" env:"LOG_LEVEL" envDefault:"debug"`
}

var (
	cfg  Config
	once sync.Once
)

func NewConfig() *Config {
	// Config loaded. Once
	once.Do(func() {
		var configPath string

		// When launched with flag, specifying the path with the configuration file
		// config loaded from .toml file
		flag.StringVar(&configPath, "config-path", "", "path to file in .toml format")
		flag.Parse()

		// Loaded environment variables
		if err := env.Parse(&cfg); err != nil {
			log.Fatalf("Can't load environment variables: %s", err)
		}

		if configPath != "" {
			// Config file decoding
			_, err := toml.DecodeFile(configPath, &cfg)
			if err != nil {
				log.Fatalf("Can't load configuration file: %v", err)
			}
		}
		log.Printf("Load config successful %v", cfg)
	})
	return &cfg
}
