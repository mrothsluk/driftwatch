package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Provider represents a supported cloud provider.
type Provider string

const (
	ProviderAWS   Provider = "aws"
	ProviderGCP   Provider = "gcp"
	ProviderAzure Provider = "azure"
)

// Config holds the top-level driftwatch configuration.
type Config struct {
	Statefile string   `yaml:"statefile"`
	Provider  Provider `yaml:"provider"`
	Output    string   `yaml:"output"` // "text" or "json"
	Filters   Filters  `yaml:"filters"`
}

// Filters allow narrowing which resources are checked.
type Filters struct {
	ResourceTypes []string `yaml:"resource_types"`
	ExcludeIDs    []string `yaml:"exclude_ids"`
}

// Load reads and validates a YAML config file from path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	return Parse(data)
}

// Parse unmarshals YAML bytes into a Config and validates it.
func Parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Statefile == "" {
		return fmt.Errorf("config: statefile is required")
	}
	switch c.Provider {
	case ProviderAWS, ProviderGCP, ProviderAzure:
		// valid
	case "":
		return fmt.Errorf("config: provider is required")
	default:
		return fmt.Errorf("config: unsupported provider %q", c.Provider)
	}
	if c.Output == "" {
		c.Output = "text"
	}
	if c.Output != "text" && c.Output != "json" {
		return fmt.Errorf("config: output must be \"text\" or \"json\", got %q", c.Output)
	}
	return nil
}
