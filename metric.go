package scale

import (
	core "github.com/grokify/prism-core"
)

// Metric unit constants (informational; other units are permitted).
const (
	UnitPercent = "percent"
	UnitCount   = "count"
	UnitMinutes = "minutes"
	UnitHours   = "hours"
	UnitScore   = "score"
)

// Metric is a single measurable fact, tagged with exactly one SCALE aspect.
// Metrics measure reality; frameworks interpret metrics — so a metric may
// additionally map to external frameworks (DORA, NIST CSF, ...) via
// Frameworks without belonging to them.
//
// A metric participates in aspect rollups only when it has both a Target and
// an Owner (see RollupEligible). Metrics without targets are tracked but
// excluded from the story.
type Metric struct {
	ID              string                  `json:"id" jsonschema:"required"`
	Name            string                  `json:"name" jsonschema:"required"`
	Description     string                  `json:"description,omitempty"`
	Aspect          string                  `json:"aspect" jsonschema:"required,enum=standards,enum=consumption,enum=automation,enum=leverage,enum=effectiveness"`
	ConsumptionKind string                  `json:"consumptionKind,omitempty" jsonschema:"enum=adoption,enum=conformance,description=Required when aspect is consumption; distinguishes 'using it' from 'using it correctly'"`
	Unit            string                  `json:"unit,omitempty"`
	Direction       string                  `json:"direction,omitempty" jsonschema:"enum=higher_is_better,enum=lower_is_better,description=Defaults to higher_is_better"`
	Target          *Target                 `json:"target,omitempty"`
	Maturity        *MaturityLadder         `json:"maturity,omitempty" jsonschema:"description=Optional ladder mapping observed values to PRISM maturity levels"`
	Owner           string                  `json:"owner,omitempty"`
	Source          *Source                 `json:"source,omitempty"`
	Frameworks      []core.FrameworkMapping `json:"frameworks,omitempty"`
	Tags            []string                `json:"tags,omitempty"`
}

// Target is the goal value for a metric. Targets make metrics
// story-eligible: without a target there is no variance, and without
// variance there is no narrative.
type Target struct {
	Value       float64 `json:"value" jsonschema:"required"`
	By          string  `json:"by,omitempty" jsonschema:"description=Date or period the target should be reached (e.g. 2026-Q4)"`
	Description string  `json:"description,omitempty"`
}

// Source identifies the system of record that produces a metric's
// observations, e.g. api-style-spec lint reports or OpenTelemetry inventory.
type Source struct {
	System      string `json:"system" jsonschema:"required,description=Producing system, e.g. api-style-spec"`
	Ref         string `json:"ref,omitempty" jsonschema:"description=Locator within the system, e.g. lint-report.conformance_level"`
	Description string `json:"description,omitempty"`
}

// RollupEligible reports whether the metric participates in aspect rollups.
// Eligibility requires both a target (so attainment is computable) and an
// owner (so the number is accountable).
func (m *Metric) RollupEligible() bool {
	return m.Target != nil && m.Owner != ""
}

// EffectiveDirection returns the metric's comparison direction, defaulting
// to higher-is-better.
func (m *Metric) EffectiveDirection() string {
	if m.Direction == "" {
		return core.SLIDirectionHigherIsBetter
	}
	return m.Direction
}
