package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mailhog/backends/config"
)

// TODO: make TLSConfig and PolicySet 'ref'able

// DefaultConfig provides a default (but relatively useless) configuration
func DefaultConfig() *Config {
	return &Config{
		Backends: map[string]config.BackendConfig{
			"local_delivery": config.BackendConfig{
				Type: "local",
				Data: map[string]interface{}{},
			},
			"local_mailbox": config.BackendConfig{
				Type: "local",
				Data: map[string]interface{}{},
			},
			"local_resolver": config.BackendConfig{
				Type: "local",
				Data: map[string]interface{}{},
			},
		},
		Delivery: &config.BackendConfig{
			Ref: "local_delivery",
		},
		Mailbox: &config.BackendConfig{
			Ref: "local_mailbox",
		},
		Resolver: &config.BackendConfig{
			Ref: "local_resolver",
		},
	}
}

// Config defines the top-level application configuration
type Config struct {
	relPath string

	Backends map[string]config.BackendConfig `json:",omitempty"`
	Delivery *config.BackendConfig           `json:",omitempty"`
	Mailbox  *config.BackendConfig           `json:",omitempty"`
	Resolver *config.BackendConfig           `json:",omitempty"`
}

// RelPath returns the path to the configuration file directory,
// used when loading external files using relative paths
func (c Config) RelPath() string {
	return c.relPath
}

var cfg = DefaultConfig()

var configFile string

// Configure returns the configuration
func Configure() *Config {
	if len(configFile) > 0 {
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Error reading %s: %s", configFile, err)
			os.Exit(1)
		}
		switch {
		case strings.HasSuffix(configFile, ".json"):
			err = json.Unmarshal(b, &cfg)
			if err != nil {
				fmt.Printf("Error parsing JSON in %s: %s", configFile, err)
				os.Exit(3)
			}
		default:
			fmt.Printf("Unsupported file type: %s\n", configFile)
			os.Exit(2)
		}

		cfg.relPath = filepath.Dir(configFile)
	}

	b, _ := json.MarshalIndent(&cfg, "", "  ")
	fmt.Println(string(b))

	return cfg
}

// RegisterFlags registers command line options
func RegisterFlags() {
	flag.StringVar(&configFile, "config-file", "", "Path to configuration file")
}
