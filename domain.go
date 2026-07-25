package scale

// Domain is a horizontal engineering discipline (api, observability,
// security, ...). Each domain owns its capabilities and may define one or
// more ordered lifecycle dimensions that act as the domain's story spine.
type Domain struct {
	SchemaURI    string            `json:"$schema,omitempty"`
	ID           string            `json:"id" jsonschema:"required"`
	Name         string            `json:"name" jsonschema:"required"`
	Description  string            `json:"description,omitempty"`
	Status       string            `json:"status,omitempty" jsonschema:"enum=draft,enum=active,enum=deprecated"`
	Owner        string            `json:"owner,omitempty"`
	Dimensions   []DomainDimension `json:"dimensions,omitempty"`
	Capabilities []Capability      `json:"capabilities,omitempty"`
	Narratives   []NarrativeBlock  `json:"narratives,omitempty"`
}

// DomainDimension is a domain-specific classification axis, most commonly an
// ordered lifecycle (Security: Design & Planning → Development → CI/CD
// Pipeline → Offensive Validation → Operations). Stage order is the slice
// order — the arc is the order, which is what lets reports read the
// dimension left-to-right as a journey.
//
// DomainDimensions are orthogonal to SCALE aspects: the dimension tells the
// domain's own story ("how good are we at each stage of security"), while
// aspects tell the horizontal platform story ("is Standardization
// improving across domains").
type DomainDimension struct {
	ID          string  `json:"id" jsonschema:"required"`
	Name        string  `json:"name" jsonschema:"required"`
	Description string  `json:"description,omitempty"`
	Stages      []Stage `json:"stages" jsonschema:"required"`
}

// Stage is one ordered step of a DomainDimension.
type Stage struct {
	ID          string           `json:"id" jsonschema:"required"`
	Name        string           `json:"name" jsonschema:"required"`
	Description string           `json:"description,omitempty"`
	Narratives  []NarrativeBlock `json:"narratives,omitempty"`
}

// Dimension returns the dimension with the given ID, or nil.
func (d *Domain) Dimension(id string) *DomainDimension {
	for i := range d.Dimensions {
		if d.Dimensions[i].ID == id {
			return &d.Dimensions[i]
		}
	}
	return nil
}

// Capability returns the capability with the given ID, or nil.
func (d *Domain) Capability(id string) *Capability {
	for i := range d.Capabilities {
		if d.Capabilities[i].ID == id {
			return &d.Capabilities[i]
		}
	}
	return nil
}

// Stage returns the stage with the given ID, or nil.
func (dd *DomainDimension) Stage(id string) *Stage {
	for i := range dd.Stages {
		if dd.Stages[i].ID == id {
			return &dd.Stages[i]
		}
	}
	return nil
}
