// Pins compute.js to values verified by the Go reference implementation's
// test suite (rollup_test.go, maturity_test.go). If these fail, compute.js
// has drifted from the SCALE specification.
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { test } from 'node:test';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

import {
  attainment,
  buildMovers,
  capabilityMaturity,
  computeRollup,
  domainRollupFor,
  metricMaturity,
  observation,
} from '../src/compute.js';

const repoRoot = join(dirname(fileURLToPath(import.meta.url)), '..', '..');

const coverageLadder = () => ({
  levels: [
    { level: 1, name: 'Tracked' },
    { level: 2, name: '>0%', threshold: 0, exclusive: true },
    { level: 3, name: '>=50%', threshold: 50 },
    { level: 4, name: '>=80%', threshold: 80 },
    { level: 5, name: '100%', threshold: 100 },
  ],
});

// Mirror of Go testFramework().
const testFramework = () => ({
  id: 'scale',
  name: 'SCALE',
  domains: [
    {
      id: 'api',
      name: 'API',
      capabilities: [
        {
          id: 'style',
          name: 'Style',
          metrics: [
            { id: 'api.c.adoption', name: 'Adoption', aspect: 'consumption', consumptionKind: 'adoption', target: { value: 100 }, owner: 'api-platform' },
            { id: 'api.c.conformance', name: 'Conformance', aspect: 'consumption', consumptionKind: 'conformance', target: { value: 100 }, owner: 'api-platform' },
            { id: 'api.e.breaking', name: 'Breaking changes', aspect: 'effectiveness', direction: 'lower_is_better', target: { value: 2 }, owner: 'api-platform' },
            { id: 'api.untracked', name: 'No target', aspect: 'leverage' },
          ],
        },
      ],
    },
    {
      id: 'obs',
      name: 'Observability',
      capabilities: [
        {
          id: 'otel',
          name: 'OTel',
          metrics: [
            { id: 'obs.c.adoption', name: 'Adoption', aspect: 'consumption', consumptionKind: 'adoption', target: { value: 100 }, owner: 'obs-platform' },
            { id: 'obs.missing', name: 'Unmeasured', aspect: 'automation', target: { value: 100 }, owner: 'obs-platform' },
          ],
        },
      ],
    },
  ],
});

test('attainment matches Go semantics', () => {
  const higher = { id: 'h', target: { value: 100 } };
  const lower = { id: 'l', direction: 'lower_is_better', target: { value: 2 } };
  assert.equal(attainment(higher, 100), 1);
  assert.equal(attainment(higher, 120), 1);
  assert.ok(Math.abs(attainment(higher, 80) - 0.8) < 1e-9);
  assert.equal(attainment(higher, 0), 0);
  assert.equal(attainment(lower, 2), 1);
  assert.equal(attainment(lower, 0), 1);
  assert.equal(attainment(lower, 4), 0.5);
});

test('computeRollup matches Go TestComputeRollup', () => {
  const f = testFramework();
  const a = {
    frameworkId: 'scale',
    period: '2026-Q3',
    observations: [
      { metricId: 'api.c.adoption', value: 90 },
      { metricId: 'api.c.conformance', value: 70 },
      { metricId: 'api.e.breaking', value: 4 },
      { metricId: 'api.untracked', value: 55 },
      { metricId: 'obs.c.adoption', value: 80 },
    ],
  };
  const r = computeRollup(f, a);
  const apiConsumption = domainRollupFor(r, 'api').aspects.find((x) => x.aspect === 'consumption');
  assert.ok(Math.abs(apiConsumption.score - 0.8) < 1e-9);
  assert.equal(apiConsumption.metricCount, 2);
  const overall = r.aspects.find((x) => x.aspect === 'consumption');
  assert.ok(Math.abs(overall.score - 0.8) < 1e-9);
  const effectiveness = r.aspects.find((x) => x.aspect === 'effectiveness');
  assert.ok(Math.abs(effectiveness.score - 0.5) < 1e-9);
  assert.deepEqual(r.excluded, ['api.untracked']);
  assert.deepEqual(r.missing, ['obs.missing']);
});

