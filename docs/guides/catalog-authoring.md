# Catalog Authoring Guide

This guide explains how to create and maintain SCALE catalogs — the framework.json, domain files, and external model definitions that feed the SCALE report.

## Catalog Structure

```text
catalog/
├── framework.json           # Framework metadata and narratives
├── domains/
│   ├── api.json             # One file per domain
│   ├── observability.json
│   └── security.json
└── external/
    ├── aws-observability.json      # Third-party maturity models
    └── newrelic-observability.json
```

`LoadFrameworkDir("catalog")` assembles these into a single `Framework` with validation.

## Framework Metadata

`framework.json` defines the framework itself:

```json
{
  "$schema": "https://productbuildershq.com/schema/scale/v0/scale-framework.schema.json",
  "id": "scale",
  "name": "SCALE",
  "description": "Platform engineering maturity across five aspects.",
  "version": "0.1.0",
  "owner": "platform-engineering",
  "updated": "2026-07-25",
  "narratives": [
    {
      "id": "scale-thesis",
      "kind": "thesis",
      "title": "Why SCALE exists",
      "body": "Maturity frameworks make it cheap to generate hundreds of metrics...",
      "owner": "platform-engineering",
      "reviewBy": "2027-01-01"
    }
  ]
}
```

Key fields:

- `id` — unique identifier, used in assessments (`frameworkId`)
- `version` — semver; bump on breaking schema changes
- `narratives` — framework-level thesis; journey/outlook go in assessments

## Writing a Domain

Each domain file defines one horizontal engineering discipline:

```json
{
  "$schema": "https://productbuildershq.com/schema/scale/v0/scale-domain.schema.json",
  "id": "observability",
  "name": "Observability",
  "description": "Telemetry coverage, golden signals, and incident response.",
  "status": "active",
  "owner": "observability-platform",
  "capabilities": [
    {
      "id": "otel-coverage",
      "name": "OTel Signal Coverage",
      "description": "Percentage of services reporting OTel signals.",
      "metrics": [...]
    }
  ]
}
```

### Domain Status

- `draft` — work in progress, not yet ready for assessment
- `active` — ready for measurement
- `deprecated` — being phased out

### Domain Dimensions (Optional)

For domains with a natural lifecycle, add dimensions:

```json
"dimensions": [
  {
    "id": "security-lifecycle",
    "name": "Security Lifecycle",
    "stages": [
      { "id": "design", "name": "Design & Planning" },
      { "id": "development", "name": "Development" },
      { "id": "cicd", "name": "CI/CD Pipeline" },
      { "id": "offensive", "name": "Offensive Validation" },
      { "id": "operations", "name": "Operations" }
    ]
  }
]
```

Stage order is slice order — the arc is the order.

## Writing Metrics

Metrics are the atomic measurement units:

```json
{
  "id": "obs.otel.metrics-coverage",
  "name": "OTel Metrics Coverage",
  "description": "Percentage of eligible services reporting OTel metrics.",
  "aspect": "consumption",
  "consumptionKind": "adoption",
  "direction": "higher_is_better",
  "unit": "percent",
  "target": { "value": 100 },
  "owner": "observability-platform",
  "maturity": {
    "levels": [
      { "level": 1, "name": "Tracked" },
      { "level": 2, "name": ">0%", "threshold": 0, "exclusive": true },
      { "level": 3, "name": ">=50%", "threshold": 50 },
      { "level": 4, "name": ">=80%", "threshold": 80 },
      { "level": 5, "name": "100%", "threshold": 100 }
    ]
  }
}
```

### Required Fields

| Field | Purpose |
|-------|---------|
| `id` | Unique across the framework; convention: `domain.capability.metric` |
| `name` | Human-readable label |
| `aspect` | One of: `standards`, `consumption`, `automation`, `leverage`, `effectiveness` |

### Rollup Participation

A metric participates in rollup scores **only when it has both**:
- `target` — what "done" looks like
- `owner` — who is accountable

Metrics without these are tracked but excluded from the story (reported in `excluded`).

### SCALE Aspects

| Aspect | Question | Example Metrics |
|--------|----------|-----------------|
| `standards` | Do executable standards exist? | Rule coverage, schema completeness |
| `consumption` | Are teams using them? | SDK adoption, conformance pass rate |
| `automation` | How much is enforced? | CI gate coverage, auto-remediation |
| `leverage` | What capacity do they create? | Generated SDKs, template reuse |
| `effectiveness` | Did outcomes improve? | MTTR, breaking changes, incidents |

### Consumption Sub-kinds

The `consumption` aspect requires `consumptionKind`:
- `adoption` — "the SDK is installed"
- `conformance` — "the telemetry passes the profile"

### Direction

- `higher_is_better` (default) — 100% coverage is good
- `lower_is_better` — 0 incidents is good

### Maturity Ladders

Ladders map observed values to PRISM maturity levels (1-5):

