package scale

import (
	core "github.com/grokify/prism-core"
)

// Status constants for domains and capabilities.
const (
	StatusDraft      = "draft"
	StatusActive     = "active"
	StatusDeprecated = "deprecated"
)

// Capability is an engineering concept a domain is trying to improve
// (e.g. "API style conformance", "OpenTelemetry instrumentation").
// Capabilities are canonical: external frameworks reference them via
// Frameworks mappings, and PRISM assesses them via PRISM refs, but the
// capability itself belongs to SCALE's domain taxonomy.
type Capability struct {
	ID          string                  `json:"id" jsonschema:"required"`
	Name        string                  `json:"name" jsonschema:"required"`
	Description string                  `json:"description,omitempty"`
	Status      string                  `json:"status,omitempty" jsonschema:"enum=draft,enum=active,enum=deprecated"`
	Owner       string                  `json:"owner,omitempty"`
	Stages      []StageRef              `json:"stages,omitempty" jsonschema:"description=Positions of this capability on the domain's lifecycle dimension(s)"`
	Metrics     []Metric                `json:"metrics,omitempty"`
	Frameworks  []core.FrameworkMapping `json:"frameworks,omitempty"`
	PRISM       *PRISMRef               `json:"prism,omitempty"`
	Narratives  []NarrativeBlock        `json:"narratives,omitempty"`
}

// StageRef places a capability on a stage of a domain dimension.
type StageRef struct {
	Dimension string `json:"dimension" jsonschema:"required,description=DomainDimension ID"`
	Stage     string `json:"stage" jsonschema:"required,description=Stage ID within the dimension"`
}

// PRISMRef links a SCALE capability to the PRISM assessment ecosystem:
// prism-capability entries, prism-maturity domains and SLIs, and the
// evidence expected at each maturity level.
type PRISMRef struct {
	CapabilityID   string            `json:"capabilityId,omitempty" jsonschema:"description=prism-capability capability ID"`
	MaturityDomain string            `json:"maturityDomain,omitempty" jsonschema:"description=prism-maturity domain"`
	SLIIDs         []string          `json:"sliIds,omitempty"`
	LevelCriteria  map[string]string `json:"levelCriteria,omitempty" jsonschema:"description=Expected evidence per maturity level, keyed M1..M5"`
}

// ValidStatus checks if a status value is valid.
func ValidStatus(status string) bool {
	switch status {
	case StatusDraft, StatusActive, StatusDeprecated, "":
		return true
	default:
		return false
	}
}
