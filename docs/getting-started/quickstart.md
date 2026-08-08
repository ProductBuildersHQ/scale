# Quick Start

This walks through rendering a SCALE story report from the embedded reference
catalog and the example assessments shipped in the repo.

## 1. Validate the catalog

```bash
scale validate -catalog catalog
```

```text
catalog OK: 4 domains, 2 external models
```

Add an assessment to validate observations against the catalog:

```bash
scale validate -catalog catalog -assessment examples/assessments/2026-q3.json
```

## 2. Render the story report

```bash
scale report \
  -catalog catalog \
  -assessment examples/assessments/2026-q3.json \
  -prev examples/assessments/2026-q2.json \
  -o scale-report-2026-q3.html
```

Open `scale-report-2026-q3.html` in a browser. The report is a pyramid, so
hundreds of metrics read as a handful of sentences:

1. **Headline** — five aspect tiles with period-over-period deltas.
2. **What moved** — top aspect movements with per-metric attribution.
3. **Domain journeys** — each domain's lifecycle arc, aspect bars, and authored
   journey / outlook narratives.
4. **External framework lens** — where current practice sits on vendor ladders.
5. **Appendix + coverage honesty** — full metric tables plus explicit lists of
   what was tracked-but-excluded and eligible-but-unmeasured.

## 3. Render through a vendor model's ladder

The same catalog and assessment can be rendered through an external maturity
model's own levels:

```bash
scale report \
  -catalog catalog \
  -assessment examples/assessments/2026-q3.json \
  -model aws-observability-maturity \
  -o aws-observability-2026-q3.html
```

## 4. Use it as a library

```go
import scale "github.com/ProductBuildersHQ/scale"

f, err := scale.LoadFrameworkDir("catalog")
if err != nil { /* ... */ }

a, err := scale.LoadAssessmentFile("examples/assessments/2026-q3.json", f)
if err != nil { /* ... */ }

r, err := scale.ComputeRollup(f, a)
if err != nil { /* ... */ }

for _, as := range r.Aspects {
    fmt.Printf("%s: %.0f%%\n", scale.AspectDisplayName(as.Aspect), as.Score*100)
}
```

## Next steps

- Learn [the five aspects](../guide/aspects.md) and the
  [information model](../guide/information-model.md).
- [Author your own catalog](../guide/catalog-authoring.md).
