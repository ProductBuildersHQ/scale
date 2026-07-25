package scale

import (
	"errors"
	"fmt"

	core "github.com/grokify/prism-core"
)

// Validate checks the framework for structural consistency: ID uniqueness,
// valid enum values, resolvable stage references, and consumption-kind
// rules. It returns all problems found, joined via errors.Join.
func (f *Framework) Validate() error {
	var errs []error

	if f.ID == "" {
		errs = append(errs, errors.New("framework: id is required"))
	}
	if f.Name == "" {
		errs = append(errs, errors.New("framework: name is required"))
	}
	errs = append(errs, validateNarratives("framework", f.Narratives)...)

	domainIDs := map[string]bool{}
	metricIDs := map[string]bool{}

	for i := range f.Domains {
		d := &f.Domains[i]
		if d.ID == "" {
			errs = append(errs, fmt.Errorf("domain[%d]: id is required", i))
			continue
		}
		if domainIDs[d.ID] {
			errs = append(errs, fmt.Errorf("domain %q: duplicate id", d.ID))
		}
		domainIDs[d.ID] = true
		errs = append(errs, d.validate(metricIDs)...)
	}

	externalIDs := map[string]bool{}
	for i := range f.ExternalModels {
		em := &f.ExternalModels[i]
		if em.ID != "" && externalIDs[em.ID] {
			errs = append(errs, fmt.Errorf("external model %q: duplicate id", em.ID))
		}
		externalIDs[em.ID] = true
		errs = append(errs, em.validate()...)
		if em.Domain != "" && !domainIDs[em.Domain] {
			errs = append(errs, fmt.Errorf("external model %q: unknown domain %q", em.ID, em.Domain))
		}
	}
	errs = append(errs, f.validateExternalRefs()...)

	return errors.Join(errs...)
}

// validateExternalRefs checks that framework mappings which name a codified
// external model reference one of that model's level IDs. Mappings to
// frameworks that are not codified here (NIST_CSF_2, DORA, ...) are not
// checked — those live in prism-core.
func (f *Framework) validateExternalRefs() []error {
	var errs []error
	check := func(loc string, mappings []core.FrameworkMapping) {
		for _, fm := range mappings {
			em := f.ExternalModel(fm.Framework)
			if em == nil || fm.Reference == "" {
				continue
			}
			if em.Level(fm.Reference) == nil {
				errs = append(errs, fmt.Errorf("%s: framework mapping %q: unknown level %q", loc, fm.Framework, fm.Reference))
			}
		}
	}
	for i := range f.Domains {
		d := &f.Domains[i]
		for j := range d.Capabilities {
			c := &d.Capabilities[j]
			check(fmt.Sprintf("domain %q: capability %q", d.ID, c.ID), c.Frameworks)
			for k := range c.Metrics {
				m := &c.Metrics[k]
				check(fmt.Sprintf("domain %q: capability %q: metric %q", d.ID, c.ID, m.ID), m.Frameworks)
			}
		}
	}
	return errs
}

func (d *Domain) validate(metricIDs map[string]bool) []error {
	var errs []error

	if d.Name == "" {
		errs = append(errs, fmt.Errorf("domain %q: name is required", d.ID))
	}
	if !ValidStatus(d.Status) {
		errs = append(errs, fmt.Errorf("domain %q: invalid status %q", d.ID, d.Status))
	}
	errs = append(errs, validateNarratives("domain "+d.ID, d.Narratives)...)

	dimensionStages := map[string]map[string]bool{}
	for i := range d.Dimensions {
		dim := &d.Dimensions[i]
		if dim.ID == "" {
			errs = append(errs, fmt.Errorf("domain %q: dimension[%d]: id is required", d.ID, i))
			continue
		}
		if _, ok := dimensionStages[dim.ID]; ok {
			errs = append(errs, fmt.Errorf("domain %q: dimension %q: duplicate id", d.ID, dim.ID))
		}
		stages := map[string]bool{}
		for j := range dim.Stages {
			s := &dim.Stages[j]
			if s.ID == "" {
				errs = append(errs, fmt.Errorf("domain %q: dimension %q: stage[%d]: id is required", d.ID, dim.ID, j))
				continue
			}
			if stages[s.ID] {
				errs = append(errs, fmt.Errorf("domain %q: dimension %q: stage %q: duplicate id", d.ID, dim.ID, s.ID))
			}
			stages[s.ID] = true
			errs = append(errs, validateNarratives(fmt.Sprintf("domain %s stage %s", d.ID, s.ID), s.Narratives)...)
		}
		dimensionStages[dim.ID] = stages
	}

	capIDs := map[string]bool{}
	for i := range d.Capabilities {
		c := &d.Capabilities[i]
		if c.ID == "" {
			errs = append(errs, fmt.Errorf("domain %q: capability[%d]: id is required", d.ID, i))
			continue
		}
		if capIDs[c.ID] {
			errs = append(errs, fmt.Errorf("domain %q: capability %q: duplicate id", d.ID, c.ID))
		}
		capIDs[c.ID] = true
		errs = append(errs, c.validate(d.ID, dimensionStages, metricIDs)...)
	}

	return errs
}