test('buildMovers matches Go CompareRollups delta', () => {
  const f = testFramework();
  const prev = computeRollup(f, {
    frameworkId: 'scale',
    period: '2026-Q2',
    observations: [
      { metricId: 'api.c.adoption', value: 70 },
      { metricId: 'api.c.conformance', value: 50 },
    ],
  });
  const curr = computeRollup(f, {
    frameworkId: 'scale',
    period: '2026-Q3',
    observations: [
      { metricId: 'api.c.adoption', value: 90 },
      { metricId: 'api.c.conformance', value: 70 },
    ],
  });
  const movers = buildMovers(f, curr, prev);
  const m = movers.find((x) => x.domain === 'API' && x.aspect === 'Consumption');
  assert.ok(m, 'expected API consumption mover');
  assert.ok(Math.abs(m.delta - 20) < 1e-9);
});

test('metricMaturity matches Go semantics (N/A, tracked floor, exclusive, lower-is-better)', () => {
  const m = { id: 'x', maturity: coverageLadder() };
  assert.equal(metricMaturity(m, null), null); // not tracked → N/A
  assert.equal(metricMaturity(m, { value: 0 }).level, 1); // infrastructure in place
  assert.equal(metricMaturity(m, { value: 0.5 }).level, 2); // one reporter proves it
  assert.equal(metricMaturity(m, { value: 45 }).level, 2);
  assert.equal(metricMaturity(m, { value: 50 }).level, 3);
  assert.equal(metricMaturity(m, { value: 85 }).level, 4);
  assert.equal(metricMaturity(m, { value: 100 }).level, 5);

  const mttr = {
    id: 'mttr',
    direction: 'lower_is_better',
    maturity: {
      levels: [
        { level: 2, name: 'Tracked' },
        { level: 3, name: '<=60m', threshold: 60 },
        { level: 4, name: '<=30m', threshold: 30 },
      ],
    },
  };
  assert.equal(metricMaturity(mttr, { value: 90 }).level, 2);
  assert.equal(metricMaturity(mttr, { value: 45 }).level, 3);
  assert.equal(metricMaturity(mttr, { value: 20 }).level, 4);
});

test('capabilityMaturity is the weakest link and N/A on untracked', () => {
  const c = {
    id: 'otel',
    metrics: [
      { id: 'm.metrics', aspect: 'consumption', consumptionKind: 'adoption', maturity: coverageLadder() },
      { id: 'm.traces', aspect: 'consumption', consumptionKind: 'adoption', maturity: coverageLadder() },
      { id: 'm.unladdered', aspect: 'leverage' },
    ],
  };
  const both = {
    observations: [
      { metricId: 'm.metrics', value: 85 },
      { metricId: 'm.traces', value: 45 },
    ],
  };
  assert.equal(capabilityMaturity(c, both).level, 2);
  const oneUntracked = { observations: [{ metricId: 'm.metrics', value: 100 }] };
  assert.equal(capabilityMaturity(c, oneUntracked), null);
});

test('domainMaturity is the weakest capability and N/A propagates', async () => {
  const { domainMaturity } = await import('../src/compute.js');
  const d = {
    id: 'obs',
    capabilities: [
      { id: 'coverage', metrics: [{ id: 'm.a', aspect: 'consumption', consumptionKind: 'adoption', maturity: coverageLadder() }] },
      { id: 'signals', metrics: [{ id: 'm.b', aspect: 'effectiveness', maturity: coverageLadder() }] },
      { id: 'unladdered', metrics: [{ id: 'm.c', aspect: 'leverage' }] },
    ],
  };
  // coverage L5, signals L3 → domain L3: one M5 metric never lifts the domain.
  const both = { observations: [{ metricId: 'm.a', value: 100 }, { metricId: 'm.b', value: 62 }] };
  assert.equal(domainMaturity(d, both).level, 3);
  // An untracked laddered capability makes the domain N/A.
  const partial = { observations: [{ metricId: 'm.a', value: 100 }] };
  assert.equal(domainMaturity(d, partial), null);
});

test('end-to-end against the repo example catalog', () => {
  const f = JSON.parse(readFileSync(join(repoRoot, 'examples', 'framework-ir.json'), 'utf8'));
  const q3 = JSON.parse(readFileSync(join(repoRoot, 'examples', 'assessments', '2026-q3.json'), 'utf8'));
  const r = computeRollup(f, q3);
  assert.equal(r.frameworkId, 'scale');
  assert.ok(r.aspects.length >= 4);
  assert.deepEqual(r.missing, ['sec.automation.pipeline-gates', 'sec.effectiveness.incident-mttr']);
  const obs = observation(q3, 'api.consumption.silver-pass-rate');
  assert.equal(obs.numerator, 81);
});
