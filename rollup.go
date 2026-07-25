package scale

import (
	"fmt"
	"math"
	"sort"

	core "github.com/grokify/prism-core"
)

// Rollup is the computed story skeleton for one assessment period: aspect
// scores per domain, cross-domain aspect scores, and the bookkeeping a
// report needs to be honest about coverage (which metrics were excluded and
// which eligible metrics went unmeasured).
//
// Rollup semantics are deliberately simple so every number remains
// explainable in a leadership meeting:
//
//   - Attainment per metric is value vs. target, clamped to [0,1].
//   - A domain's aspect score is the plain mean of attainment across that
//     domain's rollup-eligible, observed metrics with that aspect.
//   - The overall aspect score is the plain mean of domain aspect scores, so
//     metric-heavy domains do not dominate the horizontal story.
type Rollup struct {
	FrameworkID string         `json:"frameworkId"`
	Period      string         `json:"period"`
	Aspects     []AspectScore  `json:"aspects,omitempty" jsonschema:"description=Cross-domain aspect scores (mean of domain aspect scores)"`
	Domains     []DomainRollup `json:"domains,omitempty"`
	Excluded    []string       `json:"excluded,omitempty" jsonschema:"description=Metric IDs observed but not rollup-eligible (missing target or owner)"`
	Missing     []string       `json:"missing,omitempty" jsonschema:"description=Rollup-eligible metric IDs with no observation this period"`
}

// DomainRollup holds one domain's aspect scores.
type DomainRollup struct {
	DomainID string        `json:"domainId"`
	Aspects  []AspectScore `json:"aspects,omitempty"`
}

// AspectScore is a normalized [0,1] score for one SCALE aspect, with the
// per-metric contributions retained for attribution ("driven by API
// conformance, +18pts").
type AspectScore struct {
	Aspect        string             `json:"aspect"`
	Score         float64            `json:"score"`
	MetricCount   int                `json:"metricCount"`
	Contributions []MetricAttainment `json:"contributions,omitempty"`
}

// MetricAttainment is one metric's contribution to an aspect score.
type MetricAttainment struct {
	DomainID     string  `json:"domainId"`
	CapabilityID string  `json:"capabilityId"`
	MetricID     string  `json:"metricId"`
	Aspect       string  `json:"aspect"`
	Value        float64 `json:"value"`
	Attainment   float64 `json:"attainment"`
}

// AspectDelta is the movement of one aspect score between two rollups.
// An empty DomainID means the cross-domain score.
type AspectDelta struct {
	DomainID string  `json:"domainId,omitempty"`
	Aspect   string  `json:"aspect"`
	Prev     float64 `json:"prev"`
	Curr     float64 `json:"curr"`
	Delta    float64 `json:"delta"`
}

// Attainment normalizes an observed value against the metric's target to
// [0,1]. Higher-is-better: value/target. Lower-is-better: 1 when value is at
// or under target, otherwise target/value. Requires a target (see
// Metric.RollupEligible).
func Attainment(m *Metric, value float64) (float64, error) {
	if m.Target == nil {
		return 0, fmt.Errorf("metric %q: attainment requires a target", m.ID)
	}
	switch m.EffectiveDirection() {
	case core.SLIDirectionLowerIsBetter:
		if m.Target.Value < 0 {
			return 0, fmt.Errorf("metric %q: lower-is-better target must be non-negative", m.ID)
		}
		if value <= m.Target.Value {
			return 1, nil
		}
		return clamp01(m.Target.Value / value), nil
	default:
		if m.Target.Value <= 0 {
			return 0, fmt.Errorf("metric %q: higher-is-better target must be positive", m.ID)
		}
		return clamp01(value / m.Target.Value), nil
	}
}

