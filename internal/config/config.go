package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
}

type Context struct {
	Name    string         `yaml:"name"`
	Context ContextDetails `yaml:"context"`
}

type ContextDetails struct {
	API       string `yaml:"api"`
	UserToken string `yaml:"user-token"`
	Namespace string `yaml:"namespace"`
}

func Load() (*Config, error) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var rawConfig map[string]any
	err = yaml.Unmarshal(data, &rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Extract kafkactl configuration
	kafkactlConfig, ok := rawConfig["kafkactl"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("kafkactl configuration not found")
	}

	configData, err := yaml.Marshal(kafkactlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set current context if not set
	if cfg.CurrentContext == "" && len(cfg.Contexts) > 0 {
		cfg.CurrentContext = cfg.Contexts[0].Name
	}

	return &cfg, nil
}

func (c *Config) GetCurrentContext() (*Context, error) {
	for _, ctx := range c.Contexts {
		if ctx.Name == c.CurrentContext {
			return &ctx, nil
		}
	}
	return nil, fmt.Errorf("current context %s not found", c.CurrentContext)
}

func getConfigPath() string {
	// Check for environment variable override
	if configPath := os.Getenv("KAFKACTL_CONFIG"); configPath != "" {
		return configPath
	}

	// Default to ~/.kafkactl/config.yml
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".kafkactl", "config.yml")
}

func (c *Config) Save() error {
	configPath := getConfigPath()

	// Create the directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Wrap config in kafkactl key
	wrapped := map[string]any{
		"kafkactl": c,
	}

	data, err := yaml.Marshal(wrapped)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
