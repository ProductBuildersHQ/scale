# CLAUDE.md — SCALE

Project-specific guidelines for Claude Code in the SCALE repo.

## Project Overview

SCALE is a maturity-measurement framework that organizes engineering metrics
under five aspects — **S**tandards, **C**onsumption, **A**utomation,
**L**everage, **E**ffectiveness — so a large metric set reads as a small number
of aspect-level storylines. The Go library computes rollups and renders reports;
the reference framework lives in `catalog/`.

## Layout

```
catalog/          # Reference framework catalog (embedded via catalog/embed.go)
  framework.json    # Framework-level metadata (id, narratives)
  domains/*.json    # One file per domain (api, observability, platform, security)
  external/*.json   # External maturity model mappings
cmd/scale/        # CLI: validate | report | export
examples/         # Example assessments
report/           # Report rollup + IR + HTML/JSON renderers
schema/           # JSON Schemas
```

## Rebuild `framework.json` on catalog changes

The repo-root `framework.json` is a **generated artifact**: the assembled,
validated bundle of `catalog/framework.json` + `catalog/domains/*.json` +
`catalog/external/*.json`, produced by `scale export`.

We commit it so the fully-assembled framework is **browsable directly on
github.com** (a single JSON file to view/link, rather than piecing together the
catalog directory).

**Whenever anything under `catalog/` changes, regenerate the root
`framework.json` in the same change so it does not drift:**

```bash
go run ./cmd/scale export -catalog catalog -o framework.json
```

The export is deterministic — regenerating with an unchanged catalog produces a
byte-identical file. Do not hand-edit root `framework.json`; edit the catalog
sources and re-export.

## Documentation Site

Docs are a MkDocs (Material) site under `docs/`, configured by `mkdocs.yml`.
Content lives in `docs/{getting-started,guide,reference,specs,releases}/`; the
existing `docs/guide/catalog-authoring.md` is the authoring how-to.

```bash
# Preview locally (http://127.0.0.1:8000)
pip install -r docs/requirements.txt
mkdocs serve

# Build (fail on broken links / nav)
mkdocs build --strict
```

The canonical site is served at `productbuildershq.com/scale` by the
ProductBuildersHQ central docs pipeline (same as sibling repos: shared theme
CSS/JS and unified navbar are pulled from `productbuildershq.com`). The in-repo
`.github/workflows/mkdocs.yaml` additionally builds on every PR and publishes to
GitHub Pages on pushes to `main`. Keep the nav in `mkdocs.yml` in sync when
adding or moving pages, and run `mkdocs build --strict` before committing.

## Common Commands

```bash
# Validate the catalog (optionally against an assessment)
go run ./cmd/scale validate -catalog catalog [-assessment examples/assessments/2026-q3.json]

# Render a report
go run ./cmd/scale report -catalog catalog -assessment <file> [-prev <file>] -o out.html

# Rebuild the browsable framework bundle (see above)
go run ./cmd/scale export -catalog catalog -o framework.json

# Build / test
go build ./...
go test ./...
```

## Commit Conventions

Follow Conventional Commits. Catalog/data edits that trigger a re-export can pair
the source change (`feat(catalog): ...`) with the regenerated bundle
(`build(catalog): regenerate framework.json export`), or bundle both when the
change is small.
