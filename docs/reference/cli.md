# CLI Reference

The `scale` binary has three subcommands.

```text
scale validate -catalog <dir> [-assessment <file>]
scale report   -catalog <dir> -assessment <file> [-prev <file>] [-model <id>] -o <out.html>
scale export   -catalog <dir> -o <framework.json>
```

## `scale validate`

Loads and validates a catalog, optionally against an assessment.

| Flag | Default | Purpose |
| ---- | ------- | ------- |
| `-catalog` | `catalog` | Catalog directory (`framework.json`, `domains/`, `external/`). |
| `-assessment` | — | Optional assessment file to validate against the catalog. |

```bash
scale validate -catalog catalog -assessment examples/assessments/2026-q3.json
```

Common validation errors: duplicate metric IDs, invalid aspect values, maturity
levels not 1–5 or not ascending, and missing required fields.

## `scale report`

Renders the SCALE aspect story report as a self-contained HTML file.

| Flag | Default | Purpose |
| ---- | ------- | ------- |
| `-catalog` | `catalog` | Catalog directory. |
| `-assessment` | — | Assessment file to report on (required). |
| `-prev` | — | Prior-period assessment, enabling deltas and the "movers" layer. |
| `-model` | — | Render through an external model's ladder (e.g. `aws-observability-maturity`). |
| `-o` | — | Output HTML path. |

```bash
# Aspect story report with period-over-period deltas
scale report -catalog catalog \
  -assessment examples/assessments/2026-q3.json \
  -prev examples/assessments/2026-q2.json \
  -o scale-report-2026-q3.html

# External maturity-model lens
scale report -catalog catalog \
  -assessment examples/assessments/2026-q3.json \
  -model aws-observability-maturity \
  -o aws-observability-2026-q3.html
```

## `scale export`

Assembles `framework.json` + `domains/*.json` + `external/*.json` into a single
validated framework JSON IR — e.g. for the `<scale-report>` web component.

| Flag | Default | Purpose |
| ---- | ------- | ------- |
| `-catalog` | `catalog` | Catalog directory. |
| `-o` | `framework.json` | Output JSON file. |

```bash
scale export -catalog catalog -o framework.json
```

!!! info "Regenerate the committed bundle"
    The repo-root `framework.json` is a committed, generated artifact so the
    assembled framework is browsable on GitHub. Regenerate it with this command
    whenever anything under `catalog/` changes — the export is deterministic, so
    an unchanged catalog produces a byte-identical file.