func (c *Capability) validate(domainID string, dimensionStages map[string]map[string]bool, metricIDs map[string]bool) []error {
	var errs []error
	loc := fmt.Sprintf("domain %q: capability %q", domainID, c.ID)

	if c.Name == "" {
		errs = append(errs, fmt.Errorf("%s: name is required", loc))
	}
	if !ValidStatus(c.Status) {
		errs = append(errs, fmt.Errorf("%s: invalid status %q", loc, c.Status))
	}
	for _, ref := range c.Stages {
		stages, ok := dimensionStages[ref.Dimension]
		if !ok {
			errs = append(errs, fmt.Errorf("%s: stage ref: unknown dimension %q", loc, ref.Dimension))
			continue
		}
		if !stages[ref.Stage] {
			errs = append(errs, fmt.Errorf("%s: stage ref: unknown stage %q in dimension %q", loc, ref.Stage, ref.Dimension))
		}
	}
	errs = append(errs, validateNarratives(loc, c.Narratives)...)

	for i := range c.Metrics {
		m := &c.Metrics[i]
		if m.ID == "" {
			errs = append(errs, fmt.Errorf("%s: metric[%d]: id is required", loc, i))
			continue
		}
		if metricIDs[m.ID] {
			errs = append(errs, fmt.Errorf("%s: metric %q: duplicate id (metric IDs are framework-global)", loc, m.ID))
		}
		metricIDs[m.ID] = true
		errs = append(errs, m.validate(loc)...)
	}

	return errs
}

func (m *Metric) validate(loc string) []error {
	var errs []error
	mloc := fmt.Sprintf("%s: metric %q", loc, m.ID)

	if m.Name == "" {
		errs = append(errs, fmt.Errorf("%s: name is required", mloc))
	}
	if !ValidAspect(m.Aspect) {
		errs = append(errs, fmt.Errorf("%s: invalid aspect %q", mloc, m.Aspect))
	}
	if m.Aspect == AspectConsumption {
		if !ValidConsumptionKind(m.ConsumptionKind) {
			errs = append(errs, fmt.Errorf("%s: consumption metrics require consumptionKind of %q or %q", mloc, ConsumptionAdoption, ConsumptionConformance))
		}
	} else if m.ConsumptionKind != "" {
		errs = append(errs, fmt.Errorf("%s: consumptionKind is only valid for consumption metrics", mloc))
	}
	switch m.Direction {
	case "", core.SLIDirectionHigherIsBetter, core.SLIDirectionLowerIsBetter:
	default:
		errs = append(errs, fmt.Errorf("%s: invalid direction %q", mloc, m.Direction))
	}
	if m.Target != nil {
		switch m.EffectiveDirection() {
		case core.SLIDirectionHigherIsBetter:
			if m.Target.Value <= 0 {
				errs = append(errs, fmt.Errorf("%s: higher-is-better target must be positive", mloc))
			}
		case core.SLIDirectionLowerIsBetter:
			if m.Target.Value < 0 {
				errs = append(errs, fmt.Errorf("%s: lower-is-better target must be non-negative", mloc))
			}
		}
	}
	if m.Maturity != nil {
		errs = append(errs, m.Maturity.validate(mloc)...)
	}

	return errs
}

func validateNarratives(loc string, blocks []NarrativeBlock) []error {
	var errs []error
	for i := range blocks {
		n := &blocks[i]
		if n.ID == "" {
			errs = append(errs, fmt.Errorf("%s: narrative[%d]: id is required", loc, i))
		}
		if !ValidNarrativeKind(n.Kind) {
			errs = append(errs, fmt.Errorf("%s: narrative %q: invalid kind %q", loc, n.ID, n.Kind))
		}
		if n.Body == "" {
			errs = append(errs, fmt.Errorf("%s: narrative %q: body is required", loc, n.ID))
		}
		if n.Scope != nil && !ValidScopeType(n.Scope.Type) {
			errs = append(errs, fmt.Errorf("%s: narrative %q: invalid scope type %q", loc, n.ID, n.Scope.Type))
		}
	}
	return errs
}

// Validate checks an assessment against its framework: the framework ID must
// match, every observation must reference a known metric, and narrative
// blocks must be journey or outlook kinds with resolvable scopes.
func (a *Assessment) Validate(f *Framework) error {
	var errs []error

	if a.FrameworkID == "" {
		errs = append(errs, errors.New("assessment: frameworkId is required"))
	} else if f != nil && a.FrameworkID != f.ID {
		errs = append(errs, fmt.Errorf("assessment: frameworkId %q does not match framework %q", a.FrameworkID, f.ID))
	}
	if a.Period == "" {
		errs = append(errs, errors.New("assessment: period is required"))
	}

	seen := map[string]bool{}
	for i := range a.Observations {
		o := &a.Observations[i]
		if o.MetricID == "" {
			errs = append(errs, fmt.Errorf("assessment: observation[%d]: metricId is required", i))
			continue
		}
		if seen[o.MetricID] {
			errs = append(errs, fmt.Errorf("assessment: metric %q: duplicate observation", o.MetricID))
		}
		seen[o.MetricID] = true
		if f != nil {
			if _, _, m := f.Metric(o.MetricID); m == nil {
				errs = append(errs, fmt.Errorf("assessment: observation references unknown metric %q", o.MetricID))
			}
		}
	}

	for i := range a.Narratives {
		n := &a.Narratives[i]
		if !ValidNarrativeKind(n.Kind) {
			errs = append(errs, fmt.Errorf("assessment: narrative %q: invalid kind %q", n.ID, n.Kind))
		}
		if n.Kind == NarrativeThesis {
			errs = append(errs, fmt.Errorf("assessment: narrative %q: thesis blocks belong in the framework catalog, not assessments", n.ID))
		}
		if n.Body == "" {
			errs = append(errs, fmt.Errorf("assessment: narrative[%d] %q: body is required", i, n.ID))
		}
		if n.Scope != nil && !ValidScopeType(n.Scope.Type) {
			errs = append(errs, fmt.Errorf("assessment: narrative %q: invalid scope type %q", n.ID, n.Scope.Type))
		}
	}

	return errors.Join(errs...)
}
