package integration

import (
	"os"
	"os/exec"
	"testing"

	"github.com/smart-fellas/k4a/internal/config"
	"github.com/smart-fellas/k4a/internal/kafkactl"
)

// TestKafkactlInstalled checks if kafkactl is available
func TestKafkactlInstalled(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cmd := exec.Command("kafkactl", "version")
	err := cmd.Run()
	if err != nil {
		t.Skip("kafkactl not installed, skipping integration tests")
	}
}

// TestKafkactlClient_GetTopics tests fetching topics from a real Kafka cluster
func TestKafkactlClient_GetTopics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if kafkactl is installed
	cmd := exec.Command("kafkactl", "version")
	if err := cmd.Run(); err != nil {
		t.Skip("kafkactl not installed, skipping integration test")
	}

	// Check if config exists
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("No kafkactl config found: %v", err)
	}

	client := kafkactl.NewClient(cfg)
	topics, err := client.GetTopics()
	if err != nil {
		t.Logf("Failed to get topics (this may be expected if no Kafka is running): %v", err)
		return
	}

	t.Logf("Successfully retrieved %d topics", len(topics))
}

// TestKafkactlClient_GetSchemas tests fetching schemas from Schema Registry
func TestKafkactlClient_GetSchemas(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if kafkactl is installed
	cmd := exec.Command("kafkactl", "version")
	if err := cmd.Run(); err != nil {
		t.Skip("kafkactl not installed, skipping integration test")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Skipf("No kafkactl config found: %v", err)
	}

	client := kafkactl.NewClient(cfg)
	schemas, err := client.GetSchemas()
	if err != nil {
		t.Logf("Failed to get schemas (this may be expected if Schema Registry is not configured): %v", err)
		return
	}

	t.Logf("Successfully retrieved %d schemas", len(schemas))
}

// TestConfig_LoadFromFile tests loading actual config file
func TestConfig_LoadFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Try to load the actual config
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("No config file found (expected for CI): %v", err)
	}

	if cfg.CurrentContext == "" {
		t.Error("Config loaded but CurrentContext is empty")
	}

	if len(cfg.Contexts) == 0 {
		t.Error("Config loaded but no contexts found")
	}

	t.Logf("Config loaded successfully with context: %s", cfg.CurrentContext)
}

// TestKafkactlClient_ExecuteCommand tests raw command execution
func TestKafkactlClient_ExecuteCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if kafkactl is installed
	cmd := exec.Command("kafkactl", "version")
	if err := cmd.Run(); err != nil {
		t.Skip("kafkactl not installed, skipping integration test")
	}

	// Create a minimal config for testing
	cfg := &config.Config{
		CurrentContext: "test",
		Contexts: []config.Context{
			{
				Name: "test",
				Context: config.ContextDetails{
					API: os.Getenv("KAFKA_API_URL"),
				},
			},
		},
	}

	client := kafkactl.NewClient(cfg)
	output, err := client.ExecuteCommand("version")
	if err != nil {
		t.Fatalf("ExecuteCommand(version) failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("ExecuteCommand(version) returned empty output")
	}

	t.Logf("kafkactl version output: %s", string(output))
}
