# Assessment Format

Schema: `scale-assessment.schema.json`

An assessment records observed metric values for a single period, plus the
`journey` and `outlook` narratives for that period. The catalog holds the
timeless structure; the assessment holds what changed this cycle.

```json
{
  "$schema": "https://productbuildershq.com/schema/scale/v0/scale-assessment.schema.json",
  "frameworkId": "scale",
  "period": "2026-Q3",
  "asOf": "2026-09-30",
  "observations": [
    {
      "metricId": "obs.consumption.otel-adoption",
      "value": 82,
      "numerator": 107,
      "denominator": 130,
      "source": "telemetry-backend presence check",
      "observedAt": "2026-07-15"
    }
  ],
  "narratives": [
    {
      "id": "obs-journey",
      "kind": "journey",
      "title": "Observability Q3",
      "body": "OTel adoption expanded from 45% to 82%...",
      "domainId": "observability"
    }
  ]
}
```

## Top-level fields

| Field | Required | Purpose |
| ----- | -------- | ------- |
| `frameworkId` | Yes | Must match the framework `id` the assessment targets. |
| `period` | Yes | The reporting period (e.g., `2026-Q3`). |
| `asOf` | No | The as-of date for the period. |
| `observations` | Yes | Observed metric values. |
| `narratives` | No | `journey` and `outlook` blocks for this period. |

## Observation fields

| Field | Required | Purpose |
| ----- | -------- | ------- |
| `metricId` | Yes | References a metric in the catalog. |
| `value` | Yes | The observed value. |
| `numerator` / `denominator` | No | Show the fraction (e.g., "107/130"). |
| `source` | No | Where the data came from. |
| `observedAt` | No | When it was measured (RFC 3339). |

!!! tip "Denominators matter"
    `numerator` / `denominator` document *what was eligible*, not just what was
    counted — the difference between "82% of 130 services" and a bare "82%".

## Assessment narratives

Only time-bound narratives belong in an assessment:

- **`journey`** — what changed this period.
- **`outlook`** — where we're going next (may reference prism-roadmap initiative
  IDs so the next cycle can compare expected vs. actual movement).

Attach a narrative to a domain with `domainId`.

## Validation

```bash
scale validate -catalog catalog -assessment assessments/2026-q3.json
```

```text
catalog OK: 4 domains, 2 external models
assessment OK: period 2026-Q3, 24 observations
```
