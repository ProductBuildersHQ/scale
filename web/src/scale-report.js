// <scale-report> — a Lit web component rendering the SCALE platformization
// story report from the framework JSON IR (produced by `scale export`) and
// assessment JSON documents. Computation (rollups, maturity ladders, movers)
// lives in ./compute.js, which mirrors the Go reference implementation.
//
// Usage (properties, e.g. from React):
//   const el = document.querySelector('scale-report');
//   el.framework = frameworkIR; el.assessment = q3; el.prev = q2;
//
// Usage (attributes, e.g. from MkDocs / plain HTML):
//   <scale-report src-framework="framework.json"
//                 src-assessment="assessments/2026-q3.json"
//                 src-prev="assessments/2026-q2.json"></scale-report>
//
// Theming: colors are exposed as --scale-* custom properties on the host and
// follow prefers-color-scheme by default; override any of them from page CSS
// (e.g. scale-report { --scale-series: #7c3aed; }).

import { LitElement, html, css, nothing } from 'lit';

import {
  aspectLetter,
  aspectDisplayName,
  attainment,
  buildMovers,
  capabilityMaturity,
  computeRollup,
  domainMaturity,
  domainRollupFor,
  hasLadderedMetrics,
  metricMaturity,
  observation,
  rollupEligible,
} from './compute.js';

const pct = (v) => `${Math.round(v)}`;
const pts = (v) => `${v > 0 ? '+' : ''}${Math.round(v)}`;

function trimFloat(v) {
  const s = v.toFixed(1);
  return s.endsWith('.0') ? s.slice(0, -2) : s;
}

function formatValue(metric, value, numerator, denominator) {
  let base = trimFloat(value);
  if (metric.unit === 'percent') base += '%';
  else if (metric.unit) base += ` ${metric.unit}`;
  if (numerator != null && denominator != null) {
    return `${numerator}/${denominator} (${base})`;
  }
  return base;
}

export class ScaleReport extends LitElement {
  static properties = {
    framework: { attribute: false },
    assessment: { attribute: false },
    prev: { attribute: false },
    srcFramework: { type: String, attribute: 'src-framework' },
    srcAssessment: { type: String, attribute: 'src-assessment' },
    srcPrev: { type: String, attribute: 'src-prev' },
    _error: { state: true },
  };

