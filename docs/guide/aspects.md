# The Five Aspects

SCALE organizes every metric under one of five aspects, so the platform story
can be told in five sentences — each backed by drill-down evidence. Every metric
is tagged with **exactly one** aspect.

!!! warning "The aspects are frozen"
    The five aspects are part of the specification. Adding, removing, or
    redefining an aspect is a **breaking change** to the framework.

| Letter | Aspect | Question it answers | Example metrics |
| ------ | ------------- | -------------------------------------------- | ------------------------------- |
| **S** | `standards` | Do executable engineering standards exist? | Rule coverage, schema completeness |
| **C** | `consumption` | Are teams using them, and using them correctly? | SDK adoption, conformance pass rate |
| **A** | `automation` | How much is enforced automatically? | CI gate coverage, auto-remediation |
| **L** | `leverage` | What reuse and engineering capacity do they create? | Generated SDKs, template reuse |
| **E** | `effectiveness` | Did engineering and business outcomes improve? | MTTR, breaking changes, incidents |

## Consumption has two sub-kinds

Consumption is the only aspect with sub-kinds, because *adoption* and
*conformance* are different facts that tell different stories:

- **`adoption`** — "the SDK is installed."
- **`conformance`** — "the telemetry actually passes the profile."

A domain can have high adoption and low conformance (widely installed, poorly
configured) — a signal you would lose if the two collapsed into one number.
Metrics tagged `consumption` therefore **must** set `consumptionKind`.

## Two orthogonal story axes

SCALE tells two stories at once:

- **SCALE aspects** tell the horizontal platform story across domains:
  *"Standardization is improving everywhere."*
- **Domain dimensions** tell each discipline's own journey. Stage order is slice
  order — the arc *is* the order. For Security: Design & Planning → Development →
  CI/CD Pipeline → Offensive Validation → Operations.

The aspects give leadership a fixed vocabulary; the dimensions give each team its
own spine. See the [information model](information-model.md) for how these fit
together with capabilities, metrics, and narratives.
