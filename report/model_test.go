package report

import (
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/scale/catalog"
)

func TestModelHTMLAWS(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	prev := loadAssessment(t, f, "../examples/assessments/2026-q2.json")
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	html, err := ModelHTML(f, curr, "aws-observability-maturity", &Options{Prev: prev, GeneratedAt: "2026-07-21"})
	if err != nil {
		t.Fatal(err)
	}
	out := string(html)

	for _, want := range []string{
		"AWS Observability Maturity Model",
		// All four AWS stages render.
		"Foundational monitoring", "Intermediate monitoring", "Advanced observability", "Proactive observability",
		// Both observability capabilities map to intermediate → single-level position.
		"Current practice: Intermediate monitoring.",
		"► current practice",
		// The stage after the highest placement is flagged.
		"Next: <b>Advanced observability</b>",
		// Capability evidence and PRISM crosswalk appear on the ladder.
		"OpenTelemetry Instrumentation", "PRISM M3",
		// Domain evidence and journey narratives are scoped to observability.
		"SCALE evidence — Observability Best Practices",
		"adoption is no longer the bottleneck",
		// Disclaimer.
		"not a vendor assessment",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("AWS model report missing %q", want)
		}
	}
	if strings.Contains(out, "API Best Practices") {
		t.Error("AWS model report should not include other domains' sections")
	}
}

func TestModelHTMLNewRelic(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	html, err := ModelHTML(f, curr, "newrelic-observability-maturity", nil)
	if err != nil {
		t.Fatal(err)
	}
	out := string(html)

	for _, want := range []string{
		"New Relic Observability Maturity",
		// Capabilities map to reactive and proactive → spanning position.
		"Current practice spans Reactive to Proactive.",
		"Next: <b>Mastery</b>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("New Relic model report missing %q", want)
		}
	}
}

func TestModelHTMLUnknownModel(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	_, err = ModelHTML(f, curr, "nope", nil)
	if err == nil {
		t.Fatal("expected error for unknown model")
	}
	if !strings.Contains(err.Error(), "aws-observability-maturity") {
		t.Errorf("error should list available models, got: %v", err)
	}
}
