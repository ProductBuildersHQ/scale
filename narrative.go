package scale

// Narrative scope type constants, used when a NarrativeBlock is declared
// outside the element it describes (for example, in an Assessment).
const (
	ScopeFramework  = "framework"
	ScopeDomain     = "domain"
	ScopeDimension  = "dimension"
	ScopeStage      = "stage"
	ScopeCapability = "capability"
	ScopeMetric     = "metric"
)

// NarrativeBlock is authored prose attached to a framework element — the
// "literate" half of a SCALE document, analogous to markdown cells in a
// notebook. The three kinds have different lifetimes and different homes:
//
//   - thesis: why this element matters. Timeless; belongs in the framework
//     catalog. Requires an owner and should carry a reviewBy date so stale
//     prose is visibly stale.
//   - journey: where we were and what changed. Time-bound; belongs in an
//     Assessment for a specific period, anchored by asOf.
//   - outlook: where we are going and why. Forward-looking; references
//     prism-roadmap initiative IDs so reports can compare expected movement
//     against actuals in the next cycle.
type NarrativeBlock struct {
	ID          string          `json:"id" jsonschema:"required"`
	Kind        string          `json:"kind" jsonschema:"required,enum=thesis,enum=journey,enum=outlook"`
	Scope       *NarrativeScope `json:"scope,omitempty"`
	Title       string          `json:"title,omitempty"`
	Body        string          `json:"body" jsonschema:"required,description=Markdown narrative text"`
	Initiatives []string        `json:"initiatives,omitempty" jsonschema:"description=prism-roadmap initiative IDs this narrative references"`
	Owner       string          `json:"owner,omitempty"`
	AsOf        string          `json:"asOf,omitempty" jsonschema:"description=Period anchor for journey/outlook blocks (RFC 3339 date or period like 2026-Q3)"`
	ReviewBy    string          `json:"reviewBy,omitempty" jsonschema:"description=Date by which a thesis block should be re-reviewed"`
}

// NarrativeScope identifies the framework element a narrative describes when
// the block is not nested inside that element.
type NarrativeScope struct {
	Type string `json:"type" jsonschema:"required,enum=framework,enum=domain,enum=dimension,enum=stage,enum=capability,enum=metric"`
	Ref  string `json:"ref" jsonschema:"required,description=ID of the referenced element"`
}

// ValidScopeType checks if a narrative scope type is valid.
func ValidScopeType(scopeType string) bool {
	switch scopeType {
	case ScopeFramework, ScopeDomain, ScopeDimension, ScopeStage, ScopeCapability, ScopeMetric:
		return true
	default:
		return false
	}
}