// ComputeRollup computes aspect scores for an assessment against its
// framework catalog.
func ComputeRollup(f *Framework, a *Assessment) (*Rollup, error) {
	if err := a.Validate(f); err != nil {
		return nil, fmt.Errorf("assessment invalid: %w", err)
	}

	r := &Rollup{FrameworkID: f.ID, Period: a.Period}
	overall := map[string][]float64{} // aspect -> per-domain scores

	for i := range f.Domains {
		d := &f.Domains[i]
		byAspect := map[string][]MetricAttainment{}

		for j := range d.Capabilities {
			c := &d.Capabilities[j]
			for k := range c.Metrics {
				m := &c.Metrics[k]
				obs := a.Observation(m.ID)
				if !m.RollupEligible() {
					if obs != nil {
						r.Excluded = append(r.Excluded, m.ID)
					}
					continue
				}
				if obs == nil {
					r.Missing = append(r.Missing, m.ID)
					continue
				}
				att, err := Attainment(m, obs.Value)
				if err != nil {
					return nil, err
				}
				byAspect[m.Aspect] = append(byAspect[m.Aspect], MetricAttainment{
					DomainID:     d.ID,
					CapabilityID: c.ID,
					MetricID:     m.ID,
					Aspect:       m.Aspect,
					Value:        obs.Value,
					Attainment:   att,
				})
			}
		}

		if len(byAspect) == 0 {
			continue
		}
		dr := DomainRollup{DomainID: d.ID}
		for _, aspect := range AllAspects() {
			contribs := byAspect[aspect]
			if len(contribs) == 0 {
				continue
			}
			score := meanAttainment(contribs)
			dr.Aspects = append(dr.Aspects, AspectScore{
				Aspect:        aspect,
				Score:         score,
				MetricCount:   len(contribs),
				Contributions: contribs,
			})
			overall[aspect] = append(overall[aspect], score)
		}
		r.Domains = append(r.Domains, dr)
	}

	for _, aspect := range AllAspects() {
		scores := overall[aspect]
		if len(scores) == 0 {
			continue
		}
		var sum float64
		var count int
		for _, s := range scores {
			sum += s
		}
		for _, dr := range r.Domains {
			for _, as := range dr.Aspects {
				if as.Aspect == aspect {
					count += as.MetricCount
				}
			}
		}
		r.Aspects = append(r.Aspects, AspectScore{
			Aspect:      aspect,
			Score:       sum / float64(len(scores)),
			MetricCount: count,
		})
	}

	sort.Strings(r.Excluded)
	sort.Strings(r.Missing)
	return r, nil
}

// AspectScoreFor returns the cross-domain score for an aspect, or nil.
func (r *Rollup) AspectScoreFor(aspect string) *AspectScore {
	for i := range r.Aspects {
		if r.Aspects[i].Aspect == aspect {
			return &r.Aspects[i]
		}
	}
	return nil
}

// DomainRollupFor returns the rollup for a domain, or nil.
func (r *Rollup) DomainRollupFor(domainID string) *DomainRollup {
	for i := range r.Domains {
		if r.Domains[i].DomainID == domainID {
			return &r.Domains[i]
		}
	}
	return nil
}

// CompareRollups returns aspect movements between two rollups — the raw
// material for the "movers" layer of a story report. Deltas are returned for
// the cross-domain scores and every domain present in either rollup, sorted
// by absolute movement, largest first. Aspects absent from one side are
// skipped: a score appearing or disappearing is a coverage change, not a
// movement.
func CompareRollups(prev, curr *Rollup) []AspectDelta {
	var deltas []AspectDelta

	deltas = append(deltas, compareAspects("", prev.Aspects, curr.Aspects)...)

	domainIDs := map[string]bool{}
	for _, dr := range prev.Domains {
		domainIDs[dr.DomainID] = true
	}
	for _, dr := range curr.Domains {
		domainIDs[dr.DomainID] = true
	}
	for domainID := range domainIDs {
		var prevAspects, currAspects []AspectScore
		if dr := prev.DomainRollupFor(domainID); dr != nil {
			prevAspects = dr.Aspects
		}
		if dr := curr.DomainRollupFor(domainID); dr != nil {
			currAspects = dr.Aspects
		}
		deltas = append(deltas, compareAspects(domainID, prevAspects, currAspects)...)
	}

	sort.SliceStable(deltas, func(i, j int) bool {
		return math.Abs(deltas[i].Delta) > math.Abs(deltas[j].Delta)
	})
	return deltas
}

func compareAspects(domainID string, prev, curr []AspectScore) []AspectDelta {
	var deltas []AspectDelta
	for _, p := range prev {
		for _, c := range curr {
			if p.Aspect == c.Aspect {
				deltas = append(deltas, AspectDelta{
					DomainID: domainID,
					Aspect:   p.Aspect,
					Prev:     p.Score,
					Curr:     c.Score,
					Delta:    c.Score - p.Score,
				})
			}
		}
	}
	return deltas
}

func meanAttainment(contribs []MetricAttainment) float64 {
	var sum float64
	for _, c := range contribs {
		sum += c.Attainment
	}
	return sum / float64(len(contribs))
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
