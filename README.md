# SCALE

SCALE is a machine-readable framework for telling the **platform story**: how engineering best practices are **S**tandardized, **C**onsumed, **A**utomated, **L**everaged, and made **E**ffective across an organization.

SCALE is an aggregation and narrative layer over detailed maturity metrics (such as those assessed by [PRISM](https://github.com/grokify/prism-maturity)) — not another metrics catalog. Maturity frameworks make it cheap to generate hundreds of metrics, and cheap metrics destroy the story. SCALE's purpose is to make hundreds of metrics readable as five aspect-level storylines, each backed by drill-down evidence.

## The Five Aspects

Every metric is tagged with exactly one aspect. These are frozen; changing them is a breaking change.

| Letter | Aspect | Question it answers |
| ------ | ------------- | ------------------------------------------------------------------- |
| **S** | Standards | Do executable engineering standards exist? |
| **C** | Consumption | Are teams using them (adoption) and using them correctly (conformance)? |
| **A** | Automation | How much is enforced automatically? |
| **L** | Leverage | What reuse and engineering capacity do they create? |
| **E** | Effectiveness | Did engineering and business outcomes improve? |

Consumption is the only aspect with sub-kinds, because *adoption* ("the SDK is installed") and *conformance* ("the telemetry actually passes the profile") are different facts that tell different stories.

## Information Model

```text
Framework
    └── Domain (api, observability, security, ...)
        ├── DomainDimension (ordered lifecycle stages — the story spine)
        └── Capability
            └── Metric (each tagged with exactly one SCALE aspect)
```

Two orthogonal story axes:

- **SCALE aspects** tell the horizontal platform story: *"Standardization is improving across domains."*
- **Domain dimensions** tell each discipline's journey story. Stage order is slice order — the arc is the order. For Security: Design & Planning → Development → CI/CD Pipeline → Offensive Validation → Operations.

### Narratives

`NarrativeBlock`s attach at every level — the "literate" half of a SCALE document, analogous to markdown cells in a notebook. Three kinds with different lifetimes:

| Kind | Lifetime | Home | Purpose |
| --------- | ---------- | ----------------- | ---------------------------------------------------- |
| `thesis` | Timeless | Framework catalog | Why this element matters. Carries `reviewBy` so stale prose is visibly stale. |
| `journey` | Time-bound | Assessment | Where we were and what changed this period. |
| `outlook` | Forward | Assessment | Where we are going; references prism-roadmap initiative IDs so the next cycle can compare expected vs. actual movement. |

### Rollup Semantics

Rollup rules are part of the specification, not the report generator, so every score means the same thing in every report:

- A metric participates in rollups **only when it has both a target and an owner**. Metrics without targets are tracked but excluded from the story (and reported in `excluded` so coverage gaps are honest).
- Attainment is value vs. target, clamped to `[0,1]`.
- A domain's aspect score is the plain mean of attainment across its eligible, observed metrics with that aspect.
- The overall aspect score is the plain mean of domain aspect scores, so metric-heavy domains don't dominate.
- `CompareRollups` produces aspect deltas between periods — the raw material for a report's "movers" layer, with per-metric contributions retained for attribution.

Simple means are deliberate: the first question about "Standardization: 71" is *what does 71 mean*, and "the average attainment of our standards metrics" survives the meeting.

## Documents

| Document | Schema | Purpose |
| ------------- | -------------------------------- | ---------------------------------------------------------- |
| Framework | `scale-framework.schema.json` | Framework metadata, narratives, and (assembled) domains. |
| Domain | `scale-domain.schema.json` | One domain's dimensions, capabilities, metrics, narratives. |
| Assessment | `scale-assessment.schema.json` | Observed values plus journey/outlook narratives for one period. |

Catalogs are stored as `framework.json` + `domains/*.json` and assembled by `LoadFrameworkFS`. Schemas are generated from the Go types (`go run ./cmd/schemagen`) and embedded via `schema` package.

## Reference Catalog

The embedded reference catalog (`catalog/`) seeds three domains:

- **api** (active) — anchored by [api-style-spec](https://github.com/plexusone/api-style-spec): rule coverage (S), linting enrollment (C-adoption), silver conformance pass rate (C-conformance), blocking CI validation (A), generated SDKs (L), unplanned breaking changes (E).
- **observability** (active) — OpenTelemetry/Prometheus adoption, Golden Signals and RED/USE conformance, auto-instrumentation, template dashboards, MTTR and alert-based incident detection.
- **security** (draft) — the AI-era security lifecycle as the worked example of a domain dimension: AI Threat Modeling → Defensive Hygiene → Automated Posture → AI Pentesting & Red Teams → AI SOC/CSIRT/PSIRT.

```go
import "github.com/ProductBuildersHQ/scale/catalog"

f, err := catalog.Default()
```

## Usage

```go
import scale "github.com/ProductBuildersHQ/scale"

f, err := scale.LoadFrameworkDir("catalog")
if err != nil { /* ... */ }

a, err := scale.LoadAssessmentFile("assessments/2026-Q3.json", f)
if err != nil { /* ... */ }

r, err := scale.ComputeRollup(f, a)
if err != nil { /* ... */ }

for _, as := range r.Aspects {
    fmt.Printf("%s: %.0f%%\n", scale.AspectDisplayName(as.Aspect), as.Score*100)
}
```

## Viewing SCALE Metrics: the Story Report

The `scale` CLI renders a framework + assessment as a self-contained HTML story report — no external assets, light and dark themes:

```bash
go run ./cmd/scale report \
  -catalog catalog \
  -assessment examples/assessments/2026-q3.json \
  -prev examples/assessments/2026-q2.json \
  -o examples/scale-report-2026-q3.html
```

The report is a pyramid, so hundreds of metrics read as a handful of sentences:

1. **Headline** — five aspect tiles with period-over-period deltas.
2. **What moved** — top aspect movements with per-metric attribution ("driven by APIs passing silver conformance (+19 pts)").
3. **Domain journeys** — each domain's lifecycle arc, aspect bars, and the period's authored journey/outlook narratives (with prism-roadmap initiative references).
4. **External framework lens** — where current practice sits on codified vendor ladders (see below).
5. **Appendix + coverage honesty** — full metric tables, plus explicit lists of what was tracked-but-excluded (no target/owner) and eligible-but-unmeasured.

Why platform teams should do this: platform investment is invisible precisely because it succeeds — nothing breaks, nothing ships late, and at budget time there is no story. The report turns adoption/conformance/automation evidence into the narrative leadership actually reads: *what moved, why it moved, what it unlocks next*.

### External Maturity Models (AWS, New Relic, ...)

Third-party maturity models are codified as **data** (`catalog/external/*.json`), not machinery: source-faithful levels with provenance (`sourceUrl`, `retrievedAt`, `interpretation`) and an explicit per-model PRISM crosswalk — never a hard-coded offset, because published models differ structurally (AWS uses four stages, Foundational → Proactive; New Relic uses a Level 0 foundation plus a reactive → proactive → mastery progression per value driver).

Capabilities reference a model through ordinary `frameworks` mappings (`framework` = model ID, `reference` = level ID), and the report renders the ladder with a "current practice" marker — "we're at AWS Stage 2 of 4; Stage 3 requires AI/ML anomaly detection" is a roadmap sentence that writes itself.

The same catalog and assessment can also be rendered **through a vendor model's own ladder** — the "frameworks are projections" principle as a report:

```bash
go run ./cmd/scale report \
  -catalog catalog \
  -assessment examples/assessments/2026-q3.json \
  -prev examples/assessments/2026-q2.json \
  -model aws-observability-maturity \
  -o examples/aws-observability-report-2026-q3.html
```

The model report structures the page as the vendor's levels (with their PRISM crosswalks), places your capabilities and metric evidence on each rung, summarizes current position ("Current practice: Intermediate monitoring"), flags the **next up** level with what it requires, and keeps the SCALE aspect bars and journey narratives for the model's domain. Its footer states the evidence boundary: the placement is a SCALE mapping authored in the catalog, not a vendor assessment.

Note: scale deliberately does **not** import prism-capability or prism-maturity. The dependency arrow points one way — PRISM (and any other assessment platform) consumes SCALE. That keeps the spec implementation-neutral and leaves PRISM free to import scale for its own richer, fleet-level reporting.

## Web Component

`web/` provides `<scale-report>`, a Lit web component rendering the same story report client-side from the framework JSON IR (`scale export -catalog <dir> -o framework.json`) and assessment JSON — injectable into a div on React, MkDocs, or plain HTML sites, themeable via `--scale-*` CSS custom properties. `web/src/compute.js` is a dependency-free mirror of the Go rollup/maturity semantics, pinned to Go-verified values by `npm test`. See `web/README.md`.

## Relationship to PRISM and External Frameworks

- **Capabilities are canonical; frameworks are projections.** A metric maps to external frameworks (DORA, NIST CSF, OWASP, ...) via `frameworks` mappings (shared `prism-core` `FrameworkMapping` type) without belonging to them — evidence is collected once and interpreted many times.
- **PRISM assesses; SCALE narrates.** Capabilities link to prism-capability entries and prism-maturity domains/SLIs via `prism` refs, including expected evidence per maturity level (M1–M5).
- Compliance frameworks tell the auditor's story; SCALE's domain dimensions tell the transformation story. Same capabilities, same evidence, two projections.

## Status

v0.x — the information model is expected to evolve alongside its first consumers (PRISM reporting, api-style-spec evidence pipelines, observability profiles).
