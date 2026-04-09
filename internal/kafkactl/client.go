package kafkactl

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/smart-fellas/k4a/internal/cache"
	"github.com/smart-fellas/k4a/internal/config"
	"gopkg.in/yaml.v3"
)

type Client struct {
	config          *config.Config
	cache           *cache.Manager
	refreshInterval time.Duration
	useCache        bool
}

func NewClient(cfg *config.Config) *Client {
	cacheManager, err := cache.NewManager()
	if err != nil {
		// Log error but continue without cache
		fmt.Printf("Warning: failed to initialize cache: %v\n", err)
	}

	return &Client{
		config:          cfg,
		cache:           cacheManager,
		refreshInterval: cache.DefaultRefreshInterval,
		useCache:        cacheManager != nil,
	}
}

// SetRefreshInterval sets the cache refresh interval
func (c *Client) SetRefreshInterval(interval time.Duration) {
	c.refreshInterval = interval
}

// InvalidateCache invalidates all cached data
func (c *Client) InvalidateCache() error {
	if c.cache == nil {
		return nil
	}
	if err := c.cache.InvalidateAll(); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}
	return nil
}

// ExecuteCommand runs a kafkactl command and returns the output.
func (c *Client) ExecuteCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("kafkactl", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
	}

	return out.Bytes(), nil
}

// executeCommandWithCache runs a kafkactl command with caching support
func (c *Client) executeCommandWithCache(resourceType string, forceRefresh bool, args ...string) ([]byte, error) {
	// Check cache first if not forcing refresh and cache is available
	if !forceRefresh && c.useCache && c.cache != nil {
		cachedData, found, err := c.cache.Get(resourceType, c.refreshInterval, args...)
		if err != nil {
			// Log error but continue to fetch fresh data
			fmt.Printf("Warning: cache read error: %v\n", err)
		} else if found {
			return cachedData, nil
		}
	}

	// Build command arguments
	cmdArgs := append([]string{"get", resourceType}, args...)
	cmdArgs = append(cmdArgs, "-o", "yaml")

	// Execute command
	output, err := c.ExecuteCommand(cmdArgs...)
	if err != nil {
		return nil, err
	}

	// Save to cache
	if c.useCache && c.cache != nil {
		if cacheErr := c.cache.Set(resourceType, output, args...); cacheErr != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to cache data: %v\n", cacheErr)
		}
	}

	return output, nil
}

// GetTopics retrieves all topics.
func (c *Client) GetTopics(forceRefresh ...bool) ([]map[string]any, error) {
	refresh := false
	if len(forceRefresh) > 0 {
		refresh = forceRefresh[0]
	}

	output, err := c.executeCommandWithCache("topics", refresh)
	if err != nil {
		return nil, err
	}

	return c.parseYAMLList(output)
}

// GetSchemas retrieves all schemas.
func (c *Client) GetSchemas(forceRefresh ...bool) ([]map[string]any, error) {
	refresh := false
	if len(forceRefresh) > 0 {
		refresh = forceRefresh[0]
	}

	output, err := c.executeCommandWithCache("schemas", refresh)
	if err != nil {
		return nil, err
	}

	return c.parseYAMLList(output)
}

// GetConnectors retrieves all connectors.
func (c *Client) GetConnectors(forceRefresh ...bool) ([]map[string]any, error) {
	refresh := false
	if len(forceRefresh) > 0 {
		refresh = forceRefresh[0]
	}

	output, err := c.executeCommandWithCache("connectors", refresh)
	if err != nil {
		return nil, err
	}

	return c.parseYAMLList(output)
}

// GetConsumerGroups retrieves consumer groups for a topic.
func (c *Client) GetConsumerGroups(topic string, forceRefresh ...bool) ([]map[string]any, error) {
	refresh := false
	if len(forceRefresh) > 0 {
		refresh = forceRefresh[0]
	}

	output, err := c.executeCommandWithCache("consumer-groups", refresh, "--topic", topic)
	if err != nil {
		return nil, err
	}

	return c.parseYAMLList(output)
}

// GetResourceYAML retrieves the YAML for a specific resource.
func (c *Client) GetResourceYAML(resourceType, name string) (string, error) {
	output, err := c.ExecuteCommand("get", resourceType, name, "-o", "yaml")
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (c *Client) parseYAMLList(data []byte) ([]map[string]any, error) {
	// Split by document separator
	docs := strings.Split(string(data), "---")

	var results []map[string]any
	for _, doc := range docs {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		var item map[string]any
		if err := yaml.Unmarshal([]byte(doc), &item); err != nil {
			continue
		}

		results = append(results, item)
	}

	return results, nil
}
