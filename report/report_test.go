package report

import (
	"os"
	"strings"
	"testing"

	scale "github.com/ProductBuildersHQ/scale"
	"github.com/ProductBuildersHQ/scale/catalog"
)

func loadAssessment(t *testing.T, f *scale.Framework, path string) *scale.Assessment {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	a, err := scale.ParseAssessment(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.Validate(f); err != nil {
		t.Fatal(err)
	}
	return a
}

func TestHTMLReport(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	prev := loadAssessment(t, f, "../examples/assessments/2026-q2.json")
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	html, err := HTML(f, curr, &Options{Prev: prev, GeneratedAt: "2026-07-21"})
	if err != nil {
		t.Fatal(err)
	}
	out := string(html)

	for _, want := range []string{
		// Aspect tiles in S-C-A-L-E order.
		"Standards", "Consumption", "Automation", "Leverage", "Effectiveness",
		// Movers section exists when a prev assessment is given.
		"What moved",
		// Domain sections and their lifecycle arcs.
		"API Best Practices", "Observability Best Practices", "Security Best Practices",
		"Design &amp; Planning", "Offensive Validation",
		// External lens with provenance.
		"AWS Observability Maturity Model", "New Relic Observability Maturity",
		"◄ current practice", "normalized",
		// Numerator/denominator storytelling.
		"81/114",
		// Journey narrative from the assessment.
		"CI enforcement changed the conformance curve",
		// Coverage honesty: security pipeline-gates metric unmeasured in Q3.
		"sec.automation.pipeline-gates",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("report missing %q", want)
		}
	}

	if strings.Contains(out, "thesis does not belong") {
		t.Error("unexpected content")
	}
}

func TestHTMLReportWithoutPrev(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	html, err := HTML(f, curr, nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(html), "What moved") {
		t.Error("movers section should be absent without a previous assessment")
	}
}
