package report

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	scale "github.com/ProductBuildersHQ/scale"
)

// ModelHTML renders the story report through the lens of one codified
// external maturity model (e.g. the AWS Observability Maturity Model): the
// vendor's ladder becomes the page structure, with the organization's
// capabilities, evidence, and current-practice placement positioned on it.
// Same catalog, same assessment, same evidence — a different projection.
func ModelHTML(f *scale.Framework, curr *scale.Assessment, modelID string, opts *Options) ([]byte, error) {
	if opts == nil {
		opts = &Options{}
	}
	em := f.ExternalModel(modelID)
	if em == nil {
		ids := make([]string, 0, len(f.ExternalModels))
		for i := range f.ExternalModels {
			ids = append(ids, f.ExternalModels[i].ID)
		}
		return nil, fmt.Errorf("unknown external model %q (available: %s)", modelID, strings.Join(ids, ", "))
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

	data := buildModelData(f, curr, em, rollup, prevRollup, opts)

	tmpl, err := template.New("model-report").Funcs(funcMap()).Parse(modelTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}
	return buf.Bytes(), nil
}

type modelData struct {
	Model       *scale.ExternalModel
	DomainName  string
	Period      string
	PrevPeriod  string
	GeneratedAt string
	Position    string
	Next        *scale.ExternalLevel
	Levels      []modelLevel
	Bars        []aspectBar
	Journey     []scale.NarrativeBlock
	Missing     []string
	Excluded    []string
}

type modelLevel struct {
	Level        *scale.ExternalLevel
	PRISM        string
	Current      bool
	IsNext       bool
	Capabilities []modelCapability
}

type modelCapability struct {
	Name     string
	Why      string
	Maturity string
	Metrics  []metricRow
}

func buildModelData(f *scale.Framework, curr *scale.Assessment, em *scale.ExternalModel, rollup, prevRollup *scale.Rollup, opts *Options) *modelData {
	data := &modelData{
		Model:       em,
		Period:      curr.Period,
		GeneratedAt: opts.GeneratedAt,
	}
	if opts.Prev != nil {
		data.PrevPeriod = opts.Prev.Period
	}

	// Gather capability placements on this model's ladder, from any domain.
	byLevel := map[string][]modelCapability{}
	for i := range f.Domains {
		d := &f.Domains[i]
		for j := range d.Capabilities {
			c := &d.Capabilities[j]
			for _, fm := range c.Frameworks {
				if fm.Framework != em.ID || fm.Reference == "" {
					continue
				}
				byLevel[fm.Reference] = append(byLevel[fm.Reference], modelCapability{
					Name:     c.Name,
					Why:      fm.Description,
					Maturity: capabilityMaturityChip(c, curr),
					Metrics:  buildCapabilityMetricRows(c, curr),
				})
			}
		}
	}

	minIdx, maxIdx := -1, -1
	for idx := range em.Levels {
		l := &em.Levels[idx]
		ml := modelLevel{Level: l, Capabilities: byLevel[l.ID]}
		if l.PRISMLevel > 0 {
			ml.PRISM = fmt.Sprintf("M%d", l.PRISMLevel)
		}
		if len(ml.Capabilities) > 0 {
			ml.Current = true
			if minIdx == -1 {
				minIdx = idx
			}
			maxIdx = idx
		}
		data.Levels = append(data.Levels, ml)
	}

	switch {
	case maxIdx == -1:
		data.Position = "No current-practice placement is declared in the catalog for this model."
	case minIdx == maxIdx:
		data.Position = fmt.Sprintf("Current practice: %s.", em.Levels[maxIdx].Name)
	default:
		data.Position = fmt.Sprintf("Current practice spans %s to %s.", em.Levels[minIdx].Name, em.Levels[maxIdx].Name)
	}
	if maxIdx >= 0 && maxIdx+1 < len(em.Levels) {
		data.Next = &em.Levels[maxIdx+1]
		data.Levels[maxIdx+1].IsNext = true
	}

	if d := f.Domain(em.Domain); d != nil {
		data.DomainName = d.Name
		var prevDR *scale.DomainRollup
		if prevRollup != nil {
			prevDR = prevRollup.DomainRollupFor(d.ID)
		}
		data.Bars = buildDomainBars(rollup.DomainRollupFor(d.ID), prevDR)
		data.Journey = scopedNarratives(curr, scale.ScopeDomain, d.ID)
		data.Missing = filterMetricsByDomain(rollup.Missing, f, d.ID)
		data.Excluded = filterMetricsByDomain(rollup.Excluded, f, d.ID)
	}

	return data
}

func filterMetricsByDomain(metricIDs []string, f *scale.Framework, domainID string) []string {
	var out []string
	for _, id := range metricIDs {
		if d, _, _ := f.Metric(id); d != nil && d.ID == domainID {
			out = append(out, id)
		}
	}
	return out
}
