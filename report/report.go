// Package report renders a SCALE framework and assessment as a
// self-contained HTML story report: aspect headlines, movers with
// attribution, per-domain journeys along their lifecycle dimensions,
// external framework lenses (e.g. AWS and New Relic maturity ladders), and
// honest coverage bookkeeping. It deliberately depends only on the scale
// module — PRISM and other assessment platforms remain downstream consumers
// of SCALE, never dependencies of it.
package report

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	scale "github.com/ProductBuildersHQ/scale"
)

// Options configures report generation.
type Options struct {
	// Prev is an optional prior-period assessment; when present the report
	// includes deltas and the movers section.
	Prev *scale.Assessment
	// GeneratedAt is a display timestamp supplied by the caller.
	GeneratedAt string
}

// HTML renders the story report. The assessment is validated against the
// framework and rolled up internally.
func HTML(f *scale.Framework, curr *scale.Assessment, opts *Options) ([]byte, error) {
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

	data := buildData(f, curr, rollup, prevRollup, opts)

	tmpl, err := template.New("report").Funcs(funcMap()).Parse(reportTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}
	return buf.Bytes(), nil
}

type reportData struct {
	Framework   *scale.Framework
	Period      string
	PrevPeriod  string
	GeneratedAt string
	Narratives  []scale.NarrativeBlock
	Tiles       []tile
	Movers      []mover
	Domains     []domainView
	Excluded    []string
	Missing     []string
}

type tile struct {
	Letter   string
	Name     string
	HasScore bool
	Score    float64 // 0-100
	HasDelta bool
	Delta    float64 // points
}

type mover struct {
	Domain        string
	Aspect        string
	AspectLetter  string
	Prev          float64 // 0-100
	Curr          float64
	Delta         float64 // points
	Contributions []string
}

type domainView struct {
	Domain     *scale.Domain
	Maturity   string
	Thesis     []scale.NarrativeBlock
	Journey    []scale.NarrativeBlock
	Dimensions []scale.DomainDimension
	Bars       []aspectBar
	Groups     []capGroup
	Lenses     []lensView
}

type capGroup struct {
	Name     string
	Maturity string
	Rows     []metricRow
}

type aspectBar struct {
	Letter   string
	Name     string
	Score    float64 // 0-100
	HasDelta bool
	Delta    float64 // points
}

type metricRow struct {
	Name         string
	AspectLetter string
	AspectName   string
	Kind         string
	Value        string
	Target       string
	HasAttain    bool
	Attain       float64 // 0-100
	Owner        string
	Note         string
	Maturity     string
}

type lensView struct {
	Model *scale.ExternalModel
	Rows  []lensRow
}

type lensRow struct {
	Level        *scale.ExternalLevel
	PRISM        string
	Current      bool
	Capabilities []string
}

func buildData(f *scale.Framework, curr *scale.Assessment, rollup, prevRollup *scale.Rollup, opts *Options) *reportData {
	data := &reportData{
		Framework:   f,
		Period:      curr.Period,
		GeneratedAt: opts.GeneratedAt,
		Narratives:  filterKind(f.Narratives, scale.NarrativeThesis),
		Excluded:    rollup.Excluded,
		Missing:     rollup.Missing,
	}
	if opts.Prev != nil {
		data.PrevPeriod = opts.Prev.Period
	}

	for _, aspect := range scale.AllAspects() {
		t := tile{Letter: scale.AspectLetter(aspect), Name: scale.AspectDisplayName(aspect)}
		if as := rollup.AspectScoreFor(aspect); as != nil {
			t.HasScore = true
			t.Score = as.Score * 100
			if prevRollup != nil {
				if ps := prevRollup.AspectScoreFor(aspect); ps != nil {
					t.HasDelta = true
					t.Delta = (as.Score - ps.Score) * 100
				}
			}
		}
		data.Tiles = append(data.Tiles, t)
	}

	if prevRollup != nil {
		data.Movers = buildMovers(f, rollup, prevRollup)
	}

	for i := range f.Domains {
		d := &f.Domains[i]
		dv := domainView{
			Domain:     d,
			Maturity:   domainMaturityChip(d, curr),
			Thesis:     filterKind(d.Narratives, scale.NarrativeThesis),
			Journey:    scopedNarratives(curr, scale.ScopeDomain, d.ID),
			Dimensions: d.Dimensions,
			Lenses:     buildLenses(f, d),
		}
		var prevDR *scale.DomainRollup
		if prevRollup != nil {
			prevDR = prevRollup.DomainRollupFor(d.ID)
		}
		dv.Bars = buildDomainBars(rollup.DomainRollupFor(d.ID), prevDR)
		dv.Groups = buildCapGroups(d, curr)
		data.Domains = append(data.Domains, dv)
	}

	return data
}

