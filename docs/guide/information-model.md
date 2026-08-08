# Information Model

```text
Framework
    └── Domain (api, observability, security, platform, ...)
        ├── DomainDimension (ordered lifecycle stages — the story spine)
        └── Capability
            └── Metric (each tagged with exactly one SCALE aspect)
```

- **Framework** — the top-level container. Holds framework-level narratives and
  the assembled domains.
- **Domain** — one horizontal engineering discipline (API best practices,
  observability, security, platform adoption).
- **DomainDimension** — an ordered lifecycle for domains that have a natural
  arc. Stage order is slice order.
- **Capability** — a coherent unit of practice within a domain.
- **Metric** — the atomic measurement unit, tagged with exactly one
  [aspect](aspects.md).

## Narratives

`NarrativeBlock`s attach at every level — the "literate" half of a SCALE
document, analogous to markdown cells in a notebook. Three kinds with different
lifetimes:

| Kind | Lifetime | Home | Purpose |
| --------- | ---------- | ----------------- | ---------------------------------------------------- |
| `thesis` | Timeless | Framework catalog | Why this element matters. Carries `reviewBy` so stale prose is visibly stale. |
| `journey` | Time-bound | Assessment | Where we were and what changed this period. |
| `outlook` | Forward | Assessment | Where we are going; references prism-roadmap initiative IDs so the next cycle can compare expected vs. actual movement. |

Because `thesis` lives in the catalog and `journey` / `outlook` live in the
assessment, the timeless "why" is authored once, while the per-period "what
changed" is authored each cycle.

## Rollup Semantics

Rollup rules are part of the specification, not the report generator, so every
score means the same thing in every report.

- A metric participates in rollups **only when it has both a target and an
  owner.** Metrics without targets are tracked but excluded from the story (and
  reported in `excluded` so coverage gaps are honest).
- **Attainment** is value vs. target, clamped to `[0, 1]`.
- A **domain's aspect score** is the plain mean of attainment across its
  eligible, observed metrics with that aspect.
- The **overall aspect score** is the plain mean of domain aspect scores, so
  metric-heavy domains don't dominate.
- **`CompareRollups`** produces aspect deltas between periods — the raw material
  for a report's "movers" layer, with per-metric contributions retained for
  attribution.

!!! quote "Why plain means?"
    Simple means are deliberate. The first question about "Standardization: 71"
    is *what does 71 mean* — and "the average attainment of our standards
    metrics" survives the meeting.

### Coverage honesty

Two lists make the denominator explicit and keep scores honest:

- **`excluded`** — metrics that are tracked but have no target/owner, so they do
  not contribute to any score.
- **`missing`** — metrics that are rollup-eligible but were not observed this
  period.

## Maturity ladders (weakest-link)

Independently of aspect rollups, metrics can carry a `MaturityLadder` that maps
an observed value to a PRISM maturity level (1–5).

- **Capability maturity** = `min(metric rungs)`
- **Domain maturity** = `min(capability rungs)`
- **"Not tracked" is a status (N/A), not a level** — absence of measurement
  earns no rung, and N/A propagates upward.

See [Catalog Authoring](catalog-authoring.md#maturity-ladders) for how to write
ladders.

## Documents

| Document | Schema | Purpose |
| ------------- | -------------------------------- | ---------------------------------------------------------- |
| Framework | `scale-framework.schema.json` | Framework metadata, narratives, and (assembled) domains. |
| Domain | `scale-domain.schema.json` | One domain's dimensions, capabilities, metrics, narratives. |
| Assessment | `scale-assessment.schema.json` | Observed values plus journey/outlook narratives for one period. |

Catalogs are stored as `framework.json` + `domains/*.json` and assembled by
`LoadFrameworkDir`. Schemas are generated from the Go types
(`go run ./cmd/schemagen`) and embedded via the `schema` package. See the
[reference](../reference/framework.md) for the field-level formats.
