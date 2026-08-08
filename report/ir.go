package report

import (
	"encoding/json"
	"fmt"

	scale "github.com/ProductBuildersHQ/scale"
)

// ReportIR is the intermediate representation of a SCALE report.
// It contains all computed data needed to render a report in any format
// (HTML, JSON, Slack, etc.) without requiring access to the original
// Framework or Assessment objects.
type ReportIR struct {
	// Period is the assessment period (e.g., "2026-Q3").
	Period string `json:"period"`

	// PrevPeriod is the prior period if comparison data is present.
	PrevPeriod string `json:"prevPeriod,omitempty"`

	// GeneratedAt is the timestamp when this IR was built.
	GeneratedAt string `json:"generatedAt,omitempty"`

	// Framework contains framework-level metadata.
	Framework FrameworkIR `json:"framework"`

	// Aspects contains the five SCALE aspect scores.
	Aspects []AspectIR `json:"aspects"`

	// Movers lists significant period-over-period changes.
	Movers []MoverIR `json:"movers,omitempty"`

	// Domains contains per-domain detailed breakdowns.
	Domains []DomainIR `json:"domains"`

	// Coverage tracks rollup completeness.
	Coverage CoverageIR `json:"coverage"`
}

// FrameworkIR contains framework-level metadata.
type FrameworkIR struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Narratives  []NarrativeIR `json:"narratives,omitempty"`
}

// AspectIR represents one of the five SCALE aspects.
type AspectIR struct {
	// Aspect is the aspect identifier (e.g., "standards").
	Aspect string `json:"aspect"`

	// Letter is the single-letter code (S, C, A, L, E).
	Letter string `json:"letter"`

	// Name is the display name (e.g., "Standards").
	Name string `json:"name"`

	// Score is 0-100, nil if no data.
	Score *float64 `json:"score"`

	// Delta is the period-over-period change in points, nil if no prior.
	Delta *float64 `json:"delta,omitempty"`
}

// MoverIR represents a significant change between periods.
type MoverIR struct {
	Domain        string   `json:"domain"`
	Aspect        string   `json:"aspect"`
	AspectLetter  string   `json:"aspectLetter"`
	Prev          float64  `json:"prev"`
	Curr          float64  `json:"curr"`
	Delta         float64  `json:"delta"`
	Contributions []string `json:"contributions,omitempty"`
}

// DomainIR represents a single domain's report data.
type DomainIR struct {
	// ID is the domain identifier.
	ID string `json:"id"`

	// Name is the display name.
	Name string `json:"name"`

	// Description is the domain description.
	Description string `json:"description,omitempty"`

	// Status is the domain status (active, draft, deprecated).
	Status string `json:"status,omitempty"`

	// Maturity is the weakest-link maturity label if applicable.
	Maturity string `json:"maturity,omitempty"`

	// Narratives are domain-level thesis/journey narratives.
	Narratives []NarrativeIR `json:"narratives,omitempty"`

	// Dimensions are the domain's lifecycle stages.
	Dimensions []DimensionIR `json:"dimensions,omitempty"`

	// Aspects contains domain-level aspect scores.
	Aspects []AspectIR `json:"aspects"`

	// Capabilities contains per-capability breakdowns.
	Capabilities []CapabilityIR `json:"capabilities"`

	// ExternalModels contains external framework lens views.
	ExternalModels []ExternalModelIR `json:"externalModels,omitempty"`
}

// DimensionIR represents a domain lifecycle dimension.
type DimensionIR struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Stages      []StageIR `json:"stages"`
}

// StageIR represents a stage in a dimension.
type StageIR struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CapabilityIR represents a capability with its metrics.
type CapabilityIR struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	Maturity    string     `json:"maturity,omitempty"`
	Metrics     []MetricIR `json:"metrics"`
}

