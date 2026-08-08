# Framework & Domain Reference

The catalog is stored as a framework file plus one file per domain, assembled by
`LoadFrameworkDir`. JSON schemas are generated from the Go types
(`go run ./cmd/schemagen`) and embedded via the `schema` package.

```text
catalog/
├── framework.json           # Framework metadata and narratives
├── domains/
│   ├── api.json             # One file per domain
│   ├── observability.json
│   ├── platform.json
│   └── security.json
└── external/
    ├── aws-observability.json      # Third-party maturity models
    └── newrelic-observability.json
```

For a task-oriented walkthrough, see [Catalog Authoring](../guide/catalog-authoring.md).

## Framework

Schema: `scale-framework.schema.json`

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

| Field | Required | Purpose |
| ----- | -------- | ------- |
| `id` | Yes | Unique identifier; referenced by assessments as `frameworkId`. |
| `name` | Yes | Display name. |
| `version` | Yes | Semver; bump on breaking schema changes. |
| `owner` | No | Accountable team. |
| `updated` | No | Last-updated date. |
| `narratives` | No | Framework-level `thesis` blocks (journey/outlook live in assessments). |

## Domain

Schema: `scale-domain.schema.json`

```json
{
  "$schema": "https://productbuildershq.com/schema/scale/v0/scale-domain.schema.json",
  "id": "observability",
  "name": "Observability",
  "description": "Telemetry coverage, golden signals, and incident response.",
  "status": "active",
  "owner": "observability-platform",
  "dimensions": [ ... ],
  "capabilities": [ ... ],
  "narratives": [ ... ]
}
```

### Status

| Status | Meaning |
| ------ | ------- |
| `draft` | Work in progress, not yet ready for assessment. |
| `active` | Ready for measurement. |
| `deprecated` | Being phased out. |

### Dimensions

Ordered lifecycle stages — the domain's story spine. **Stage order is slice
order — the arc is the order.**

```json
"dimensions": [
  {
    "id": "security-lifecycle",
    "name": "Security Lifecycle",
    "stages": [
      { "id": "design",      "name": "Design & Planning" },
      { "id": "development",  "name": "Development" },
      { "id": "cicd",         "name": "CI/CD Pipeline" },
      { "id": "offensive",    "name": "Offensive Validation" },
      { "id": "operations",   "name": "Operations" }
    ]
  }
]
```

### Capabilities & metrics

A capability groups metrics within a domain. Each metric is tagged with exactly
one [aspect](../guide/aspects.md) and participates in rollups only when it has
both a `target` and an `owner`.

```json
{
  "id": "otel-coverage",
  "name": "OTel Signal Coverage",
  "metrics": [
    {
      "id": "obs.consumption.otel-adoption",
      "name": "OTel Adoption",
      "aspect": "consumption",
      "consumptionKind": "adoption",
      "direction": "higher_is_better",
      "unit": "percent",
      "target": { "value": 100 },
      "owner": "observability-platform",
      "maturity": { "levels": [ ... ] }
    }
  ]
}
```

See [Catalog Authoring → Writing Metrics](../guide/catalog-authoring.md#writing-metrics)
for the full field list, maturity ladders, and lower-is-better handling.

## External models

Schema: `scale-external-model.schema.json`. Codify third-party maturity models
as data, with source-faithful levels, provenance, and a PRISM crosswalk.
Capabilities reference them through ordinary `frameworks` mappings
(`framework` = model ID, `reference` = level ID). See
[Catalog Authoring → External Models](../guide/catalog-authoring.md#external-models).