  static styles = css`
    :host {
      display: block;
      color-scheme: light dark;
      --_page: var(--scale-page, #f9f9f7);
      --_surface: var(--scale-surface, #fcfcfb);
      --_ink-1: var(--scale-ink-1, #0b0b0b);
      --_ink-2: var(--scale-ink-2, #52514e);
      --_muted: var(--scale-muted, #898781);
      --_grid: var(--scale-grid, #e1e0d9);
      --_baseline: var(--scale-baseline, #c3c2b7);
      --_series: var(--scale-series, #2a78d6);
      --_good: var(--scale-good, #006300);
      --_bad: var(--scale-bad, #d03b3b);
      --_border: var(--scale-border, rgba(11, 11, 11, 0.1));
      background: var(--_page);
      color: var(--_ink-1);
      font: 15px/1.55 system-ui, -apple-system, 'Segoe UI', sans-serif;
    }
    @media (prefers-color-scheme: dark) {
      :host {
        --_page: var(--scale-page, #0d0d0d);
        --_surface: var(--scale-surface, #1a1a19);
        --_ink-1: var(--scale-ink-1, #ffffff);
        --_ink-2: var(--scale-ink-2, #c3c2b7);
        --_muted: var(--scale-muted, #898781);
        --_grid: var(--scale-grid, #2c2c2a);
        --_baseline: var(--scale-baseline, #383835);
        --_series: var(--scale-series, #3987e5);
        --_good: var(--scale-good, #0ca30c);
        --_bad: var(--scale-bad, #d03b3b);
        --_border: var(--scale-border, rgba(255, 255, 255, 0.1));
      }
    }
    * { box-sizing: border-box; }
    main { max-width: 920px; margin: 0 auto; padding: 24px 20px 48px; }
    h1 { font-size: 26px; margin: 0 0 2px; }
    h2 { font-size: 19px; margin: 40px 0 12px; }
    h3 { font-size: 16px; margin: 24px 0 8px; }
    h4 { margin: 0 0 4px; font-size: 14px; }
    .sub { color: var(--_ink-2); margin: 0 0 4px; }
    .meta { color: var(--_muted); font-size: 13px; }
    section.card {
      background: var(--_surface);
      border: 1px solid var(--_border);
      border-radius: 10px;
      padding: 20px 22px;
      margin-top: 16px;
    }
    .narrative { border-left: 3px solid var(--_grid); padding: 2px 0 2px 14px; margin: 12px 0; }
    .narrative p { margin: 0; color: var(--_ink-2); }
    .narrative .meta { margin-top: 6px; }
    .tiles { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 12px; margin-top: 16px; }
    .tile { background: var(--_surface); border: 1px solid var(--_border); border-radius: 10px; padding: 14px 16px; }
    .tile .label { font-size: 13px; color: var(--_ink-2); }
    .tile .label b { color: var(--_ink-1); }
    .tile .value { font-size: 30px; font-weight: 650; margin-top: 2px; }
    .tile .value small { font-size: 15px; font-weight: 400; color: var(--_muted); }
    .delta { font-size: 13px; font-weight: 600; }
    .delta.up { color: var(--_good); }
    .delta.down { color: var(--_bad); }
    .delta.flat { color: var(--_muted); font-weight: 400; }
    .none { color: var(--_muted); }
    .barrow { display: grid; grid-template-columns: 130px 1fr 90px; align-items: center; gap: 10px; margin: 7px 0; }
    .barrow .name { font-size: 13px; color: var(--_ink-2); }
    .track { background: var(--_grid); border-radius: 4px; height: 10px; }
    .fill { background: var(--_series); border-radius: 0 4px 4px 0; height: 10px; min-width: 2px; }
    .barval { font-size: 13px; font-variant-numeric: tabular-nums; }
    .stages { display: flex; flex-wrap: wrap; align-items: center; gap: 6px; margin: 10px 0 4px; }
    .stage {
      border: 1px solid var(--_baseline);
      border-radius: 999px;
      padding: 3px 12px;
      font-size: 13px;
      color: var(--_ink-2);
      white-space: nowrap;
    }
    .stage b { color: var(--_ink-1); font-weight: 600; }
    .arrow { color: var(--_muted); }
    table { border-collapse: collapse; width: 100%; margin-top: 10px; font-size: 13.5px; }
    th { text-align: left; color: var(--_muted); font-weight: 500; border-bottom: 1px solid var(--_baseline); padding: 6px 10px 6px 0; }
    td { border-bottom: 1px solid var(--_grid); padding: 7px 10px 7px 0; vertical-align: top; }
    td.num { font-variant-numeric: tabular-nums; white-space: nowrap; }
    .chip {
      display: inline-block;
      border: 1px solid var(--_baseline);
      border-radius: 4px;
      padding: 0 6px;
      font-size: 12px;
      color: var(--_ink-2);
      white-space: nowrap;
    }
    .minitrack { background: var(--_grid); border-radius: 3px; height: 8px; width: 90px; display: inline-block; vertical-align: middle; }
    .minifill { background: var(--_series); border-radius: 0 3px 3px 0; height: 8px; min-width: 2px; }
    .warn { color: var(--_ink-2); }
    .warn .icon { color: var(--_ink-1); }
    tr.current { font-weight: 650; }
    .marker { color: var(--_series); font-weight: 650; }
    .level-cap { margin: 12px 0 0; }
    ul { margin: 8px 0; padding-left: 22px; }
    li { margin: 4px 0; }
    a { color: var(--_series); }
    footer { margin-top: 48px; color: var(--_muted); font-size: 13px; border-top: 1px solid var(--_grid); padding-top: 14px; }
    .error { color: var(--_bad); padding: 16px; }
  `;

  willUpdate(changed) {
    if (changed.has('srcFramework') && this.srcFramework) {
      this.#fetchInto('framework', this.srcFramework);
    }
    if (changed.has('srcAssessment') && this.srcAssessment) {
      this.#fetchInto('assessment', this.srcAssessment);
    }
    if (changed.has('srcPrev') && this.srcPrev) {
      this.#fetchInto('prev', this.srcPrev);
    }
  }

