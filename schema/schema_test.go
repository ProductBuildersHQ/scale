package schema

import "testing"

func TestEmbeddedSchemas(t *testing.T) {
	names := []string{
		"scale-framework.schema.json",
		"scale-domain.schema.json",
		"scale-assessment.schema.json",
		"scale-external-model.schema.json",
	}
	for _, name := range names {
		m, err := Map(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if _, ok := m["$id"]; !ok {
			t.Errorf("%s: missing $id", name)
		}
		if m["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
			t.Errorf("%s: unexpected $schema %v", name, m["$schema"])
		}
	}
}