// MetricIR represents a single metric's computed state.
type MetricIR struct {
	// ID is the metric identifier.
	ID string `json:"id"`

	// Name is the display name.
	Name string `json:"name"`

	// Description is the metric description.
	Description string `json:"description,omitempty"`

	// Aspect is the SCALE aspect (standards, consumption, etc.).
	Aspect string `json:"aspect"`

	// AspectLetter is the single-letter code.
	AspectLetter string `json:"aspectLetter"`

	// ConsumptionKind is the sub-type for consumption metrics.
	ConsumptionKind string `json:"consumptionKind,omitempty"`

	// Unit is the metric unit (percent, count, score, etc.).
	Unit string `json:"unit,omitempty"`

	// Direction indicates if lower is better.
	Direction string `json:"direction,omitempty"`

	// Owner is the metric owner.
	Owner string `json:"owner,omitempty"`

	// Value is the observed value, formatted as string.
	Value string `json:"value"`

	// RawValue is the numeric observed value, nil if not measured.
	RawValue *float64 `json:"rawValue,omitempty"`

	// Numerator is the count numerator if applicable.
	Numerator *int `json:"numerator,omitempty"`

	// Denominator is the count denominator if applicable.
	Denominator *int `json:"denominator,omitempty"`

	// Target is the target value, formatted as string.
	Target string `json:"target,omitempty"`

	// TargetRaw is the numeric target value.
	TargetRaw *float64 `json:"targetRaw,omitempty"`

	// Attainment is 0-100 if computable, nil otherwise.
	Attainment *float64 `json:"attainment,omitempty"`

	// Maturity is the maturity ladder position if applicable.
	Maturity string `json:"maturity,omitempty"`

	// Note provides context (e.g., "not measured").
	Note string `json:"note,omitempty"`

	// RollupEligible indicates if this metric contributes to rollup.
	RollupEligible bool `json:"rollupEligible"`
}

// NarrativeIR represents a narrative block.
type NarrativeIR struct {
	ID    string `json:"id,omitempty"`
	Kind  string `json:"kind"`
	Title string `json:"title,omitempty"`
	Body  string `json:"body"`
	Owner string `json:"owner,omitempty"`
}

// ExternalModelIR represents an external framework lens view.
type ExternalModelIR struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Publisher   string            `json:"publisher,omitempty"`
	URL         string            `json:"url,omitempty"`
	Description string            `json:"description,omitempty"`
	Levels      []ExternalLevelIR `json:"levels"`
}

