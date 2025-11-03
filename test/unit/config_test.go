package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smart-fellas/k4a/internal/config"
)

func TestConfig_GetCurrentContext(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
		wantCtx string
	}{
		{
			name: "get existing context",
			config: &config.Config{
				CurrentContext: "dev",
				Contexts: []config.Context{
					{
						Name: "dev",
						Context: config.ContextDetails{
							API:       "http://localhost:8080",
							UserToken: "token123",
							Namespace: "default",
						},
					},
				},
			},
			wantErr: false,
			wantCtx: "dev",
		},
		{
			name: "context not found",
			config: &config.Config{
				CurrentContext: "nonexistent",
				Contexts: []config.Context{
					{
						Name: "dev",
						Context: config.ContextDetails{
							API: "http://localhost:8080",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty contexts",
			config: &config.Config{
				CurrentContext: "dev",
				Contexts:       []config.Context{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := tt.config.GetCurrentContext()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ctx.Name != tt.wantCtx {
				t.Errorf("GetCurrentContext() context name = %v, want %v", ctx.Name, tt.wantCtx)
			}
		})
	}
}

func TestConfig_Save(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	// Set environment variable to use temp config
	oldEnv := os.Getenv("KAFKACTL_CONFIG")
	os.Setenv("KAFKACTL_CONFIG", configPath)
	defer os.Setenv("KAFKACTL_CONFIG", oldEnv)

	cfg := &config.Config{
		CurrentContext: "test",
		Contexts: []config.Context{
			{
				Name: "test",
				Context: config.ContextDetails{
					API:       "http://localhost:8080",
					UserToken: "test-token",
					Namespace: "default",
				},
			},
		},
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created at %s", configPath)
	}

	// Verify file contents
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Config file is empty")
	}
}
