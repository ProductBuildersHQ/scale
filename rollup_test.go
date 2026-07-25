package scale

import (
	"math"
	"testing"
)

func testFramework() *Framework {
	return &Framework{
		ID:   "scale",
		Name: "SCALE",
		Domains: []Domain{
			{
				ID:   "api",
				Name: "API",
				Capabilities: []Capability{
					{
						ID:   "style",
						Name: "Style",
						Metrics: []Metric{
							{
								ID:              "api.c.adoption",
								Name:            "Adoption",
								Aspect:          AspectConsumption,
								ConsumptionKind: ConsumptionAdoption,
								Target:          &Target{Value: 100},
								Owner:           "api-platform",
							},
							{
								ID:              "api.c.conformance",
								Name:            "Conformance",
								Aspect:          AspectConsumption,
								ConsumptionKind: ConsumptionConformance,
								Target:          &Target{Value: 100},
								Owner:           "api-platform",
							},
							{
								ID:        "api.e.breaking",
								Name:      "Breaking changes",
								Aspect:    AspectEffectiveness,
								Direction: "lower_is_better",
								Target:    &Target{Value: 2},
								Owner:     "api-platform",
							},
							{
								ID:     "api.untracked",
								Name:   "No target — excluded from rollups",
								Aspect: AspectLeverage,
							},
						},
					},
				},
			},
			{
				ID:   "obs",
				Name: "Observability",
				Capabilities: []Capability{
					{
						ID:   "otel",
						Name: "OTel",
						Metrics: []Metric{
							{
								ID:              "obs.c.adoption",
								Name:            "Adoption",
								Aspect:          AspectConsumption,
								ConsumptionKind: ConsumptionAdoption,
								Target:          &Target{Value: 100},
								Owner:           "obs-platform",
							},
							{
								ID:     "obs.missing",
								Name:   "Eligible but unmeasured",
								Aspect: AspectAutomation,
								Target: &Target{Value: 100},
								Owner:  "obs-platform",
							},
						},
					},
				},
			},
		},
	}
}

func TestFrameworkValidate(t *testing.T) {
	f := testFramework()
	if err := f.Validate(); err != nil {
		t.Fatalf("valid framework failed validation: %v", err)
	}
}

func TestValidateRejectsBadMetrics(t *testing.T) {
	f := testFramework()
	f.Domains[0].Capabilities[0].Metrics[0].ConsumptionKind = "bogus"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid consumptionKind")
	}

	f = testFramework()
	f.Domains[0].Capabilities[0].Metrics[2].ConsumptionKind = ConsumptionAdoption
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for consumptionKind on non-consumption metric")
	}

	f = testFramework()
	f.Domains[1].Capabilities[0].Metrics[0].ID = "api.c.adoption"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for duplicate metric ID across domains")
	}
}

func TestAttainment(t *testing.T) {
	higher := &Metric{ID: "h", Target: &Target{Value: 100}}
	lower := &Metric{ID: "l", Direction: "lower_is_better", Target: &Target{Value: 2}}

	tests := []struct {
		name   string
		metric *Metric
		value  float64
		want   float64
	}{
		{"higher at target", higher, 100, 1},
		{"higher over target clamps", higher, 120, 1},
		{"higher partial", higher, 80, 0.8},
		{"higher zero", higher, 0, 0},
		{"lower at target", lower, 2, 1},
		{"lower under target", lower, 0, 1},
		{"lower over target", lower, 4, 0.5},
	}
	for _, tt := range tests {
		got, err := Attainment(tt.metric, tt.value)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}

	if _, err := Attainment(&Metric{ID: "no-target"}, 1); err == nil {
		t.Error("expected error for metric without target")
	}
}

