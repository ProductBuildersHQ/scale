// Package scale defines the SCALE framework: a machine-readable specification
// for telling the platform story across horizontal engineering domains.
//
// SCALE measures how engineering best practices are Standardized, Consumed,
// Automated, Leveraged, and made Effective across an organization. It is an
// aggregation and narrative layer over detailed maturity metrics (such as
// those produced by PRISM), not another metrics catalog: its purpose is to
// make hundreds of metrics readable as a small number of aspect-level
// storylines such as "Standardization is improving."
//
// The information model is:
//
//	Framework
//	    └── Domain (api, observability, security, ...)
//	        ├── DomainDimension (ordered lifecycle stages — the story spine)
//	        └── Capability
//	            └── Metric (each tagged with exactly one SCALE aspect)
//
// NarrativeBlocks attach at every level and come in three kinds with
// different lifetimes: thesis (timeless — why this matters), journey
// (time-bound — what changed this period), and outlook (forward — where we
// are going, linked to prism-roadmap initiatives).
//
// Rollup semantics are part of the specification, not the report generator:
// a metric participates in aspect rollups only when it has both a target and
// an owner, attainment is normalized to [0,1], and aspect scores are simple
// means so that every score remains explainable.
package scale
