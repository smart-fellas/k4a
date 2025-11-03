package unit

import (
	"testing"

	"github.com/smart-fellas/k4a/internal/utils"
)

func TestExtractValue(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"nested": map[string]interface{}{
			"value": "nested-value",
			"deep": map[string]interface{}{
				"item": "deep-value",
			},
		},
		"count": 42,
	}

	tests := []struct {
		name    string
		path    string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "simple key",
			path:    "name",
			want:    "test",
			wantErr: false,
		},
		{
			name:    "nested key",
			path:    "nested.value",
			want:    "nested-value",
			wantErr: false,
		},
		{
			name:    "deep nested key",
			path:    "nested.deep.item",
			want:    "deep-value",
			wantErr: false,
		},
		{
			name:    "number value",
			path:    "count",
			want:    42,
			wantErr: false,
		},
		{
			name:    "nonexistent key",
			path:    "missing",
			wantErr: true,
		},
		{
			name:    "invalid nested path",
			path:    "name.value",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ExtractValue(data, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ExtractValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractString(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"nested": map[string]interface{}{
			"value": "nested-value",
		},
		"count": 42,
	}

	tests := []struct {
		name         string
		path         string
		defaultValue string
		want         string
	}{
		{
			name:         "existing string",
			path:         "name",
			defaultValue: "default",
			want:         "test",
		},
		{
			name:         "nested string",
			path:         "nested.value",
			defaultValue: "default",
			want:         "nested-value",
		},
		{
			name:         "nonexistent key returns default",
			path:         "missing",
			defaultValue: "default",
			want:         "default",
		},
		{
			name:         "number converts to string",
			path:         "count",
			defaultValue: "default",
			want:         "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.ExtractString(data, tt.path, tt.defaultValue)
			if got != tt.want {
				t.Errorf("ExtractString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractInt(t *testing.T) {
	data := map[string]interface{}{
		"count":      42,
		"countInt64": int64(100),
		"countFloat": 3.14,
		"nested": map[string]interface{}{
			"number": 99,
		},
		"notANumber": "string",
	}

	tests := []struct {
		name         string
		path         string
		defaultValue int
		want         int
	}{
		{
			name:         "int value",
			path:         "count",
			defaultValue: 0,
			want:         42,
		},
		{
			name:         "int64 value",
			path:         "countInt64",
			defaultValue: 0,
			want:         100,
		},
		{
			name:         "float value",
			path:         "countFloat",
			defaultValue: 0,
			want:         3,
		},
		{
			name:         "nested int",
			path:         "nested.number",
			defaultValue: 0,
			want:         99,
		},
		{
			name:         "nonexistent key returns default",
			path:         "missing",
			defaultValue: 123,
			want:         123,
		},
		{
			name:         "string value returns default",
			path:         "notANumber",
			defaultValue: 456,
			want:         456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.ExtractInt(data, tt.path, tt.defaultValue)
			if got != tt.want {
				t.Errorf("ExtractInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterResources(t *testing.T) {
	resources := []map[string]interface{}{
		{
			"metadata": map[string]interface{}{
				"name": "kafka-topic-1",
			},
		},
		{
			"metadata": map[string]interface{}{
				"name": "kafka-topic-2",
			},
		},
		{
			"metadata": map[string]interface{}{
				"name": "redis-queue",
			},
		},
		{
			"metadata": map[string]interface{}{
				"name": "kafka-stream",
			},
		},
	}

	tests := []struct {
		name   string
		search string
		want   int // number of expected results
	}{
		{
			name:   "empty search returns all",
			search: "",
			want:   4,
		},
		{
			name:   "filter by kafka",
			search: "kafka",
			want:   3,
		},
		{
			name:   "filter by topic",
			search: "topic",
			want:   2,
		},
		{
			name:   "filter by redis",
			search: "redis",
			want:   1,
		},
		{
			name:   "no matches",
			search: "nonexistent",
			want:   0,
		},
		{
			name:   "case insensitive",
			search: "KAFKA",
			want:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.FilterResources(resources, tt.search)
			if len(got) != tt.want {
				t.Errorf("FilterResources() returned %d resources, want %d", len(got), tt.want)
			}
		})
	}
}