// ExternalLevelIR represents a level in an external framework.
type ExternalLevelIR struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	PRISMLevel   string   `json:"prismLevel,omitempty"`
	IsCurrent    bool     `json:"isCurrent,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

// CoverageIR tracks rollup completeness.
type CoverageIR struct {
	// Missing lists metric IDs that are rollup-eligible but not observed.
	Missing []string `json:"missing,omitempty"`

	// Excluded lists metric IDs excluded from rollup (no target/owner).
	Excluded []string `json:"excluded,omitempty"`
}

// JSON builds a ReportIR and returns it as formatted JSON.
func JSON(f *scale.Framework, curr *scale.Assessment, opts *Options) ([]byte, error) {
	ir, err := BuildIR(f, curr, opts)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(ir, "", "  ")
}

// BuildIR constructs a ReportIR from the framework, assessment, and rollup.
func BuildIR(f *scale.Framework, curr *scale.Assessment, opts *Options) (*ReportIR, error) {
	if opts == nil {
		opts = &Options{}
	}

	rollup, err := scale.ComputeRollup(f, curr)
	if err != nil {
		return nil, fmt.Errorf("computing rollup: %w", err)
	}

	var prevRollup *scale.Rollup
	if opts.Prev != nil {
		prevRollup, err = scale.ComputeRollup(f, opts.Prev)
		if err != nil {
			return nil, fmt.Errorf("computing previous rollup: %w", err)
		}
	}

	ir := &ReportIR{
		Period:      curr.Period,
		GeneratedAt: opts.GeneratedAt,
		Framework: FrameworkIR{
			ID:          f.ID,
			Name:        f.Name,
			Description: f.Description,
		},
		Coverage: CoverageIR{
			Missing:  rollup.Missing,
			Excluded: rollup.Excluded,
		},
	}

	if opts.Prev != nil {
		ir.PrevPeriod = opts.Prev.Period
	}

	// Framework narratives
	for _, n := range f.Narratives {
		if n.Kind == scale.NarrativeThesis {
			ir.Framework.Narratives = append(ir.Framework.Narratives, NarrativeIR{
				ID:    n.ID,
				Kind:  n.Kind,
				Title: n.Title,
				Body:  n.Body,
				Owner: n.Owner,
			})
		}
	}

	// Build aspect tiles
	for _, aspect := range scale.AllAspects() {
		air := AspectIR{
			Aspect: aspect,
			Letter: scale.AspectLetter(aspect),
			Name:   scale.AspectDisplayName(aspect),
		}
		if as := rollup.AspectScoreFor(aspect); as != nil {
			score := as.Score * 100
			air.Score = &score
			if prevRollup != nil {
				if ps := prevRollup.AspectScoreFor(aspect); ps != nil {
					delta := (as.Score - ps.Score) * 100
					air.Delta = &delta
				}
			}
		}
		ir.Aspects = append(ir.Aspects, air)
	}

	// Build movers
	if prevRollup != nil {
		ir.Movers = buildMoversIR(f, rollup, prevRollup)
	}

	// Build domains
	for i := range f.Domains {
		d := &f.Domains[i]
		dir := buildDomainIR(f, d, curr, rollup, prevRollup)
		ir.Domains = append(ir.Domains, dir)
	}

	return ir, nil
}

func buildMoversIR(f *scale.Framework, rollup, prevRollup *scale.Rollup) []MoverIR {
	prevAttain := map[string]float64{}
	for _, dr := range prevRollup.Domains {
		for _, as := range dr.Aspects {
			for _, c := range as.Contributions {
				prevAttain[c.MetricID] = c.Attainment
			}
		}
	}

	deltas := scale.CompareRollups(prevRollup, rollup)
	var movers []MoverIR
	for _, delta := range deltas {
		if delta.DomainID == "" || delta.Delta == 0 {
			continue
		}
		d := f.Domain(delta.DomainID)
		if d == nil {
			continue
		}
		m := MoverIR{
			Domain:       d.Name,
			Aspect:       scale.AspectDisplayName(delta.Aspect),
			AspectLetter: scale.AspectLetter(delta.Aspect),
			Prev:         delta.Prev * 100,
			Curr:         delta.Curr * 100,
			Delta:        delta.Delta * 100,
		}
		if dr := rollup.DomainRollupFor(delta.DomainID); dr != nil {
			for _, as := range dr.Aspects {
				if as.Aspect != delta.Aspect {
					continue
				}
				type contribDelta struct {
					name  string
					delta float64
				}
				var cds []contribDelta
				for _, c := range as.Contributions {
					prev, ok := prevAttain[c.MetricID]
					if !ok {
						continue
					}
					_, _, metric := f.Metric(c.MetricID)
					name := c.MetricID
					if metric != nil {
						name = metric.Name
					}
					cds = append(cds, contribDelta{name: name, delta: (c.Attainment - prev) * 100})
				}
				for i, cd := range cds {
					if i >= 2 || cd.delta == 0 {
						break
					}
					m.Contributions = append(m.Contributions, fmt.Sprintf("%s (%+.0f pts)", cd.name, cd.delta))
				}
			}
		}
		movers = append(movers, m)
		if len(movers) == 6 {
			break
		}
	}
	return movers
}

func buildDomainIR(f *scale.Framework, d *scale.Domain, a *scale.Assessment, rollup, prevRollup *scale.Rollup) DomainIR {
	dir := DomainIR{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Status:      d.Status,
		Maturity:    domainMaturityChip(d, a),
	}

	// Domain narratives (thesis from catalog)
	for _, n := range d.Narratives {
		if n.Kind == scale.NarrativeThesis {
			dir.Narratives = append(dir.Narratives, NarrativeIR{
				ID:    n.ID,
				Kind:  n.Kind,
				Title: n.Title,
				Body:  n.Body,
				Owner: n.Owner,
			})
		}
	}

	// Journey narratives from assessment
	for _, n := range a.Narratives {
		if n.Scope != nil && n.Scope.Type == scale.ScopeDomain && n.Scope.Ref == d.ID {
			dir.Narratives = append(dir.Narratives, NarrativeIR{
				ID:    n.ID,
				Kind:  n.Kind,
				Title: n.Title,
				Body:  n.Body,
				Owner: n.Owner,
			})
		}
	}

	// Dimensions
	for _, dim := range d.Dimensions {
		dimIR := DimensionIR{
			ID:          dim.ID,
			Name:        dim.Name,
			Description: dim.Description,
		}
		for _, s := range dim.Stages {
			dimIR.Stages = append(dimIR.Stages, StageIR{
				ID:          s.ID,
				Name:        s.Name,
				Description: s.Description,
			})
		}
		dir.Dimensions = append(dir.Dimensions, dimIR)
	}

	// Domain aspect scores
	if dr := rollup.DomainRollupFor(d.ID); dr != nil {
		var prevDR *scale.DomainRollup
		if prevRollup != nil {
			prevDR = prevRollup.DomainRollupFor(d.ID)
		}
		for _, as := range dr.Aspects {
			air := AspectIR{
				Aspect: as.Aspect,
				Letter: scale.AspectLetter(as.Aspect),
				Name:   scale.AspectDisplayName(as.Aspect),
			}
			score := as.Score * 100
			air.Score = &score
			if prevDR != nil {
				for _, pas := range prevDR.Aspects {
					if pas.Aspect == as.Aspect {
						delta := (as.Score - pas.Score) * 100
						air.Delta = &delta
					}
				}
			}
			dir.Aspects = append(dir.Aspects, air)
		}
	}

	// Capabilities
	for i := range d.Capabilities {
		c := &d.Capabilities[i]
		cir := buildCapabilityIR(c, a)
		dir.Capabilities = append(dir.Capabilities, cir)
	}

	// External models
	for i := range f.ExternalModels {
		em := &f.ExternalModels[i]
		if em.Domain != d.ID {
			continue
		}
		current := map[string][]string{}
		for j := range d.Capabilities {
			cap := &d.Capabilities[j]
			for _, fm := range cap.Frameworks {
				if fm.Framework == em.ID && fm.Reference != "" {
					current[fm.Reference] = append(current[fm.Reference], cap.Name)
				}
			}
		}
		emir := ExternalModelIR{
			ID:          em.ID,
			Name:        em.Name,
			Publisher:   em.Publisher,
			URL:         em.SourceURL,
			Description: em.Description,
		}
		for j := range em.Levels {
			l := &em.Levels[j]
			lir := ExternalLevelIR{
				ID:           l.ID,
				Name:         l.Name,
				Description:  l.Description,
				Capabilities: current[l.ID],
				IsCurrent:    len(current[l.ID]) > 0,
			}
			if l.PRISMLevel > 0 {
				lir.PRISMLevel = fmt.Sprintf("M%d", l.PRISMLevel)
			}
			emir.Levels = append(emir.Levels, lir)
		}
		dir.ExternalModels = append(dir.ExternalModels, emir)
	}

	return dir
}

func buildCapabilityIR(c *scale.Capability, a *scale.Assessment) CapabilityIR {
	cir := CapabilityIR{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Status:      c.Status,
		Maturity:    capabilityMaturityChip(c, a),
	}

	for j := range c.Metrics {
		m := &c.Metrics[j]
		mir := MetricIR{
			ID:              m.ID,
			Name:            m.Name,
			Description:     m.Description,
			Aspect:          m.Aspect,
			AspectLetter:    scale.AspectLetter(m.Aspect),
			ConsumptionKind: m.ConsumptionKind,
			Unit:            m.Unit,
			Direction:       m.Direction,
			Owner:           m.Owner,
			RollupEligible:  m.RollupEligible(),
		}

		if m.Target != nil {
			mir.Target = formatValue(m, m.Target.Value, nil, nil)
			mir.TargetRaw = &m.Target.Value
		}

		obs := a.Observation(m.ID)
		if m.Maturity != nil {
			switch rung := scale.MetricMaturity(m, obs); {
			case rung != nil:
				mir.Maturity = fmt.Sprintf("L%d · %s", rung.Level, rung.Name)
			case obs == nil:
				mir.Maturity = "N/A · not tracked"
			default:
				mir.Maturity = "below ladder"
			}
		}

		switch {
		case obs == nil && m.RollupEligible():
			mir.Value = "—"
			mir.Note = "not measured"
		case obs == nil:
			mir.Value = "—"
			mir.Note = "tracked only (no target/owner)"
		default:
			mir.Value = formatValue(m, obs.Value, obs.Numerator, obs.Denominator)
			mir.RawValue = &obs.Value
			mir.Numerator = obs.Numerator
			mir.Denominator = obs.Denominator
			if m.RollupEligible() {
				if att, err := scale.Attainment(m, obs.Value); err == nil {
					attPct := att * 100
					mir.Attainment = &attPct
				}
			} else {
				mir.Note = "tracked only (no target/owner)"
			}
		}

		cir.Metrics = append(cir.Metrics, mir)
	}

	return cir
}
