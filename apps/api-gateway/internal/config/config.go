package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	Version    string     `yaml:"version" env-default:"v0.0.1"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Clients    Clients    `yaml:"clients"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"127.0.0.1:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Clients struct {
	SSO   GRPCClient `yaml:"sso"`
	Tasks GRPCClient `yaml:"tasks"`
}

type GRPCClient struct {
	Address      string        `yaml:"address" env-required:"true"`
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	RetriesCount int           `yaml:"retries_count" env-default:"3"`
	Insecure     bool          `yaml:"insecure" env-default:"false"`
}

// MustLoad retrive path to config
// if there is no config path provided, it will panic.
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

// MustLoadByPath loads the configuration from the specified path.
// If config cannot be loaded, it will panic.
func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable
// Priority: flag > env > default
// Default value: is empty string
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