func TestComputeRollup(t *testing.T) {
	f := testFramework()
	a := &Assessment{
		FrameworkID: "scale",
		Period:      "2026-Q3",
		Observations: []Observation{
			{MetricID: "api.c.adoption", Value: 90},
			{MetricID: "api.c.conformance", Value: 70},
			{MetricID: "api.e.breaking", Value: 4},
			{MetricID: "api.untracked", Value: 55},
			{MetricID: "obs.c.adoption", Value: 80},
		},
	}

	r, err := ComputeRollup(f, a)
	if err != nil {
		t.Fatal(err)
	}

	// api consumption: mean(0.9, 0.7) = 0.8
	apiRollup := r.DomainRollupFor("api")
	if apiRollup == nil {
		t.Fatal("missing api domain rollup")
	}
	var apiConsumption *AspectScore
	for i := range apiRollup.Aspects {
		if apiRollup.Aspects[i].Aspect == AspectConsumption {
			apiConsumption = &apiRollup.Aspects[i]
		}
	}
	if apiConsumption == nil {
		t.Fatal("missing api consumption score")
	}
	if math.Abs(apiConsumption.Score-0.8) > 1e-9 {
		t.Errorf("api consumption score: got %v, want 0.8", apiConsumption.Score)
	}
	if apiConsumption.MetricCount != 2 {
		t.Errorf("api consumption metric count: got %d, want 2", apiConsumption.MetricCount)
	}

	// Cross-domain consumption: mean of domain scores mean(0.8, 0.8) = 0.8.
	overall := r.AspectScoreFor(AspectConsumption)
	if overall == nil {
		t.Fatal("missing overall consumption score")
	}
	if math.Abs(overall.Score-0.8) > 1e-9 {
		t.Errorf("overall consumption score: got %v, want 0.8", overall.Score)
	}

	// api.e.breaking: lower_is_better, value 4 vs target 2 → 0.5.
	effectiveness := r.AspectScoreFor(AspectEffectiveness)
	if effectiveness == nil {
		t.Fatal("missing overall effectiveness score")
	}
	if math.Abs(effectiveness.Score-0.5) > 1e-9 {
		t.Errorf("overall effectiveness score: got %v, want 0.5", effectiveness.Score)
	}

	if len(r.Excluded) != 1 || r.Excluded[0] != "api.untracked" {
		t.Errorf("excluded: got %v, want [api.untracked]", r.Excluded)
	}
	if len(r.Missing) != 1 || r.Missing[0] != "obs.missing" {
		t.Errorf("missing: got %v, want [obs.missing]", r.Missing)
	}
}

func TestCompareRollups(t *testing.T) {
	f := testFramework()
	prev := &Assessment{
		FrameworkID: "scale",
		Period:      "2026-Q2",
		Observations: []Observation{
			{MetricID: "api.c.adoption", Value: 70},
			{MetricID: "api.c.conformance", Value: 50},
		},
	}
	curr := &Assessment{
		FrameworkID: "scale",
		Period:      "2026-Q3",
		Observations: []Observation{
			{MetricID: "api.c.adoption", Value: 90},
			{MetricID: "api.c.conformance", Value: 70},
		},
	}

	prevRollup, err := ComputeRollup(f, prev)
	if err != nil {
		t.Fatal(err)
	}
	currRollup, err := ComputeRollup(f, curr)
	if err != nil {
		t.Fatal(err)
	}

	deltas := CompareRollups(prevRollup, currRollup)
	if len(deltas) == 0 {
		t.Fatal("expected deltas")
	}
	// Consumption moved from mean(0.7,0.5)=0.6 to mean(0.9,0.7)=0.8.
	found := false
	for _, d := range deltas {
		if d.DomainID == "api" && d.Aspect == AspectConsumption {
			found = true
			if math.Abs(d.Delta-0.2) > 1e-9 {
				t.Errorf("api consumption delta: got %v, want 0.2", d.Delta)
			}
		}
	}
	if !found {
		t.Error("missing api consumption delta")
	}
}

func TestAssessmentValidate(t *testing.T) {
	f := testFramework()

	a := &Assessment{
		FrameworkID:  "scale",
		Period:       "2026-Q3",
		Observations: []Observation{{MetricID: "nope", Value: 1}},
	}
	if err := a.Validate(f); err == nil {
		t.Error("expected error for unknown metric reference")
	}

	a = &Assessment{
		FrameworkID: "scale",
		Period:      "2026-Q3",
		Narratives: []NarrativeBlock{
			{ID: "n1", Kind: NarrativeThesis, Body: "thesis does not belong here"},
		},
	}
	if err := a.Validate(f); err == nil {
		t.Error("expected error for thesis narrative in assessment")
	}
}