  async #fetchInto(prop, url) {
    try {
      const resp = await fetch(url);
      if (!resp.ok) throw new Error(`${url}: HTTP ${resp.status}`);
      this[prop] = await resp.json();
    } catch (err) {
      this._error = String(err);
    }
  }

  render() {
    if (this._error) return html`<div class="error">scale-report: ${this._error}</div>`;
    if (!this.framework || !this.assessment) {
      return html`<div class="meta" style="padding:16px">Loading SCALE report…</div>`;
    }
    let rollup;
    let prevRollup = null;
    try {
      rollup = computeRollup(this.framework, this.assessment);
      if (this.prev) prevRollup = computeRollup(this.framework, this.prev);
    } catch (err) {
      return html`<div class="error">scale-report: ${String(err)}</div>`;
    }
    const f = this.framework;
    const movers = prevRollup ? buildMovers(f, rollup, prevRollup) : [];
    return html`
      <main>
        <h1>${f.name} — Platformization Report</h1>
        <p class="sub">${f.description}</p>
        <p class="meta">
          Period ${this.assessment.period}${this.prev ? html` · compared with ${this.prev.period}` : nothing}
        </p>
        ${this.#renderTiles(rollup, prevRollup)}
        ${(f.narratives ?? []).filter((n) => n.kind === 'thesis').map((n) => this.#renderNarrative(n, true))}
        ${movers.length ? this.#renderMovers(movers) : nothing}
        ${(f.domains ?? []).map((d) => this.#renderDomain(d, rollup, prevRollup))}
        ${this.#renderCoverage(rollup)}
        <footer>
          Aspect scores are the mean attainment (value vs. target, clamped to 0–1) of targeted, owned
          metrics; overall scores average domain scores so metric-heavy domains don't dominate. Rollup
          and maturity semantics are defined by the SCALE specification. External model levels belong to
          their publishers; PRISM crosswalks are explicit SCALE mappings.
        </footer>
      </main>
    `;
  }

  #renderTiles(rollup, prevRollup) {
    const tiles = ['standards', 'consumption', 'automation', 'leverage', 'effectiveness'].map((aspect) => {
      const score = rollup.aspects.find((a) => a.aspect === aspect) ?? null;
      const prev = prevRollup?.aspects.find((a) => a.aspect === aspect) ?? null;
      return { aspect, score, delta: score && prev ? (score.score - prev.score) * 100 : null };
    });
    return html`<div class="tiles">
      ${tiles.map(
        (t) => html`<div class="tile" title="${aspectDisplayName(t.aspect)} aspect score">
          <div class="label"><b>${aspectLetter(t.aspect)}</b> · ${aspectDisplayName(t.aspect)}</div>
          ${t.score
            ? html`<div class="value">${pct(t.score.score * 100)}<small>/100</small></div>
                ${t.delta === null
                  ? nothing
                  : t.delta > 0.5
                    ? html`<div class="delta up">▲ ${pts(t.delta)} pts</div>`
                    : t.delta < -0.5
                      ? html`<div class="delta down">▼ ${pts(t.delta)} pts</div>`
                      : html`<div class="delta flat">— flat</div>`}`
            : html`<div class="value none">–</div>
                <div class="meta">no eligible metrics</div>`}
        </div>`,
      )}
    </div>`;
  }

  #renderNarrative(n, card = false) {
    const body = html`
      ${n.title ? html`<h4>${n.title}${n.kind === 'outlook' ? ' (outlook)' : ''}</h4>` : nothing}
      <p>${n.body}</p>
      ${n.initiatives?.length
        ? html`<div class="meta">initiatives: ${n.initiatives.join(', ')}</div>`
        : nothing}
      ${n.owner
        ? html`<div class="meta">owner: ${n.owner}${n.reviewBy ? ` · review by ${n.reviewBy}` : ''}</div>`
        : nothing}
    `;
    return card
      ? html`<section class="card narrative">${body}</section>`
      : html`<div class="narrative">${body}</div>`;
  }

  #renderMovers(movers) {
    return html`<h2>What moved</h2>
      <section class="card">
        <ul>
          ${movers.map(
            (m) => html`<li>
              <b>${m.domain}</b> ${m.aspect}
              <span class="delta ${m.delta > 0 ? 'up' : 'down'}">
                ${m.delta > 0 ? '▲' : '▼'} ${pct(m.prev)} → ${pct(m.curr)} (${pts(m.delta)} pts)
              </span>
              ${m.contributions.length
                ? html`<span class="meta"> — driven by ${m.contributions.join(', ')}</span>`
                : nothing}
            </li>`,
          )}
        </ul>
      </section>`;
  }

  #renderDomain(d, rollup, prevRollup) {
    const journey = (this.assessment.narratives ?? []).filter(
      (n) => n.scope?.type === 'domain' && n.scope?.ref === d.id,
    );
    const dr = domainRollupFor(rollup, d.id);
    const prevDR = prevRollup ? domainRollupFor(prevRollup, d.id) : null;
    const lenses = (this.framework.externalModels ?? []).filter((em) => em.domain === d.id);
    let domainChip = '';
    if ((d.capabilities ?? []).some((c) => hasLadderedMetrics(c))) {
      const rung = domainMaturity(d, this.assessment);
      domainChip = rung ? `L${rung.level} · ${rung.name}` : 'N/A';
    }
    return html`
      <h2>
        ${d.name}${domainChip
          ? html` <span class="chip" title="Domain maturity: lowest rung across laddered capabilities (weakest link)">${domainChip}</span>`
          : nothing}${d.status === 'draft' ? html` <span class="chip">draft</span>` : nothing}
      </h2>
      ${(d.narratives ?? []).filter((n) => n.kind === 'thesis').map((n) => this.#renderNarrative(n))}
      ${(d.dimensions ?? []).map(
        (dim) => html`<div class="stages" title="${dim.name}: the domain's journey, left to right">
          ${dim.stages.map(
            (s, i) => html`${i ? html`<span class="arrow">→</span>` : nothing}<span class="stage"
                ><b>${s.name}</b></span
              >`,
          )}
        </div>`,
      )}
      ${dr ? this.#renderBars(dr, prevDR) : nothing}
      ${journey.map((n) => this.#renderNarrative(n, true))}
      ${lenses.length
        ? html`<h3>External framework lens</h3>${lenses.map((em) => this.#renderLens(em, d))}`
        : nothing}
      ${this.#renderCapabilities(d)}
    `;
  }

  #renderBars(dr, prevDR) {
    return html`<section class="card">
      ${dr.aspects.map((as) => {
        const prev = prevDR?.aspects.find((p) => p.aspect === as.aspect) ?? null;
        const delta = prev ? (as.score - prev.score) * 100 : null;
        return html`<div class="barrow" title="${aspectDisplayName(as.aspect)}: ${pct(as.score * 100)}/100">
          <span class="name"><b>${aspectLetter(as.aspect)}</b> ${aspectDisplayName(as.aspect)}</span>
          <span class="track"><span class="fill" style="width: ${pct(as.score * 100)}%"></span></span>
          <span class="barval"
            >${pct(as.score * 100)}${delta === null
              ? nothing
              : delta > 0.5
                ? html` <span class="delta up">▲${pts(delta)}</span>`
                : delta < -0.5
                  ? html` <span class="delta down">▼${pts(delta)}</span>`
                  : nothing}</span
          >
        </div>`;
      })}
    </section>`;
  }

  #renderLens(em, d) {
    const current = new Map();
    for (const c of d.capabilities ?? []) {
      for (const fm of c.frameworks ?? []) {
        if (fm.framework === em.id && fm.reference) {
          if (!current.has(fm.reference)) current.set(fm.reference, []);
          current.get(fm.reference).push(c.name);
        }
      }
    }
    return html`<section class="card">
      <p class="sub">
        <b>${em.name}</b> (${em.publisher})${em.sourceUrl
          ? html` · <a href="${em.sourceUrl}">source</a>`
          : nothing}${em.interpretation ? html` · <span class="chip">${em.interpretation}</span>` : nothing}
      </p>
      <table>
        <tr><th>Level</th><th>PRISM</th><th>Where we are</th></tr>
        ${em.levels.map((l) => {
          const caps = current.get(l.id) ?? [];
          return html`<tr class="${caps.length ? 'current' : ''}">
            <td>
              ${l.ordinal}. ${l.name}
              ${l.description ? html`<div class="meta">${l.description}</div>` : nothing}
            </td>
            <td class="num">${l.prismLevel ? `M${l.prismLevel}` : ''}</td>
            <td>${caps.length ? html`<span class="marker">◄ current practice</span> — ${caps.join(', ')}` : nothing}</td>
          </tr>`;
        })}
      </table>
    </section>`;
  }

  #renderCapabilities(d) {
    const groups = (d.capabilities ?? []).filter((c) => c.metrics?.length);
    if (!groups.length) return nothing;
    return html`<section class="card">
      ${groups.map((c) => {
        const hasLadder = hasLadderedMetrics(c);
        const rung = hasLadder ? capabilityMaturity(c, this.assessment) : null;
        const chip = hasLadder ? (rung ? `L${rung.level} · ${rung.name}` : 'N/A') : '';
        return html`<div class="level-cap">
          <b>${c.name}</b>
          ${chip
            ? html` <span class="chip" title="Capability maturity: lowest rung across laddered metrics (weakest link)">${chip}</span>`
            : nothing}
          <table>
            <tr><th>Metric</th><th>Aspect</th><th>Value</th><th>Target</th><th>Attainment</th><th>Owner</th></tr>
            ${c.metrics.map((m) => this.#renderMetricRow(m))}
          </table>
        </div>`;
      })}
    </section>`;
  }

  #renderMetricRow(m) {
    const obs = observation(this.assessment, m.id);
    let maturity = '';
    if (m.maturity) {
      const rung = metricMaturity(m, obs);
      maturity = rung ? `L${rung.level} · ${rung.name}` : obs ? 'below ladder' : 'N/A · not tracked';
    }
    let value = '—';
    let note = '';
    let attain = null;
    if (!obs) {
      note = rollupEligible(m) ? 'not measured' : 'tracked only (no target/owner)';
    } else {
      value = formatValue(m, obs.value, obs.numerator, obs.denominator);
      if (rollupEligible(m)) attain = attainment(m, obs.value) * 100;
      else note = 'tracked only (no target/owner)';
    }
    return html`<tr>
      <td>
        ${m.name}
        ${maturity ? html` <span class="chip" title="Maturity per this metric's ladder">${maturity}</span>` : nothing}
        ${note ? html`<div class="meta">${note}</div>` : nothing}
      </td>
      <td>
        <span class="chip" title="${aspectDisplayName(m.aspect)}${m.consumptionKind ? ` (${m.consumptionKind})` : ''}">
          ${aspectLetter(m.aspect)}${m.consumptionKind ? ` · ${m.consumptionKind}` : ''}
        </span>
      </td>
      <td class="num">${value}</td>
      <td class="num">${m.target ? formatValue(m, m.target.value, null, null) : ''}</td>
      <td>
        ${attain === null
          ? nothing
          : html`<span class="minitrack"><span class="minifill" style="width: ${pct(attain)}%"></span></span>
              <span class="barval">${pct(attain)}%</span>`}
      </td>
      <td class="meta">${m.owner ?? ''}</td>
    </tr>`;
  }

  #renderCoverage(rollup) {
    if (!rollup.missing.length && !rollup.excluded.length) return nothing;
    return html`<h2>Coverage</h2>
      <section class="card">
        ${rollup.missing.length
          ? html`<p class="warn">
              <span class="icon">⚠</span> <b>Not measured this period</b> (rollup-eligible, no
              observation): ${rollup.missing.join(', ')}
            </p>`
          : nothing}
        ${rollup.excluded.length
          ? html`<p class="warn">
              <span class="icon">⚠</span> <b>Tracked but excluded from rollups</b> (missing target or
              owner): ${rollup.excluded.join(', ')}
            </p>`
          : nothing}
      </section>`;
  }
}

customElements.define('scale-report', ScaleReport);
