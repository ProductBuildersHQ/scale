# SCALE Roadmap

Roadmap organized by themed phases. Status is derived from member RMI statuses.

## Phase 1: Foundation (v0.1.0)

Core specification and tooling for platform maturity measurement.

| RMI | Status | Description |
|-----|--------|-------------|
| RMI-SCALE-001 | Done | Core SCALE types (Framework, Domain, Capability, Metric) |
| RMI-SCALE-002 | Done | Maturity ladder computation with weakest-link semantics |
| RMI-SCALE-003 | Done | Rollup aggregation and period comparison |
| RMI-SCALE-004 | Done | External model support (AWS, New Relic) |
| RMI-SCALE-005 | Done | HTML story report with aspect tiles and movers |
| RMI-SCALE-006 | Done | CLI (validate, report, export) |
| RMI-SCALE-007 | Done | Lit web component (`<scale-report>`) |
| RMI-SCALE-008 | Done | Reference catalog (api, observability, security) |

## Phase 2: Documentation & Integration

Authoring guides and ecosystem integration.

| RMI | Status | Description |
|-----|--------|-------------|
| RMI-SCALE-010 | Done | Catalog authoring guide |
| RMI-SCALE-011 | Done | Register SCALE as framework constant in prism-core |
| RMI-SCALE-012 | Done | Add SCALE aspect types to prism-core |
| RMI-SCALE-013 | Done | Add `scaleAspect` field to prism-maturity SLI type |
| RMI-SCALE-014 | Planned | prism-maturity SCALE rollup integration |

## Phase 3: Distribution

Binary distribution and package publishing.

| RMI | Status | Description |
|-----|--------|-------------|
| RMI-SCALE-020 | Planned | GoReleaser configuration for multi-platform binaries |
| RMI-SCALE-021 | Planned | Homebrew tap (`brew install scale`) |
| RMI-SCALE-022 | Planned | npm publish `@productbuildershq/scale-report` |

## Phase 4: Additional Domains

Expand the reference catalog with more engineering domains.

| RMI | Status | Description |
|-----|--------|-------------|
| RMI-SCALE-030 | Planned | Developer Experience (DevEx) domain |
| RMI-SCALE-031 | Planned | Data Engineering domain |
| RMI-SCALE-032 | Planned | Cost Optimization domain |
| RMI-SCALE-033 | Planned | Infrastructure domain |
| RMI-SCALE-034 | Planned | AI/ML Engineering domain |

## Phase 5: Advanced Features

Enhanced reporting and analysis capabilities.

| RMI | Status | Description |
|-----|--------|-------------|
| RMI-SCALE-040 | Planned | MkDocs documentation site with Material theme |
| RMI-SCALE-041 | Planned | Multi-period trend analysis and sparklines |
| RMI-SCALE-042 | Planned | PDF export for executive reports |
| RMI-SCALE-043 | Planned | Slack/Teams integration for report distribution |
| RMI-SCALE-044 | Planned | Assessment automation via CI pipeline |

## Phase 6: Fleet Reporting

Enterprise-scale multi-team reporting.

| RMI | Status | Description |
|-----|--------|-------------|
| RMI-SCALE-050 | Planned | Multi-team aggregation (org-level rollups) |
| RMI-SCALE-051 | Planned | Team comparison views |
| RMI-SCALE-052 | Planned | Portfolio dashboard generation |
| RMI-SCALE-053 | Planned | Time-series metrics export for BI tools |

## Completed

- v0.1.0 (2026-07-25): Phase 1 complete
- v0.1.1 (2026-07-25): Phase 2 RMI-010 through RMI-013 complete
