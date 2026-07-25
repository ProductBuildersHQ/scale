package scale

// Assessment is the time-bound half of a SCALE document set: observed metric
// values for one reporting period, plus the journey and outlook narratives
// for that period. The Framework catalog is the "textbook"; an Assessment is
// the "executed notebook" — computed evidence interleaved with authored
// prose about what happened and what comes next.
type Assessment struct {
	SchemaURI    string           `json:"$schema,omitempty"`
	FrameworkID  string           `json:"frameworkId" jsonschema:"required"`
	Period       string           `json:"period" jsonschema:"required,description=Reporting period, e.g. 2026-Q3"`
	AsOf         string           `json:"asOf,omitempty" jsonschema:"description=RFC 3339 date the observations were taken"`
	Observations []Observation    `json:"observations,omitempty"`
	Narratives   []NarrativeBlock `json:"narratives,omitempty" jsonschema:"description=Journey and outlook blocks for this period; use scope to anchor each block to a framework element"`
}

// Observation is a single observed metric value. Numerator and denominator
// are optional but recommended for coverage-style metrics because "487/520
// APIs" tells a better story than "93.7%".
type Observation struct {
	MetricID    string  `json:"metricId" jsonschema:"required"`
	Value       float64 `json:"value" jsonschema:"required"`
	Numerator   *int    `json:"numerator,omitempty"`
	Denominator *int    `json:"denominator,omitempty"`
	AsOf        string  `json:"asOf,omitempty"`
	Note        string  `json:"note,omitempty"`
}

// Observation returns the observation for the given metric ID, or nil.
func (a *Assessment) Observation(metricID string) *Observation {
	for i := range a.Observations {
		if a.Observations[i].MetricID == metricID {
			return &a.Observations[i]
		}
	}
	return nil
}
