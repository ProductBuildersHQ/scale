// Command schemagen generates JSON Schema from Go types.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	scale "github.com/ProductBuildersHQ/scale"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	schemaDir := "schema"
	if len(os.Args) > 1 {
		schemaDir = os.Args[1]
	}

	//nolint:gosec // G703: Path from CLI arg
	if err := os.MkdirAll(schemaDir, 0o755); err != nil {
		return fmt.Errorf("creating schema directory: %w", err)
	}

	schemas := []struct {
		name     string
		typ      any
		schemaID string
	}{
		{
			name:     "scale-framework.schema.json",
			typ:      &scale.Framework{},
			schemaID: "https://productbuildershq.com/schema/scale/v0/scale-framework.schema.json",
		},
		{
			name:     "scale-domain.schema.json",
			typ:      &scale.Domain{},
			schemaID: "https://productbuildershq.com/schema/scale/v0/scale-domain.schema.json",
		},
		{
			name:     "scale-assessment.schema.json",
			typ:      &scale.Assessment{},
			schemaID: "https://productbuildershq.com/schema/scale/v0/scale-assessment.schema.json",
		},
		{
			name:     "scale-external-model.schema.json",
			typ:      &scale.ExternalModel{},
			schemaID: "https://productbuildershq.com/schema/scale/v0/scale-external-model.schema.json",
		},
	}

	for _, s := range schemas {
		if err := generateSchema(schemaDir, s.name, s.typ, s.schemaID); err != nil {
			return fmt.Errorf("generating %s: %w", s.name, err)
		}
		fmt.Printf("Generated %s\n", filepath.Join(schemaDir, s.name))
	}

	return nil
}

func generateSchema(dir, filename string, typ any, schemaID string) error {
	r := &jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,
		ExpandedStruct:             true,
	}

	schema := r.Reflect(typ)
	schema.ID = jsonschema.ID(schemaID)
	schema.Version = "https://json-schema.org/draft/2020-12/schema"

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling schema: %w", err)
	}

	path := filepath.Join(dir, filename)
	//nolint:gosec // G703: Path from CLI arg
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing schema: %w", err)
	}

	return nil
}
