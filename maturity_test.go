package scale

import "testing"

// coverageLadder is the recommended shape: "not tracked" is a status (N/A),
// so all five levels carry attainment signal. L1 requires no other team
// (measurement reports a number, even 0); L2 ">0%" proves the pipeline
// works with at least one real reporter.
func coverageLadder() *MaturityLadder {
	t0, t50, t80, t100 := 0.0, 50.0, 80.0, 100.0
	return &MaturityLadder{Levels: []MaturityRung{
		{Level: 1, Name: "Tracked"},
		{Level: 2, Name: ">0%", Threshold: &t0, Exclusive: true},
		{Level: 3, Name: ">=50%", Threshold: &t50},
		{Level: 4, Name: ">=80%", Threshold: &t80},
		{Level: 5, Name: "100%", Threshold: &t100},
	}}
}

func TestMetricMaturity(t *testing.T) {
	m := &Metric{ID: "obs.otel.metrics", Name: "OTel metrics", Aspect: AspectConsumption,
		ConsumptionKind: ConsumptionAdoption, Maturity: coverageLadder()}

	if rung := MetricMaturity(m, nil); rung != nil {
		t.Errorf("no observation should be N/A (nil), got L%d", rung.Level)
	}

	tests := []struct {
		name  string
		value float64
		level int
	}{
		{"infrastructure in place, zero reporters", 0, 1},
		{"one reporter proves the pipeline", 0.5, 2},
		{"below majority", 45, 2},
		{"at 50", 50, 3},
		{"at 62", 62, 3},
		{"at 85", 85, 4},
		{"at 100", 100, 5},
	}
	for _, tt := range tests {
		rung := MetricMaturity(m, &Observation{MetricID: m.ID, Value: tt.value})
		if rung == nil {
			t.Fatalf("%s: nil rung", tt.name)
		}
		if rung.Level != tt.level {
			t.Errorf("%s: got L%d, want L%d", tt.name, rung.Level, tt.level)
		}
	}

	if MetricMaturity(&Metric{ID: "no-ladder"}, &Observation{Value: 50}) != nil {
		t.Error("metric without ladder should return nil")
	}

	// All-threshold ladder with a value below the lowest rung: below the ladder.
	t50 := 50.0
	allThresholds := &Metric{ID: "x", Maturity: &MaturityLadder{Levels: []MaturityRung{
		{Level: 3, Name: ">=50%", Threshold: &t50},
	}}}
	if MetricMaturity(allThresholds, &Observation{MetricID: "x", Value: 10}) != nil {
		t.Error("value below an all-threshold ladder should return nil")
	}
}

func TestMetricMaturityLowerIsBetter(t *testing.T) {
	t60, t30 := 60.0, 30.0
	m := &Metric{ID: "obs.mttr", Name: "MTTR", Aspect: AspectEffectiveness,
		Direction: "lower_is_better",
		Maturity: &MaturityLadder{Levels: []MaturityRung{
			{Level: 2, Name: "Tracked"},
			{Level: 3, Name: "<=60m", Threshold: &t60},
			{Level: 4, Name: "<=30m", Threshold: &t30},
		}},
	}

	if got := MetricMaturity(m, &Observation{MetricID: m.ID, Value: 90}).Level; got != 2 {
		t.Errorf("90m: got L%d, want L2", got)
	}
	if got := MetricMaturity(m, &Observation{MetricID: m.ID, Value: 45}).Level; got != 3 {
		t.Errorf("45m: got L%d, want L3", got)
	}
	if got := MetricMaturity(m, &Observation{MetricID: m.ID, Value: 20}).Level; got != 4 {
		t.Errorf("20m: got L%d, want L4", got)
	}
}

