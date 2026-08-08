# Reports & Web Component

SCALE renders a framework + assessment as a story report — server-side as
self-contained HTML, or client-side via a web component.

## The HTML story report

```bash
scale report \
  -catalog catalog \
  -assessment examples/assessments/2026-q3.json \
  -prev examples/assessments/2026-q2.json \
  -o scale-report-2026-q3.html
```

The output is a single self-contained HTML file — no external assets, with light
and dark themes. It is structured as a pyramid:

1. **Headline** — five aspect tiles with period-over-period deltas.
2. **What moved** — top aspect movements with per-metric attribution
   (*"driven by APIs passing silver conformance (+19 pts)"*).
3. **Domain journeys** — each domain's lifecycle arc, aspect bars, and the
   period's authored journey / outlook narratives.
4. **External framework lens** — where current practice sits on codified vendor
   ladders.
5. **Appendix + coverage honesty** — full metric tables, plus explicit lists of
   what was tracked-but-excluded (no target/owner) and eligible-but-unmeasured.

## External maturity model lens

Third-party maturity models (AWS, New Relic, …) are codified as **data**
(`catalog/external/*.json`), not machinery: source-faithful levels with
provenance (`sourceUrl`, `retrievedAt`, `interpretation`) and an explicit
per-model PRISM crosswalk — never a hard-coded offset, because published models
differ structurally.

The same catalog and assessment can be rendered **through a vendor model's own
ladder**:

```bash
scale report \
  -catalog catalog \
  -assessment examples/assessments/2026-q3.json \
  -prev examples/assessments/2026-q2.json \
  -model aws-observability-maturity \
  -o aws-observability-2026-q3.html
```

The model report structures the page as the vendor's levels (with their PRISM
crosswalks), places your capabilities and metric evidence on each rung,
summarizes current position, and flags the **next up** level with what it
requires. Its footer states the evidence boundary: the placement is a SCALE
mapping authored in the catalog, not a vendor assessment.

## JSON IR and `HTMLFromIR`

The report is built from an intermediate representation (`ReportIR`) that holds
all computed data — aspects, movers, domains, capabilities, coverage — without
requiring the original `Framework` or `Assessment`.

- `report.BuildIR(f, a, opts)` builds the IR.
- `report.JSON(f, a, opts)` renders it as JSON.
- `report.HTMLFromIR(ir)` renders HTML from a pre-built or cached IR.

This makes it cheap to cache the computed report, serve it as JSON, or render it
in more than one format from a single computation.

## Web component

`web/` provides `<scale-report>`, a Lit web component that renders the same story
report client-side from the framework JSON IR and an assessment JSON — injectable
into a div on React, MkDocs, or plain HTML sites, and themeable via `--scale-*`
CSS custom properties.

Export the assembled framework IR for the component with:

```bash
scale export -catalog catalog -o framework.json
```

`web/src/compute.js` is a dependency-free mirror of the Go rollup / maturity
semantics, pinned to Go-verified values by `npm test`. See `web/README.md` for
embedding details.

!!! info "Exporting `framework.json`"
    `scale export` assembles `framework.json` + `domains/*.json` +
    `external/*.json` into one validated JSON IR. The repo-root `framework.json`
    is committed as a generated artifact so the assembled framework is browsable
    on GitHub — regenerate it whenever the catalog changes (see `CLAUDE.md`).