```json
"maturity": {
  "levels": [
    { "level": 1, "name": "Tracked" },
    { "level": 2, "name": ">0%", "threshold": 0, "exclusive": true },
    { "level": 3, "name": ">=50%", "threshold": 50 },
    { "level": 4, "name": ">=80%", "threshold": 80 },
    { "level": 5, "name": "100%", "threshold": 100 }
  ]
}
```

Rung semantics:
- **No threshold** — satisfied by any observation ("Tracked")
- **With threshold** — satisfied when value meets it (direction-aware)
- **`exclusive: true`** — strict comparison (`>` or `<` instead of `>=` or `<=`)

**"Not tracked" is a status, not a level.** A metric with no observation has no maturity (N/A).

**Weakest-link rule:**
- Capability maturity = min(metric rungs)
- Domain maturity = min(capability rungs)
- N/A propagates upward

### Lower-is-better Ladders

For metrics where lower is better (e.g., latency, error rate):

```json
{
  "id": "obs.gs.latency",
  "name": "Platform P95 Latency",
  "direction": "lower_is_better",
  "maturity": {
    "levels": [
      { "level": 1, "name": "Tracked" },
      { "level": 2, "name": "<=2000ms", "threshold": 2000 },
      { "level": 3, "name": "<=1000ms", "threshold": 1000 },
      { "level": 4, "name": "<=500ms", "threshold": 500 },
      { "level": 5, "name": "<=250ms", "threshold": 250 }
    ]
  }
}
```

## External Models

Codify third-party maturity models as data:

```json
{
  "$schema": "https://productbuildershq.com/schema/scale/v0/scale-external-model.schema.json",
  "id": "aws-observability-maturity",
  "name": "AWS Observability Maturity Model",
  "sourceUrl": "https://docs.aws.amazon.com/...",
  "retrievedAt": "2026-07-01",
  "interpretation": "Mapped to SCALE capabilities based on AWS documentation.",
  "levels": [
    {
      "id": "foundational",
      "name": "Foundational",
      "order": 1,
      "description": "Basic monitoring in place.",
      "prism": { "level": 2 }
    },
    {
      "id": "intermediate",
      "name": "Intermediate",
      "order": 2,
      "description": "Centralized observability platform.",
      "prism": { "level": 3 }
    }
  ]
}
```

Capabilities reference external models via `frameworks` mappings:

```json
{
  "id": "otel-coverage",
  "frameworks": [
    {
      "framework": "aws-observability-maturity",
      "reference": "intermediate"
    }
  ]
}
```

## Validation

Validate your catalog before committing:

```bash
scale validate -catalog catalog
```

Common validation errors:
- Duplicate metric IDs
- Invalid aspect values
- Maturity levels not 1-5 or not ascending
- Missing required fields

Validate with an assessment:

```bash
scale validate -catalog catalog -assessment assessments/2026-q3.json
```

## Writing Assessments

Assessments record observations for a period:

```json
{
  "$schema": "https://productbuildershq.com/schema/scale/v0/scale-assessment.schema.json",
  "frameworkId": "scale",
  "period": "2026-Q3",
  "observations": [
    {
      "metricId": "obs.otel.metrics-coverage",
      "value": 62,
      "numerator": 124,
      "denominator": 200,
      "source": "telemetry-backend presence check",
      "observedAt": "2026-07-15"
    }
  ],
  "narratives": [
    {
      "id": "obs-journey",
      "kind": "journey",
      "title": "Observability Q3",
      "body": "OTel coverage expanded from 45% to 62%...",
      "domainId": "observability"
    }
  ]
}
```

### Observation Fields

| Field | Required | Purpose |
|-------|----------|---------|
| `metricId` | Yes | References a metric in the catalog |
| `value` | Yes | The observed value |
| `numerator` / `denominator` | No | Show the fraction (e.g., "124/200") |
| `source` | No | Where the data came from |
| `observedAt` | No | When measured (RFC 3339) |

### Assessment Narratives

- `journey` — what changed this period (time-bound)
- `outlook` — where we're going (forward-looking)

Attach to a domain with `domainId`.

## Generating Reports

```bash
# Aspect story report
scale report -catalog catalog -assessment assessments/2026-q3.json \
  -prev assessments/2026-q2.json -o reports/2026-q3.html

# External model lens
scale report -catalog catalog -assessment assessments/2026-q3.json \
  -model aws-observability-maturity -o reports/aws-2026-q3.html
```

## Exporting for Web Component

```bash
scale export -catalog catalog -o web/framework.json
```

The exported JSON IR can be consumed by `<scale-report>`.

## Best Practices

1. **Metric IDs are stable** — changing them breaks assessment history
2. **One aspect per metric** — if a metric spans aspects, split it
3. **Denominators matter** — document what's eligible, not just what's counted
4. **Ladders need L1** — give L1 no threshold to avoid "below the ladder" gaps
5. **Owners are accountable** — the owner field drives rollup participation
6. **Review dates prevent stale prose** — set `reviewBy` on thesis narratives