func TestCapabilityMaturity(t *testing.T) {
	c := &Capability{ID: "otel", Name: "OTel", Metrics: []Metric{
		{ID: "m.metrics", Name: "metrics", Aspect: AspectConsumption, ConsumptionKind: ConsumptionAdoption, Maturity: coverageLadder()},
		{ID: "m.traces", Name: "traces", Aspect: AspectConsumption, ConsumptionKind: ConsumptionAdoption, Maturity: coverageLadder()},
		{ID: "m.unladdered", Name: "unladdered", Aspect: AspectLeverage},
	}}

	a := &Assessment{FrameworkID: "f", Period: "p", Observations: []Observation{
		{MetricID: "m.metrics", Value: 85}, // L4
		{MetricID: "m.traces", Value: 45},  // L2
	}}
	rung := CapabilityMaturity(c, a)
	if rung == nil || rung.Level != 2 {
		t.Fatalf("weakest link should be L2, got %v", rung)
	}

	// An untracked laddered metric drags the capability to N/A.
	a = &Assessment{FrameworkID: "f", Period: "p", Observations: []Observation{
		{MetricID: "m.metrics", Value: 100},
	}}
	if got := CapabilityMaturity(c, a); got != nil {
		t.Errorf("untracked laddered metric should yield nil (N/A), got L%d", got.Level)
	}

	// No laddered metrics at all: nil, meaning "no maturity claim defined".
	plain := &Capability{ID: "p", Name: "p", Metrics: []Metric{{ID: "m.p", Name: "p", Aspect: AspectLeverage}}}
	if CapabilityMaturity(plain, a) != nil {
		t.Error("capability without ladders should return nil")
	}
}

func TestDomainMaturity(t *testing.T) {
	d := &Domain{ID: "obs", Name: "Obs", Capabilities: []Capability{
		{ID: "coverage", Name: "Coverage", Metrics: []Metric{
			{ID: "m.a", Name: "a", Aspect: AspectConsumption, ConsumptionKind: ConsumptionAdoption, Maturity: coverageLadder()},
		}},
		{ID: "signals", Name: "Signals", Metrics: []Metric{
			{ID: "m.b", Name: "b", Aspect: AspectEffectiveness, Maturity: coverageLadder()},
		}},
		{ID: "unladdered", Name: "Unladdered", Metrics: []Metric{
			{ID: "m.c", Name: "c", Aspect: AspectLeverage},
		}},
	}}

	// Weakest capability caps the domain: coverage L5, signals L3 → domain L3.
	a := &Assessment{FrameworkID: "f", Period: "p", Observations: []Observation{
		{MetricID: "m.a", Value: 100},
		{MetricID: "m.b", Value: 62},
	}}
	rung := DomainMaturity(d, a)
	if rung == nil || rung.Level != 3 {
		t.Fatalf("domain maturity: want L3, got %v", rung)
	}

	// An N/A capability (untracked laddered metric) makes the domain N/A —
	// a single M5 metric never carries the domain.
	a = &Assessment{FrameworkID: "f", Period: "p", Observations: []Observation{
		{MetricID: "m.a", Value: 100},
	}}
	if got := DomainMaturity(d, a); got != nil {
		t.Errorf("untracked capability should yield domain N/A, got L%d", got.Level)
	}

	// A domain with no laddered capabilities makes no claim.
	plain := &Domain{ID: "p", Name: "p", Capabilities: []Capability{
		{ID: "c", Name: "c", Metrics: []Metric{{ID: "m.p", Name: "p", Aspect: AspectLeverage}}},
	}}
	if DomainMaturity(plain, a) != nil {
		t.Error("domain without ladders should return nil")
	}
}

func TestMaturityLadderValidation(t *testing.T) {
	f := testFramework()
	f.Domains[0].Capabilities[0].Metrics[0].Maturity = coverageLadder()
	if err := f.Validate(); err != nil {
		t.Fatalf("valid ladder failed validation: %v", err)
	}

	// Non-ascending levels are invalid.
	bad := coverageLadder()
	bad.Levels[2].Level = 2
	f.Domains[0].Capabilities[0].Metrics[0].Maturity = bad
	if err := f.Validate(); err == nil {
		t.Error("expected error for non-ascending levels")
	}

	// Out-of-range level is invalid.
	bad = coverageLadder()
	bad.Levels[4].Level = 6
	f.Domains[0].Capabilities[0].Metrics[0].Maturity = bad
	if err := f.Validate(); err == nil {
		t.Error("expected error for level > 5")
	}
}
