package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/scale/catalog"
)

func TestBuildIR(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	ir, err := BuildIR(f, curr, &Options{GeneratedAt: "2026-08-04"})
	if err != nil {
		t.Fatalf("BuildIR: %v", err)
	}

	if ir.Period != "2026-Q3" {
		t.Errorf("Period = %q, want %q", ir.Period, "2026-Q3")
	}

	if len(ir.Aspects) != 5 {
		t.Errorf("len(Aspects) = %d, want 5", len(ir.Aspects))
	}

	// Check we have domains
	if len(ir.Domains) == 0 {
		t.Error("no domains in IR")
	}

	// Check platform domain has capabilities
	var platform *DomainIR
	for i := range ir.Domains {
		if ir.Domains[i].ID == "platform" {
			platform = &ir.Domains[i]
			break
		}
	}
	if platform == nil {
		t.Fatal("platform domain not found")
	}
	if len(platform.Capabilities) == 0 {
		t.Error("platform domain has no capabilities")
	}
}

func TestBuildIRWithPrev(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	prev := loadAssessment(t, f, "../examples/assessments/2026-q2.json")
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	ir, err := BuildIR(f, curr, &Options{Prev: prev, GeneratedAt: "2026-08-04"})
	if err != nil {
		t.Fatalf("BuildIR: %v", err)
	}

	if ir.PrevPeriod != "2026-Q2" {
		t.Errorf("PrevPeriod = %q, want %q", ir.PrevPeriod, "2026-Q2")
	}

	// Should have movers with prev data
	// (may be empty if no changes, but shouldn't error)
}

func TestIRJSONRoundtrip(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	ir, err := BuildIR(f, curr, &Options{GeneratedAt: "2026-08-04"})
	if err != nil {
		t.Fatalf("BuildIR: %v", err)
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(ir, "", "  ")
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	// Deserialize back
	var ir2 ReportIR
	if err := json.Unmarshal(data, &ir2); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	// Basic equality checks
	if ir2.Period != ir.Period {
		t.Errorf("Period mismatch after roundtrip")
	}
	if len(ir2.Aspects) != len(ir.Aspects) {
		t.Errorf("Aspects count mismatch after roundtrip")
	}
	if len(ir2.Domains) != len(ir.Domains) {
		t.Errorf("Domains count mismatch after roundtrip")
	}
}

func TestHTMLFromIR(t *testing.T) {
	f, err := catalog.Default()
	if err != nil {
		t.Fatal(err)
	}
	curr := loadAssessment(t, f, "../examples/assessments/2026-q3.json")

	ir, err := BuildIR(f, curr, &Options{GeneratedAt: "2026-08-04"})
	if err != nil {
		t.Fatalf("BuildIR: %v", err)
	}

	html, err := HTMLFromIR(ir)
	if err != nil {
		t.Fatalf("HTMLFromIR: %v", err)
	}

	if len(html) == 0 {
		t.Error("empty HTML output")
	}

	// Should contain key elements
	s := string(html)
	for _, want := range []string{"SCALE", "2026-Q3", "Platform Adoption"} {
		if !strings.Contains(s, want) {
			t.Errorf("HTML missing %q", want)
		}
	}
}
