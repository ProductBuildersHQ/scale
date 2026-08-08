# SCALE

SCALE is a machine-readable framework for telling the **platform story**: how
engineering best practices are **S**tandardized, **C**onsumed, **A**utomated,
**L**everaged, and made **E**ffective across an organization.

SCALE is an aggregation and narrative layer over detailed maturity metrics (such
as those assessed by [PRISM](https://github.com/grokify/prism-maturity)) — not
another metrics catalog. Maturity frameworks make it cheap to generate hundreds
of metrics, and cheap metrics destroy the story. SCALE's purpose is to make
hundreds of metrics readable as **five aspect-level storylines**, each backed by
drill-down evidence.

## The Five Aspects

Every metric is tagged with exactly one aspect. These are frozen; changing them
is a breaking change.

| Letter | Aspect | Question it answers |
| ------ | ------------- | ------------------------------------------------------------------- |
| **S** | Standards | Do executable engineering standards exist? |
| **C** | Consumption | Are teams using them (adoption) and using them correctly (conformance)? |
| **A** | Automation | How much is enforced automatically? |
| **L** | Leverage | What reuse and engineering capacity do they create? |
| **E** | Effectiveness | Did engineering and business outcomes improve? |

## Why platform teams should do this

Platform investment is invisible precisely because it succeeds — nothing breaks,
nothing ships late, and at budget time there is no story. The SCALE report turns
adoption / conformance / automation evidence into the narrative leadership
actually reads: *what moved, why it moved, what it unlocks next*.

## Where to go next

- **[Quick Start](getting-started/quickstart.md)** — render your first story
  report in a few minutes.
- **[The Five Aspects](guide/aspects.md)** — the conceptual core of the
  framework.
- **[Information Model](guide/information-model.md)** — Framework → Domain →
  Capability → Metric, narratives, and rollup semantics.
- **[Catalog Authoring](guide/catalog-authoring.md)** — write domains, metrics,
  maturity ladders, and external models.
- **[Reference](reference/framework.md)** — the Framework, Domain, and
  Assessment file formats.

## Relationship to PRISM and external frameworks

- **Capabilities are canonical; frameworks are projections.** A metric maps to
  external frameworks (DORA, NIST CSF, OWASP, …) without belonging to them —
  evidence is collected once and interpreted many times.
- **PRISM assesses; SCALE narrates.** The dependency arrow points one way: PRISM
  (and any other assessment platform) consumes SCALE, keeping the spec
  implementation-neutral.

!!! note "Status"
    v0.x — the information model is expected to evolve alongside its first
    consumers (PRISM reporting, api-style-spec evidence pipelines, observability
    profiles).
