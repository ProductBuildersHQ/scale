package scale

// SCALE aspect constants. Every metric is tagged with exactly one aspect.
// These five values are frozen; changing them is a breaking change to the
// specification.
const (
	AspectStandards     = "standards"     // S — Do executable engineering standards exist?
	AspectConsumption   = "consumption"   // C — Are teams using them (adoption) and using them correctly (conformance)?
	AspectAutomation    = "automation"    // A — How much is enforced automatically?
	AspectLeverage      = "leverage"      // L — What reuse and engineering capacity do they create?
	AspectEffectiveness = "effectiveness" // E — Did engineering and business outcomes improve?
)

// Consumption metric kinds. Consumption is the only aspect with sub-kinds
// because adoption ("SDK installed") and conformance ("telemetry actually
// passes the profile") are different facts that tell different stories.
const (
	ConsumptionAdoption    = "adoption"
	ConsumptionConformance = "conformance"
)

// Narrative kind constants. See NarrativeBlock.
const (
	NarrativeThesis  = "thesis"  // timeless: why this item matters
	NarrativeJourney = "journey" // time-bound: what changed in a period
	NarrativeOutlook = "outlook" // forward: where we are going, tied to roadmap initiatives
)

// AllAspects returns the five SCALE aspects in canonical S-C-A-L-E order.
func AllAspects() []string {
	return []string{
		AspectStandards,
		AspectConsumption,
		AspectAutomation,
		AspectLeverage,
		AspectEffectiveness,
	}
}

// ValidAspect checks if an aspect value is valid.
func ValidAspect(aspect string) bool {
	switch aspect {
	case AspectStandards, AspectConsumption, AspectAutomation, AspectLeverage, AspectEffectiveness:
		return true
	default:
		return false
	}
}

// ValidConsumptionKind checks if a consumption kind value is valid.
func ValidConsumptionKind(kind string) bool {
	switch kind {
	case ConsumptionAdoption, ConsumptionConformance:
		return true
	default:
		return false
	}
}

// AspectLetter returns the single-letter abbreviation for an aspect.
func AspectLetter(aspect string) string {
	letters := map[string]string{
		AspectStandards:     "S",
		AspectConsumption:   "C",
		AspectAutomation:    "A",
		AspectLeverage:      "L",
		AspectEffectiveness: "E",
	}
	if l, ok := letters[aspect]; ok {
		return l
	}
	return "?"
}

// AspectDisplayName returns a human-readable name for an aspect.
func AspectDisplayName(aspect string) string {
	names := map[string]string{
		AspectStandards:     "Standards",
		AspectConsumption:   "Consumption",
		AspectAutomation:    "Automation",
		AspectLeverage:      "Leverage",
		AspectEffectiveness: "Effectiveness",
	}
	if n, ok := names[aspect]; ok {
		return n
	}
	return aspect
}

// AspectSortWeight returns the canonical S-C-A-L-E ordering weight.
func AspectSortWeight(aspect string) int {
	weights := map[string]int{
		AspectStandards:     1,
		AspectConsumption:   2,
		AspectAutomation:    3,
		AspectLeverage:      4,
		AspectEffectiveness: 5,
	}
	if w, ok := weights[aspect]; ok {
		return w
	}
	return 99
}

// ValidNarrativeKind checks if a narrative kind value is valid.
func ValidNarrativeKind(kind string) bool {
	switch kind {
	case NarrativeThesis, NarrativeJourney, NarrativeOutlook:
		return true
	default:
		return false
	}
}
