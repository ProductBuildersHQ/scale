// SCALE rollup and maturity computation — a JavaScript mirror of the Go
// reference implementation (rollup.go, maturity.go). The semantics are
// defined by the SCALE specification; if this file and the Go implementation
// disagree, the Go implementation wins and this file has a bug. The node
// test (test/compute.test.mjs) pins this module to values verified by the
// Go test suite.
//
// This module is dependency-free so it can be imported by the <scale-report>
// Lit component, a React app, or node without pulling in a renderer.

export const ASPECTS = [
  'standards',
  'consumption',
  'automation',
  'leverage',
  'effectiveness',
];

const ASPECT_META = {
  standards: { letter: 'S', name: 'Standards' },
  consumption: { letter: 'C', name: 'Consumption' },
  automation: { letter: 'A', name: 'Automation' },
  leverage: { letter: 'L', name: 'Leverage' },
  effectiveness: { letter: 'E', name: 'Effectiveness' },
};

export function aspectLetter(aspect) {
  return ASPECT_META[aspect]?.letter ?? '?';
}

export function aspectDisplayName(aspect) {
  return ASPECT_META[aspect]?.name ?? aspect;
}

function isLowerBetter(metric) {
  return metric.direction === 'lower_is_better';
}

export function rollupEligible(metric) {
  return Boolean(metric.target && metric.owner);
}

export function observation(assessment, metricId) {
  return (assessment.observations ?? []).find((o) => o.metricId === metricId) ?? null;
}

const clamp01 = (v) => Math.min(1, Math.max(0, v));

// attainment normalizes an observed value against the metric's target to
// [0,1]. Higher-is-better: value/target. Lower-is-better: 1 at or under
// target, otherwise target/value.
export function attainment(metric, value) {
  const target = metric.target?.value;
  if (target === undefined || target === null) {
    throw new Error(`metric ${metric.id}: attainment requires a target`);
  }
  if (isLowerBetter(metric)) {
    if (target < 0) throw new Error(`metric ${metric.id}: lower-is-better target must be non-negative`);
    if (value <= target) return 1;
    return clamp01(target / value);
  }
  if (target <= 0) throw new Error(`metric ${metric.id}: higher-is-better target must be positive`);
  return clamp01(value / target);
}

// metricMaturity returns the ladder rung a metric holds, or null when the
// metric has no ladder, there is no observation (not tracked — N/A), or the
// value satisfies no rung (below the ladder). A rung without a threshold is
// satisfied by the observation's existence; exclusive rungs compare strictly.
export function metricMaturity(metric, obs) {
  const rungs = metric.maturity?.levels;
  if (!rungs?.length || !obs) return null;
  const lower = isLowerBetter(metric);
  let best = null;
  for (const r of rungs) {
    if (r.threshold === undefined || r.threshold === null) {
      best = r;
    } else if (!lower && (r.exclusive ? obs.value > r.threshold : obs.value >= r.threshold)) {
      best = r;
    } else if (lower && (r.exclusive ? obs.value < r.threshold : obs.value <= r.threshold)) {
      best = r;
    }
  }
  return best;
}

// capabilityMaturity is the weakest link: the lowest rung across the
// capability's laddered metrics, or null when there are no laddered metrics
// or any laddered metric has no maturity (not tracked / below the ladder).
export function capabilityMaturity(capability, assessment) {
  let min = null;
  for (const m of capability.metrics ?? []) {
    if (!m.maturity) continue;
    const rung = metricMaturity(m, observation(assessment, m.id));
    if (!rung) return null;
    if (!min || rung.level < min.level) min = rung;
  }
  return min;
}

export function hasLadderedMetrics(capability) {
  return (capability.metrics ?? []).some((m) => m.maturity);
}

// domainMaturity is the weakest link across the domain's laddered
// capabilities. Maturity claims are conjunctive: a domain cannot claim a
// level its weakest capability has not earned, and a single high metric
// never lifts the domain. Returns null when no capability carries ladders
// (no claim) or when any laddered capability is N/A.
export function domainMaturity(domain, assessment) {
  let min = null;
  for (const c of domain.capabilities ?? []) {
    if (!hasLadderedMetrics(c)) continue;
    const rung = capabilityMaturity(c, assessment);
    if (!rung) return null;
    if (!min || rung.level < min.level) min = rung;
  }
  return min;
}

export function findMetric(framework, metricId) {
  for (const d of framework.domains ?? []) {
    for (const c of d.capabilities ?? []) {
      for (const m of c.metrics ?? []) {
        if (m.id === metricId) return { domain: d, capability: c, metric: m };
      }
    }
  }
  return null;
}

