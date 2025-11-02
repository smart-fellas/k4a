package kafkactl

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/smart-fellas/k4a/internal/config"
	"gopkg.in/yaml.v3"
)

type Client struct {
	config *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{config: cfg}
}

// ExecuteCommand runs a kafkactl command and returns the output
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

// GetTopics retrieves all topics
func (c *Client) GetTopics() ([]map[string]interface{}, error) {
	output, err := c.ExecuteCommand("get", "topics", "-o", "yaml")
	if err != nil {
		return nil, err
	}

	return c.parseYAMLList(output)
}

// GetSchemas retrieves all schemas
func (c *Client) GetSchemas() ([]map[string]interface{}, error) {
	output, err := c.ExecuteCommand("get", "schemas", "-o", "yaml")
	if err != nil {
		return nil, err
	}

	return c.parseYAMLList(output)
}

// GetConnectors retrieves all connectors
func (c *Client) GetConnectors() ([]map[string]interface{}, error) {
	output, err := c.ExecuteCommand("get", "connectors", "-o", "yaml")
	if err != nil {
		return nil, err
	}

	return c.parseYAMLList(output)
}

// GetConsumerGroups retrieves consumer groups for a topic
func (c *Client) GetConsumerGroups(topic string) ([]map[string]interface{}, error) {
	output, err := c.ExecuteCommand("get", "consumer-groups", "--topic", topic, "-o", "yaml")
	if err != nil {
		return nil, err
	}

	return c.parseYAMLList(output)
}

// GetResourceYAML retrieves the YAML for a specific resource
func (c *Client) GetResourceYAML(resourceType, name string) (string, error) {
	output, err := c.ExecuteCommand("get", resourceType, name, "-o", "yaml")
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (c *Client) parseYAMLList(data []byte) ([]map[string]interface{}, error) {
	// Split by document separator
	docs := strings.Split(string(data), "---")

	var results []map[string]interface{}
	for _, doc := range docs {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		var item map[string]interface{}
		if err := yaml.Unmarshal([]byte(doc), &item); err != nil {
			continue
		}

		results = append(results, item)
	}

	return results, nil
}
