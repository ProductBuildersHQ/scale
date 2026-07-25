package catalog

import (
	"testing"

	scale "github.com/ProductBuildersHQ/scale"
)

func TestDefaultCatalogLoadsAndValidates(t *testing.T) {
	f, err := Default()
	if err != nil {
		t.Fatal(err)
	}

	for _, domainID := range []string{"api", "observability", "security"} {
		if f.Domain(domainID) == nil {
			t.Errorf("missing domain %q", domainID)
		}
	}

	for _, modelID := range []string{"aws-observability-maturity", "newrelic-observability-maturity"} {
		em := f.ExternalModel(modelID)
		if em == nil {
			t.Errorf("missing external model %q", modelID)
			continue
		}
		if em.Domain != "observability" {
			t.Errorf("external model %q: domain %q, want observability", modelID, em.Domain)
		}
	}

	sec := f.Domain("security")
	dim := sec.Dimension("security-lifecycle")
	if dim == nil {
		t.Fatal("missing security-lifecycle dimension")
	}
	if len(dim.Stages) != 5 {
		t.Errorf("security lifecycle stages: got %d, want 5", len(dim.Stages))
	}
	if dim.Stages[0].ID != "design-planning" || dim.Stages[4].ID != "operations" {
		t.Errorf("security lifecycle stage order wrong: first %q, last %q", dim.Stages[0].ID, dim.Stages[4].ID)
	}
}

func TestSeedDomainsCoverAllAspects(t *testing.T) {
	f, err := Default()
	if err != nil {
		t.Fatal(err)
	}

	// The two active seed domains must exercise all five aspects.
	for _, domainID := range []string{"api", "observability"} {
		d := f.Domain(domainID)
		aspects := map[string]bool{}
		for _, c := range d.Capabilities {
			for _, m := range c.Metrics {
				aspects[m.Aspect] = true
				if !m.RollupEligible() {
					t.Errorf("domain %s: metric %s is not rollup-eligible (needs target and owner)", domainID, m.ID)
				}
			}
		}
		for _, aspect := range scale.AllAspects() {
			if !aspects[aspect] {
				t.Errorf("domain %s: no metric for aspect %s", domainID, aspect)
			}
		}
	}
}