func buildMovers(f *scale.Framework, rollup, prevRollup *scale.Rollup) []mover {
	prevAttain := map[string]float64{}
	for _, dr := range prevRollup.Domains {
		for _, as := range dr.Aspects {
			for _, c := range as.Contributions {
				prevAttain[c.MetricID] = c.Attainment
			}
		}
	}

	deltas := scale.CompareRollups(prevRollup, rollup)
	var movers []mover
	for _, delta := range deltas {
		if delta.DomainID == "" || delta.Delta == 0 {
			continue
		}
		d := f.Domain(delta.DomainID)
		if d == nil {
			continue
		}
		m := mover{
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
				sort.SliceStable(cds, func(i, j int) bool {
					return abs(cds[i].delta) > abs(cds[j].delta)
				})
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

func buildDomainBars(dr, prevDR *scale.DomainRollup) []aspectBar {
	if dr == nil {
		return nil
	}
	var bars []aspectBar
	for _, as := range dr.Aspects {
		bar := aspectBar{
			Letter: scale.AspectLetter(as.Aspect),
			Name:   scale.AspectDisplayName(as.Aspect),
			Score:  as.Score * 100,
		}
		if prevDR != nil {
			for _, pas := range prevDR.Aspects {
				if pas.Aspect == as.Aspect {
					bar.HasDelta = true
					bar.Delta = (as.Score - pas.Score) * 100
				}
			}
		}
		bars = append(bars, bar)
	}
	return bars
}

func buildCapGroups(d *scale.Domain, a *scale.Assessment) []capGroup {
	var groups []capGroup
	for i := range d.Capabilities {
		c := &d.Capabilities[i]
		rows := buildCapabilityMetricRows(c, a)
		if len(rows) == 0 {
			continue
		}
		groups = append(groups, capGroup{
			Name:     c.Name,
			Maturity: capabilityMaturityChip(c, a),
			Rows:     rows,
		})
	}
	return groups
}

// capabilityMaturityChip renders the capability's weakest-link maturity, or
// "N/A" when a laddered metric is untracked — a capability cannot claim a
// level its weakest metric has not earned. Empty when no metric has a ladder.
func capabilityMaturityChip(c *scale.Capability, a *scale.Assessment) string {
	if !c.HasLadderedMetrics() {
		return ""
	}
	if rung := scale.CapabilityMaturity(c, a); rung != nil {
		return fmt.Sprintf("L%d · %s", rung.Level, rung.Name)
	}
	return "N/A"
}

// domainMaturityChip renders the domain's weakest-link maturity across its
// laddered capabilities, or "N/A" when any laddered capability has no
// maturity. Empty when no capability carries ladders.
func domainMaturityChip(d *scale.Domain, a *scale.Assessment) string {
	hasLadder := false
	for i := range d.Capabilities {
		if d.Capabilities[i].HasLadderedMetrics() {
			hasLadder = true
			break
		}
	}
	if !hasLadder {
		return ""
	}
	if rung := scale.DomainMaturity(d, a); rung != nil {
		return fmt.Sprintf("L%d · %s", rung.Level, rung.Name)
	}
	return "N/A"
}

func buildCapabilityMetricRows(c *scale.Capability, a *scale.Assessment) []metricRow {
	var rows []metricRow
	for j := range c.Metrics {
		m := &c.Metrics[j]
		row := metricRow{
			Name:         m.Name,
			AspectLetter: scale.AspectLetter(m.Aspect),
			AspectName:   scale.AspectDisplayName(m.Aspect),
			Kind:         m.ConsumptionKind,
			Owner:        m.Owner,
		}
		if m.Target != nil {
			row.Target = formatValue(m, m.Target.Value, nil, nil)
		}
		obs := a.Observation(m.ID)
		if m.Maturity != nil {
			switch rung := scale.MetricMaturity(m, obs); {
			case rung != nil:
				row.Maturity = fmt.Sprintf("L%d · %s", rung.Level, rung.Name)
			case obs == nil:
				row.Maturity = "N/A · not tracked"
			default:
				row.Maturity = "below ladder"
			}
		}
		switch {
		case obs == nil && m.RollupEligible():
			row.Value = "—"
			row.Note = "not measured"
		case obs == nil:
			row.Value = "—"
			row.Note = "tracked only (no target/owner)"
		default:
			row.Value = formatValue(m, obs.Value, obs.Numerator, obs.Denominator)
			if m.RollupEligible() {
				if att, err := scale.Attainment(m, obs.Value); err == nil {
					row.HasAttain = true
					row.Attain = att * 100
				}
			} else {
				row.Note = "tracked only (no target/owner)"
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func buildLenses(f *scale.Framework, d *scale.Domain) []lensView {
	var lenses []lensView
	for i := range f.ExternalModels {
		em := &f.ExternalModels[i]
		if em.Domain != d.ID {
			continue
		}
		current := map[string][]string{}
		for j := range d.Capabilities {
			c := &d.Capabilities[j]
			for _, fm := range c.Frameworks {
				if fm.Framework == em.ID && fm.Reference != "" {
					current[fm.Reference] = append(current[fm.Reference], c.Name)
				}
			}
		}
		lv := lensView{Model: em}
		for j := range em.Levels {
			l := &em.Levels[j]
			row := lensRow{Level: l, Capabilities: current[l.ID]}
			row.Current = len(row.Capabilities) > 0
			if l.PRISMLevel > 0 {
				row.PRISM = fmt.Sprintf("M%d", l.PRISMLevel)
			}
			lv.Rows = append(lv.Rows, row)
		}
		lenses = append(lenses, lv)
	}
	return lenses
}

func formatValue(m *scale.Metric, value float64, num, den *int) string {
	base := trimFloat(value)
	switch m.Unit {
	case scale.UnitPercent:
		base += "%"
	case "":
	default:
		base += " " + m.Unit
	}
	if num != nil && den != nil {
		return fmt.Sprintf("%d/%d (%s)", *num, *den, base)
	}
	return base
}

func trimFloat(v float64) string {
	s := fmt.Sprintf("%.1f", v)
	if s[len(s)-2:] == ".0" {
		return s[:len(s)-2]
	}
	return s
}

func filterKind(blocks []scale.NarrativeBlock, kind string) []scale.NarrativeBlock {
	var out []scale.NarrativeBlock
	for _, b := range blocks {
		if b.Kind == kind {
			out = append(out, b)
		}
	}
	return out
}

func scopedNarratives(a *scale.Assessment, scopeType, ref string) []scale.NarrativeBlock {
	var out []scale.NarrativeBlock
	for _, b := range a.Narratives {
		if b.Scope != nil && b.Scope.Type == scopeType && b.Scope.Ref == ref {
			out = append(out, b)
		}
	}
	return out
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"pct": func(v float64) string {
			return fmt.Sprintf("%.0f", v)
		},
		"pts": func(v float64) string {
			return fmt.Sprintf("%+.0f", v)
		},
		"kindLabel": func(kind string) string {
			if kind == "" {
				return ""
			}
			return kind
		},
	}
}
