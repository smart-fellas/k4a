package utils

import (
	"fmt"
	"strings"
)

// ExtractValue safely extracts a value from nested maps.
func ExtractValue(data map[string]any, path string) (any, error) {
	keys := strings.Split(path, ".")
	current := data

	for i, key := range keys {
		if i == len(keys)-1 {
			if val, ok := current[key]; ok {
				return val, nil
			}
			return nil, fmt.Errorf("key %s not found", key)
		}

		if next, ok := current[key].(map[string]any); ok {
			current = next
		} else {
			return nil, fmt.Errorf("invalid path at %s", key)
		}
	}

	return nil, fmt.Errorf("path %s not found", path)
}

// ExtractString safely extracts a string value from nested maps.
func ExtractString(data map[string]any, path, defaultValue string) string {
	val, err := ExtractValue(data, path)
	if err != nil {
		return defaultValue
	}

	if str, ok := val.(string); ok {
		return str
	}

	return fmt.Sprintf("%v", val)
}

// ExtractInt safely extracts an int value from nested maps.
func ExtractInt(data map[string]any, path string, defaultValue int) int {
	val, err := ExtractValue(data, path)
	if err != nil {
		return defaultValue
	}

	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return defaultValue
	}
}

// FilterResources filters resources based on a search string.
func FilterResources(resources []map[string]any, search string) []map[string]any {
	if search == "" {
		return resources
	}

	search = strings.ToLower(search)
	var filtered []map[string]any

	for _, resource := range resources {
		name := ExtractString(resource, "metadata.name", "")
		if strings.Contains(strings.ToLower(name), search) {
			filtered = append(filtered, resource)
		}
	}

	return filtered
}
