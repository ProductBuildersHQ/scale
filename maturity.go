package scale

import (
	"fmt"

	core "github.com/grokify/prism-core"
)

// MaturityLadder maps a metric's observed value to a maturity level (1-5),
// making PRISM-style maturity computable from SCALE observations instead of
// judged.
//
// "Not tracked" is a status, not a level: a metric with no observation has
// no maturity (MetricMaturity returns nil, reported as "N/A · not tracked").
// Absence of measurement earns no rung. A typical coverage ladder therefore
// spends all five levels on attainment:
//
//	L1 "Tracked"   (no threshold — measurement reports a number, even 0;
//	                within the measuring team's own control)
//	L2 ">0%"       (threshold 0, exclusive — at least one real reporter,
//	                the pipeline is proven end-to-end)
//	L3 ">=50%"     (threshold 50)
//	L4 ">=80%"     (threshold 80)
//	L5 "100%"      (threshold 100)
//
// Rung semantics:
//
//   - A rung without a threshold is satisfied by the existence of an
//     observation ("Tracked").
//   - A rung with a threshold is satisfied when the observed value meets it,
//     direction-aware: >= for higher-is-better metrics, <= for
//     lower-is-better metrics; with exclusive true the comparison is strict
//     (> or <), which is how ">0%" is expressed.
//
// The metric's maturity is the highest satisfied rung; if an observation
// exists but satisfies no rung (an all-threshold ladder with a value below
// the lowest), the metric is "below the ladder" and MetricMaturity returns
// nil — give L1 no threshold to avoid that gap.
type MaturityLadder struct {
	Levels []MaturityRung `json:"levels" jsonschema:"required"`
}

// MaturityRung is one level of a MaturityLadder.
type MaturityRung struct {
	Level       int      `json:"level" jsonschema:"required,description=PRISM maturity level 1-5"`
	Name        string   `json:"name" jsonschema:"required"`
	Threshold   *float64 `json:"threshold,omitempty" jsonschema:"description=Observed-value threshold (>= for higher-is-better, <= for lower-is-better); omit for a state rung satisfied by any observation (Tracked)"`
	Exclusive   bool     `json:"exclusive,omitempty" jsonschema:"description=Strict comparison (> or <) instead of >= / <=; use with threshold 0 to express '>0%'"`
	Description string   `json:"description,omitempty"`
}

// MetricMaturity returns the maturity rung a metric holds given its
// observation, per the ladder semantics above. It returns nil when the
// metric has no ladder, when there is no observation (not tracked — N/A),
// or when the observed value satisfies no rung (below the ladder).
func MetricMaturity(m *Metric, obs *Observation) *MaturityRung {
	if m.Maturity == nil || len(m.Maturity.Levels) == 0 || obs == nil {
		return nil
	}
	higher := m.EffectiveDirection() != core.SLIDirectionLowerIsBetter
	var best *MaturityRung
	rungs := m.Maturity.Levels
	for i := range rungs {
		r := &rungs[i]
		switch {
		case r.Threshold == nil:
			best = r // satisfied by the observation's existence
		case higher && r.Exclusive && obs.Value > *r.Threshold:
			best = r
		case higher && !r.Exclusive && obs.Value >= *r.Threshold:
			best = r
		case !higher && r.Exclusive && obs.Value < *r.Threshold:
			best = r
		case !higher && !r.Exclusive && obs.Value <= *r.Threshold:
			best = r
		}
	}
	return best
}

// CapabilityMaturity returns the capability's maturity as the weakest link:
// the lowest rung held across the capability's laddered metrics. It returns
// nil when the capability has no laddered metrics, or when any laddered
// metric has no maturity (not tracked or below the ladder) — a capability
// cannot claim a level its weakest metric has not earned.
func CapabilityMaturity(c *Capability, a *Assessment) *MaturityRung {
	var min *MaturityRung
	for i := range c.Metrics {
		m := &c.Metrics[i]
		if m.Maturity == nil {
			continue
		}
		rung := MetricMaturity(m, a.Observation(m.ID))
		if rung == nil {
			return nil
		}
		if min == nil || rung.Level < min.Level {
			min = rung
		}
	}
	return min
}

// HasLadderedMetrics reports whether any of the capability's metrics carry a
// maturity ladder — i.e. whether the capability makes a maturity claim at all.
func (c *Capability) HasLadderedMetrics() bool {
	for i := range c.Metrics {
		if c.Metrics[i].Maturity != nil {
			return true
		}
	}
	return false
}

// DomainMaturity returns the domain's maturity as the weakest link across
// its capabilities that carry maturity ladders: the lowest capability rung.
// Maturity claims are conjunctive — a domain cannot claim a level its
// weakest capability has not earned, and a single high-scoring metric never
// lifts the domain. It returns nil when no capability carries ladders (no
// claim is made), or when any laddered capability is N/A (an untracked or
// below-ladder metric) — the domain cannot claim a level across the board
// while part of the board is unmeasured.
func DomainMaturity(d *Domain, a *Assessment) *MaturityRung {
	var min *MaturityRung
	for i := range d.Capabilities {
		c := &d.Capabilities[i]
		if !c.HasLadderedMetrics() {
			continue
		}
		rung := CapabilityMaturity(c, a)
		if rung == nil {
			return nil
		}
		if min == nil || rung.Level < min.Level {
			min = rung
		}
	}
	return min
}

func (ml *MaturityLadder) validate(loc string) []error {
	var errs []error
	if len(ml.Levels) < 2 {
		errs = append(errs, fmt.Errorf("%s: maturity ladder needs at least 2 levels", loc))
	}
	prevLevel := 0
	for i := range ml.Levels {
		r := &ml.Levels[i]
		if r.Level < 1 || r.Level > 5 {
			errs = append(errs, fmt.Errorf("%s: maturity level %d must be 1-5", loc, r.Level))
		}
		if r.Level <= prevLevel {
			errs = append(errs, fmt.Errorf("%s: maturity levels must be strictly ascending (level %d after %d)", loc, r.Level, prevLevel))
		}
		prevLevel = r.Level
		if r.Name == "" {
			errs = append(errs, fmt.Errorf("%s: maturity level %d: name is required", loc, r.Level))
		}
	}
	return errs
}
