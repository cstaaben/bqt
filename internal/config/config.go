// Package config defines configuration for bqt.
package config

import (
	"fmt"
	"time"

	"github.com/cstaaben/bqt/internal/formatter"
	"github.com/spf13/viper"
)

// Config holds the configuration for bqt.
type Config struct {
	Verobosity      int           `json:"verobosity" mapstructure:"verobosity"`
	ProjectID       string        `json:"project_id" mapstructure:"project_id"`
	CredentialsFile string        `json:"credentials_file" mapstructure:"credentials_file"`
	BatchPriority   bool          `json:"batch_priority" mapstructure:"batch_priority"`
	Timeout         time.Duration `json:"timeout" mapstructure:"timeout"`
	Format          string        `json:"format" mapstructure:"format"`
}

// New creates a new Config, parsing configFile if provided.
func New(configFile string) (*Config, error) {
	cfg := &Config{}

	if configFile != "" {
		viper.AddConfigPath(configFile)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, cfg.validate()
}

// validate ensures the config is valid to use. Only checks that require a hard-stop should be here.
func (c *Config) validate() error {
	if _, supported := formatter.SupportedFormats[c.Format]; !supported {
		return fmt.Errorf("unsupported format: %s", c.Format)
	}

	return nil
}