// computeRollup mirrors Go ComputeRollup: attainment per eligible observed
// metric; domain aspect score = plain mean; overall aspect score = plain
// mean of domain scores so metric-heavy domains don't dominate.
export function computeRollup(framework, assessment) {
  const rollup = {
    frameworkId: framework.id,
    period: assessment.period,
    aspects: [],
    domains: [],
    excluded: [],
    missing: [],
  };
  const overall = new Map(); // aspect -> [domain scores]

  for (const d of framework.domains ?? []) {
    const byAspect = new Map();
    for (const c of d.capabilities ?? []) {
      for (const m of c.metrics ?? []) {
        const obs = observation(assessment, m.id);
        if (!rollupEligible(m)) {
          if (obs) rollup.excluded.push(m.id);
          continue;
        }
        if (!obs) {
          rollup.missing.push(m.id);
          continue;
        }
        const contrib = {
          domainId: d.id,
          capabilityId: c.id,
          metricId: m.id,
          aspect: m.aspect,
          value: obs.value,
          attainment: attainment(m, obs.value),
        };
        if (!byAspect.has(m.aspect)) byAspect.set(m.aspect, []);
        byAspect.get(m.aspect).push(contrib);
      }
    }
    if (byAspect.size === 0) continue;
    const domainRollup = { domainId: d.id, aspects: [] };
    for (const aspect of ASPECTS) {
      const contribs = byAspect.get(aspect);
      if (!contribs?.length) continue;
      const score = contribs.reduce((s, c) => s + c.attainment, 0) / contribs.length;
      domainRollup.aspects.push({
        aspect,
        score,
        metricCount: contribs.length,
        contributions: contribs,
      });
      if (!overall.has(aspect)) overall.set(aspect, []);
      overall.get(aspect).push(score);
    }
    rollup.domains.push(domainRollup);
  }

  for (const aspect of ASPECTS) {
    const scores = overall.get(aspect);
    if (!scores?.length) continue;
    let metricCount = 0;
    for (const dr of rollup.domains) {
      for (const as of dr.aspects) {
        if (as.aspect === aspect) metricCount += as.metricCount;
      }
    }
    rollup.aspects.push({
      aspect,
      score: scores.reduce((s, v) => s + v, 0) / scores.length,
      metricCount,
    });
  }

  rollup.excluded.sort();
  rollup.missing.sort();
  return rollup;
}

export function aspectScoreFor(rollup, aspect) {
  return rollup.aspects.find((a) => a.aspect === aspect) ?? null;
}

export function domainRollupFor(rollup, domainId) {
  return rollup.domains.find((d) => d.domainId === domainId) ?? null;
}

// compareRollups mirrors Go CompareRollups: deltas for the cross-domain
// scores (domainId '') and every domain, sorted by absolute movement.
// Aspects absent from one side are skipped — coverage change, not movement.
export function compareRollups(prev, curr) {
  const deltas = [];
  const compare = (domainId, prevAspects, currAspects) => {
    for (const p of prevAspects) {
      for (const c of currAspects) {
        if (p.aspect === c.aspect) {
          deltas.push({
            domainId,
            aspect: p.aspect,
            prev: p.score,
            curr: c.score,
            delta: c.score - p.score,
          });
        }
      }
    }
  };
  compare('', prev.aspects, curr.aspects);
  const domainIds = new Set([
    ...prev.domains.map((d) => d.domainId),
    ...curr.domains.map((d) => d.domainId),
  ]);
  for (const domainId of domainIds) {
    compare(
      domainId,
      domainRollupFor(prev, domainId)?.aspects ?? [],
      domainRollupFor(curr, domainId)?.aspects ?? [],
    );
  }
  deltas.sort((a, b) => Math.abs(b.delta) - Math.abs(a.delta));
  return deltas;
}

// buildMovers mirrors the Go report's movers: domain-level deltas with the
// top per-metric attainment contributions for attribution.
export function buildMovers(framework, rollup, prevRollup, limit = 6) {
  const prevAttain = new Map();
  for (const dr of prevRollup.domains) {
    for (const as of dr.aspects) {
      for (const c of as.contributions ?? []) prevAttain.set(c.metricId, c.attainment);
    }
  }
  const movers = [];
  for (const delta of compareRollups(prevRollup, rollup)) {
    if (!delta.domainId || delta.delta === 0) continue;
    const domain = (framework.domains ?? []).find((d) => d.id === delta.domainId);
    if (!domain) continue;
    const contributions = [];
    const dr = domainRollupFor(rollup, delta.domainId);
    const as = dr?.aspects.find((a) => a.aspect === delta.aspect);
    if (as) {
      const cds = (as.contributions ?? [])
        .filter((c) => prevAttain.has(c.metricId))
        .map((c) => ({
          name: findMetric(framework, c.metricId)?.metric.name ?? c.metricId,
          delta: (c.attainment - prevAttain.get(c.metricId)) * 100,
        }))
        .sort((a, b) => Math.abs(b.delta) - Math.abs(a.delta));
      for (const cd of cds.slice(0, 2)) {
        if (cd.delta === 0) break;
        contributions.push(`${cd.name} (${cd.delta > 0 ? '+' : ''}${Math.round(cd.delta)} pts)`);
      }
    }
    movers.push({
      domain: domain.name,
      aspect: aspectDisplayName(delta.aspect),
      prev: delta.prev * 100,
      curr: delta.curr * 100,
      delta: delta.delta * 100,
      contributions,
    });
    if (movers.length === limit) break;
  }
  return movers;
}
