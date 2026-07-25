package scale

import (
	"errors"
	"fmt"
)

// Interpretation classifications for codified external models, so a SCALE
// normalization is never mistaken for an official vendor definition.
const (
	InterpretationDirect     = "direct"     // verbatim from the source
	InterpretationNormalized = "normalized" // restructured into SCALE/PRISM concepts
	InterpretationInferred   = "inferred"   // not explicit in the source
)

// ExternalModel is a source-faithful codification of a third-party maturity
// model (e.g. the AWS Observability Maturity Model), kept as data rather
// than machinery. The source's own levels remain authoritative; the PRISM
// crosswalk on each level is an explicit, per-model mapping — never a
// hard-coded offset — because published models differ in structure (AWS uses
// four stages; New Relic uses a Level 0 foundation plus a
// reactive→proactive→mastery progression per value driver).
//
// Capabilities and metrics reference an external model through their
// ordinary FrameworkMapping entries: Framework holds the model ID and
// Reference holds a level ID. Reports can then render the "framework lens":
// where current practice sits on each vendor's own ladder.
type ExternalModel struct {
	SchemaURI      string          `json:"$schema,omitempty"`
	ID             string          `json:"id" jsonschema:"required"`
	Name           string          `json:"name" jsonschema:"required"`
	Publisher      string          `json:"publisher" jsonschema:"required"`
	Domain         string          `json:"domain,omitempty" jsonschema:"description=SCALE domain ID this model maps to"`
	Description    string          `json:"description,omitempty"`
	SourceURL      string          `json:"sourceUrl,omitempty"`
	RetrievedAt    string          `json:"retrievedAt,omitempty" jsonschema:"description=Date the source was last checked"`
	Interpretation string          `json:"interpretation,omitempty" jsonschema:"enum=direct,enum=normalized,enum=inferred"`
	Levels         []ExternalLevel `json:"levels" jsonschema:"required"`
}

// ExternalLevel is one level of an external maturity model, in the source's
// own terms, with an optional explicit crosswalk to the PRISM 5-level scale.
type ExternalLevel struct {
	ID          string `json:"id" jsonschema:"required"`
	Ordinal     int    `json:"ordinal" jsonschema:"required,description=The source's own level number"`
	Name        string `json:"name" jsonschema:"required,description=The source's own level label"`
	Description string `json:"description,omitempty"`
	PRISMLevel  int    `json:"prismLevel,omitempty" jsonschema:"description=Explicit crosswalk to PRISM maturity (1-5); 0 means unmapped"`
}

// Level returns the level with the given ID, or nil.
func (em *ExternalModel) Level(id string) *ExternalLevel {
	for i := range em.Levels {
		if em.Levels[i].ID == id {
			return &em.Levels[i]
		}
	}
	return nil
}

// ValidInterpretation checks if an interpretation value is valid.
func ValidInterpretation(interpretation string) bool {
	switch interpretation {
	case InterpretationDirect, InterpretationNormalized, InterpretationInferred, "":
		return true
	default:
		return false
	}
}

func (em *ExternalModel) validate() []error {
	var errs []error
	loc := fmt.Sprintf("external model %q", em.ID)

	if em.ID == "" {
		errs = append(errs, errors.New("external model: id is required"))
	}
	if em.Name == "" {
		errs = append(errs, fmt.Errorf("%s: name is required", loc))
	}
	if em.Publisher == "" {
		errs = append(errs, fmt.Errorf("%s: publisher is required", loc))
	}
	if !ValidInterpretation(em.Interpretation) {
		errs = append(errs, fmt.Errorf("%s: invalid interpretation %q", loc, em.Interpretation))
	}
	if len(em.Levels) == 0 {
		errs = append(errs, fmt.Errorf("%s: at least one level is required", loc))
	}
	levelIDs := map[string]bool{}
	for i := range em.Levels {
		l := &em.Levels[i]
		if l.ID == "" {
			errs = append(errs, fmt.Errorf("%s: level[%d]: id is required", loc, i))
			continue
		}
		if levelIDs[l.ID] {
			errs = append(errs, fmt.Errorf("%s: level %q: duplicate id", loc, l.ID))
		}
		levelIDs[l.ID] = true
		if l.Name == "" {
			errs = append(errs, fmt.Errorf("%s: level %q: name is required", loc, l.ID))
		}
		if l.PRISMLevel < 0 || l.PRISMLevel > 5 {
			errs = append(errs, fmt.Errorf("%s: level %q: prismLevel must be 0-5", loc, l.ID))
		}
	}
	return errs
}
