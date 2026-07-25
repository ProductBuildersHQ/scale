// Package schema embeds JSON Schema files for runtime access.
package schema

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed *.schema.json
var schemaFS embed.FS

// FrameworkSchema returns the JSON Schema for scale.Framework.
func FrameworkSchema() ([]byte, error) {
	return schemaFS.ReadFile("scale-framework.schema.json")
}

// DomainSchema returns the JSON Schema for scale.Domain.
func DomainSchema() ([]byte, error) {
	return schemaFS.ReadFile("scale-domain.schema.json")
}

// AssessmentSchema returns the JSON Schema for scale.Assessment.
func AssessmentSchema() ([]byte, error) {
	return schemaFS.ReadFile("scale-assessment.schema.json")
}

// ExternalModelSchema returns the JSON Schema for scale.ExternalModel.
func ExternalModelSchema() ([]byte, error) {
	return schemaFS.ReadFile("scale-external-model.schema.json")
}

// Schema returns a schema by filename.
func Schema(name string) ([]byte, error) {
	data, err := schemaFS.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("schema %q not found: %w", name, err)
	}
	return data, nil
}

// Map returns a schema parsed as a map.
func Map(name string) (map[string]any, error) {
	data, err := Schema(name)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing schema %q: %w", name, err)
	}
	return m, nil
}
