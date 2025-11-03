package models

import "time"

// Resource is the base interface for all Kafka resources.
type Resource interface {
	GetName() string
	GetKind() string
	GetNamespace() string
}

// BaseResource contains common fields for all resources.
type BaseResource struct {
	APIVersion string           `yaml:"apiVersion" json:"apiVersion"`
	Kind       string           `yaml:"kind" json:"kind"`
	Metadata   ResourceMetadata `yaml:"metadata" json:"metadata"`
}

// ResourceMetadata contains metadata for resources.
type ResourceMetadata struct {
	Name        string            `yaml:"name" json:"name"`
	Namespace   string            `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	CreatedAt   *time.Time        `yaml:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time        `yaml:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

func (r *BaseResource) GetName() string {
	return r.Metadata.Name
}

func (r *BaseResource) GetKind() string {
	return r.Kind
}

func (r *BaseResource) GetNamespace() string {
	return r.Metadata.Namespace
}
