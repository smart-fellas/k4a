package unit

import (
	"testing"
	"time"

	"github.com/smart-fellas/k4a/pkg/models"
)

func TestBaseResource_GetName(t *testing.T) {
	resource := models.BaseResource{
		Metadata: models.ResourceMetadata{
			Name: "test-resource",
		},
	}

	got := resource.GetName()
	want := "test-resource"

	if got != want {
		t.Errorf("GetName() = %v, want %v", got, want)
	}
}

func TestBaseResource_GetKind(t *testing.T) {
	resource := models.BaseResource{
		Kind: "Topic",
	}

	got := resource.GetKind()
	want := "Topic"

	if got != want {
		t.Errorf("GetKind() = %v, want %v", got, want)
	}
}

func TestBaseResource_GetNamespace(t *testing.T) {
	resource := models.BaseResource{
		Metadata: models.ResourceMetadata{
			Namespace: "production",
		},
	}

	got := resource.GetNamespace()
	want := "production"

	if got != want {
		t.Errorf("GetNamespace() = %v, want %v", got, want)
	}
}

func TestResourceMetadata_Fields(t *testing.T) {
	now := time.Now()
	metadata := models.ResourceMetadata{
		Name:      "test-topic",
		Namespace: "default",
		Labels: map[string]string{
			"env":  "dev",
			"team": "platform",
		},
		Annotations: map[string]string{
			"description": "test topic",
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "name is set correctly",
			test: func(t *testing.T) {
				if metadata.Name != "test-topic" {
					t.Errorf("Name = %v, want %v", metadata.Name, "test-topic")
				}
			},
		},
		{
			name: "namespace is set correctly",
			test: func(t *testing.T) {
				if metadata.Namespace != "default" {
					t.Errorf("Namespace = %v, want %v", metadata.Namespace, "default")
				}
			},
		},
		{
			name: "labels are set correctly",
			test: func(t *testing.T) {
				if len(metadata.Labels) != 2 {
					t.Errorf("Labels length = %v, want %v", len(metadata.Labels), 2)
				}
				if metadata.Labels["env"] != "dev" {
					t.Errorf("Labels[env] = %v, want %v", metadata.Labels["env"], "dev")
				}
			},
		},
		{
			name: "annotations are set correctly",
			test: func(t *testing.T) {
				if metadata.Annotations["description"] != "test topic" {
					t.Errorf("Annotations[description] = %v, want %v", metadata.Annotations["description"], "test topic")
				}
			},
		},
		{
			name: "timestamps are set correctly",
			test: func(t *testing.T) {
				if metadata.CreatedAt == nil {
					t.Error("CreatedAt should not be nil")
				}
				if metadata.UpdatedAt == nil {
					t.Error("UpdatedAt should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestBaseResource_FullResource(t *testing.T) {
	now := time.Now()
	resource := models.BaseResource{
		APIVersion: "kafka.michelin.io/v1",
		Kind:       "Topic",
		Metadata: models.ResourceMetadata{
			Name:      "my-topic",
			Namespace: "production",
			Labels: map[string]string{
				"app": "myapp",
			},
			Annotations: map[string]string{
				"owner": "platform-team",
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	if resource.GetName() != "my-topic" {
		t.Errorf("GetName() = %v, want %v", resource.GetName(), "my-topic")
	}

	if resource.GetKind() != "Topic" {
		t.Errorf("GetKind() = %v, want %v", resource.GetKind(), "Topic")
	}

	if resource.GetNamespace() != "production" {
		t.Errorf("GetNamespace() = %v, want %v", resource.GetNamespace(), "production")
	}

	if resource.APIVersion != "kafka.michelin.io/v1" {
		t.Errorf("APIVersion = %v, want %v", resource.APIVersion, "kafka.michelin.io/v1")
	}
}
