package scale

// Framework is the root SCALE document: framework-level metadata and
// narratives plus the domain catalog. In a repository layout, domains are
// typically stored as separate files (domains/*.json) and assembled by
// LoadFrameworkFS, mirroring the scale-spec per-domain structure.
type Framework struct {
	SchemaURI      string           `json:"$schema,omitempty"`
	ID             string           `json:"id" jsonschema:"required"`
	Name           string           `json:"name" jsonschema:"required"`
	Description    string           `json:"description,omitempty"`
	Version        string           `json:"version,omitempty" jsonschema:"description=Version of this framework document, e.g. 0.1.0"`
	Owner          string           `json:"owner,omitempty"`
	Updated        string           `json:"updated,omitempty" jsonschema:"description=RFC 3339 date of last update"`
	Narratives     []NarrativeBlock `json:"narratives,omitempty"`
	Domains        []Domain         `json:"domains,omitempty"`
	ExternalModels []ExternalModel  `json:"externalModels,omitempty" jsonschema:"description=Codified third-party maturity models referenced by framework mappings"`
}

// ExternalModel returns the external model with the given ID, or nil.
func (f *Framework) ExternalModel(id string) *ExternalModel {
	for i := range f.ExternalModels {
		if f.ExternalModels[i].ID == id {
			return &f.ExternalModels[i]
		}
	}
	return nil
}

// Domain returns the domain with the given ID, or nil.
func (f *Framework) Domain(id string) *Domain {
	for i := range f.Domains {
		if f.Domains[i].ID == id {
			return &f.Domains[i]
		}
	}
	return nil
}

// Metric returns the metric with the given ID along with its enclosing
// domain and capability, or nils if not found. Metric IDs are unique across
// the framework (enforced by Validate).
func (f *Framework) Metric(id string) (*Domain, *Capability, *Metric) {
	for i := range f.Domains {
		d := &f.Domains[i]
		for j := range d.Capabilities {
			c := &d.Capabilities[j]
			for k := range c.Metrics {
				if c.Metrics[k].ID == id {
					return d, c, &c.Metrics[k]
				}
			}
		}
	}
	return nil, nil, nil
}
