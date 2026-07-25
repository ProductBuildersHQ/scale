package report

const pageCSS = `
:root {
  color-scheme: light;
  --page:      #f9f9f7;
  --surface:   #fcfcfb;
  --ink-1:     #0b0b0b;
  --ink-2:     #52514e;
  --muted:     #898781;
  --grid:      #e1e0d9;
  --baseline:  #c3c2b7;
  --series:    #2a78d6;
  --good:      #006300;
  --bad:       #d03b3b;
  --border:    rgba(11,11,11,0.10);
}
@media (prefers-color-scheme: dark) {
  :root:where(:not([data-theme="light"])) {
    color-scheme: dark;
    --page:      #0d0d0d;
    --surface:   #1a1a19;
    --ink-1:     #ffffff;
    --ink-2:     #c3c2b7;
    --muted:     #898781;
    --grid:      #2c2c2a;
    --baseline:  #383835;
    --series:    #3987e5;
    --good:      #0ca30c;
    --bad:       #d03b3b;
    --border:    rgba(255,255,255,0.10);
  }
}
:root[data-theme="dark"] {
  color-scheme: dark;
  --page:      #0d0d0d;
  --surface:   #1a1a19;
  --ink-1:     #ffffff;
  --ink-2:     #c3c2b7;
  --muted:     #898781;
  --grid:      #2c2c2a;
  --baseline:  #383835;
  --series:    #3987e5;
  --good:      #0ca30c;
  --bad:       #d03b3b;
  --border:    rgba(255,255,255,0.10);
}
* { box-sizing: border-box; }
body {
  margin: 0;
  background: var(--page);
  color: var(--ink-1);
  font: 15px/1.55 system-ui, -apple-system, "Segoe UI", sans-serif;
}
main { max-width: 920px; margin: 0 auto; padding: 32px 20px 64px; }
h1 { font-size: 26px; margin: 0 0 2px; }
h2 { font-size: 19px; margin: 40px 0 12px; }
h3 { font-size: 16px; margin: 24px 0 8px; }
.sub { color: var(--ink-2); margin: 0 0 4px; }
.meta { color: var(--muted); font-size: 13px; }
section.card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 20px 22px;
  margin-top: 16px;
}
section.card h3 { margin: 0 0 6px; }
section.card.current { border-left: 3px solid var(--series); }
.narrative { border-left: 3px solid var(--grid); padding: 2px 0 2px 14px; margin: 12px 0; }
.narrative h4 { margin: 0 0 4px; font-size: 14px; }
.narrative p { margin: 0; color: var(--ink-2); }
.narrative .meta { margin-top: 6px; }
.tiles { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 12px; margin-top: 16px; }
.tile {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 14px 16px;
}
.tile .label { font-size: 13px; color: var(--ink-2); }
.tile .label b { color: var(--ink-1); }
.tile .value { font-size: 30px; font-weight: 650; margin-top: 2px; }
.tile .value small { font-size: 15px; font-weight: 400; color: var(--muted); }
.delta { font-size: 13px; font-weight: 600; }
.delta.up   { color: var(--good); }
.delta.down { color: var(--bad); }
.delta.flat { color: var(--muted); font-weight: 400; }
.none { color: var(--muted); }
.barrow { display: grid; grid-template-columns: 130px 1fr 90px; align-items: center; gap: 10px; margin: 7px 0; }
.barrow .name { font-size: 13px; color: var(--ink-2); }
.track { background: var(--grid); border-radius: 4px; height: 10px; position: relative; }
.fill  { background: var(--series); border-radius: 0 4px 4px 0; height: 10px; min-width: 2px; }
.barval { font-size: 13px; font-variant-numeric: tabular-nums; }
.stages { display: flex; flex-wrap: wrap; align-items: center; gap: 6px; margin: 10px 0 4px; }
.stage {
  border: 1px solid var(--baseline);
  border-radius: 999px;
  padding: 3px 12px;
  font-size: 13px;
  color: var(--ink-2);
  white-space: nowrap;
}
.stage b { color: var(--ink-1); font-weight: 600; }
.arrow { color: var(--muted); }
table { border-collapse: collapse; width: 100%; margin-top: 10px; font-size: 13.5px; }
th { text-align: left; color: var(--muted); font-weight: 500; border-bottom: 1px solid var(--baseline); padding: 6px 10px 6px 0; }
td { border-bottom: 1px solid var(--grid); padding: 7px 10px 7px 0; vertical-align: top; }
td.num { font-variant-numeric: tabular-nums; white-space: nowrap; }
.chip {
  display: inline-block;
  border: 1px solid var(--baseline);
  border-radius: 4px;
  padding: 0 6px;
  font-size: 12px;
  color: var(--ink-2);
  white-space: nowrap;
}
.minitrack { background: var(--grid); border-radius: 3px; height: 8px; width: 90px; display: inline-block; vertical-align: middle; }
.minifill { background: var(--series); border-radius: 0 3px 3px 0; height: 8px; min-width: 2px; }
.warn { color: var(--ink-2); }
.warn .icon { color: var(--ink-1); }
.current { font-weight: 650; }
.marker { color: var(--series); font-weight: 650; }
.level-cap { margin: 12px 0 0; }
ul { margin: 8px 0; padding-left: 22px; }
li { margin: 4px 0; }
a { color: var(--series); }
footer { margin-top: 48px; color: var(--muted); font-size: 13px; border-top: 1px solid var(--grid); padding-top: 14px; }
`

const reportTemplate = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Framework.Name}} — {{.Period}}</title>
<style>` + pageCSS + `</style>
</head>
<body>
<main>
<h1>{{.Framework.Name}} — Platform Report</h1>
<p class="sub">{{.Framework.Description}}</p>
<p class="meta">Period {{.Period}}{{if .PrevPeriod}} · compared with {{.PrevPeriod}}{{end}}{{if .GeneratedAt}} · generated {{.GeneratedAt}}{{end}}</p>

<div class="tiles">
{{range .Tiles}}
  <div class="tile" title="{{.Name}} aspect score: mean attainment of targeted, owned metrics">
    <div class="label"><b>{{.Letter}}</b> · {{.Name}}</div>
    {{if .HasScore}}
      <div class="value">{{pct .Score}}<small>/100</small></div>
      {{if .HasDelta}}
        {{if gt .Delta 0.5}}<div class="delta up">▲ {{pts .Delta}} pts</div>
        {{else if lt .Delta -0.5}}<div class="delta down">▼ {{pts .Delta}} pts</div>
        {{else}}<div class="delta flat">— flat</div>{{end}}
      {{end}}
    {{else}}
      <div class="value none">–</div>
      <div class="meta">no eligible metrics</div>
    {{end}}
  </div>
{{end}}
</div>

{{range .Narratives}}
<section class="card narrative">
  {{if .Title}}<h4>{{.Title}}</h4>{{end}}
  <p>{{.Body}}</p>
  {{if .Owner}}<div class="meta">owner: {{.Owner}}{{if .ReviewBy}} · review by {{.ReviewBy}}{{end}}</div>{{end}}
</section>
{{end}}

{{if .Movers}}
<h2>What moved</h2>
<section class="card">
  <ul>
  {{range .Movers}}
    <li>
      <b>{{.Domain}}</b> {{.Aspect}}
      {{if gt .Delta 0.0}}<span class="delta up">▲ {{pct .Prev}} → {{pct .Curr}} ({{pts .Delta}} pts)</span>
      {{else}}<span class="delta down">▼ {{pct .Prev}} → {{pct .Curr}} ({{pts .Delta}} pts)</span>{{end}}
      {{if .Contributions}}<span class="meta"> — driven by {{range $i, $c := .Contributions}}{{if $i}}, {{end}}{{$c}}{{end}}</span>{{end}}
    </li>
  {{end}}
  </ul>
</section>
{{end}}

{{range .Domains}}
<h2>{{.Domain.Name}}{{if .Maturity}} <span class="chip" title="Domain maturity: the lowest rung across its laddered capabilities (weakest link)">{{.Maturity}}</span>{{end}}{{if eq .Domain.Status "draft"}} <span class="chip">draft</span>{{end}}</h2>

{{range .Thesis}}
<div class="narrative">
  {{if .Title}}<h4>{{.Title}}</h4>{{end}}
  <p>{{.Body}}</p>
</div>
{{end}}

{{range .Dimensions}}
<div class="stages" title="{{.Name}}: the domain's journey, read left to right">
  {{range $i, $s := .Stages}}{{if $i}}<span class="arrow">→</span>{{end}}<span class="stage"><b>{{$s.Name}}</b></span>{{end}}
</div>
{{end}}

{{if .Bars}}
<section class="card">
  {{range .Bars}}
  <div class="barrow" title="{{.Name}}: {{pct .Score}}/100">
    <span class="name"><b>{{.Letter}}</b> {{.Name}}</span>
    <span class="track"><span class="fill" style="width: {{pct .Score}}%"></span></span>
    <span class="barval">{{pct .Score}}
      {{- if .HasDelta}}
        {{- if gt .Delta 0.5}} <span class="delta up">▲{{pts .Delta}}</span>
        {{- else if lt .Delta -0.5}} <span class="delta down">▼{{pts .Delta}}</span>
        {{- end}}
      {{- end}}</span>
  </div>
  {{end}}
</section>
{{end}}

{{range .Journey}}
<section class="card narrative">
  {{if .Title}}<h4>{{.Title}}{{if eq .Kind "outlook"}} (outlook){{end}}</h4>{{end}}
  <p>{{.Body}}</p>
  {{if .Initiatives}}<div class="meta">initiatives: {{range $i, $init := .Initiatives}}{{if $i}}, {{end}}{{$init}}{{end}}</div>{{end}}
</section>
{{end}}

{{if .Lenses}}
<h3>External framework lens</h3>
{{range .Lenses}}
<section class="card">
  <p class="sub"><b>{{.Model.Name}}</b> ({{.Model.Publisher}}){{if .Model.SourceURL}} · <a href="{{.Model.SourceURL}}">source</a>{{end}}{{if .Model.Interpretation}} · <span class="chip">{{.Model.Interpretation}}</span>{{end}}</p>
  <table>
    <tr><th>Level</th><th>PRISM</th><th>Where we are</th></tr>
    {{range .Rows}}
    <tr{{if .Current}} class="current"{{end}}>
      <td>{{.Level.Ordinal}}. {{.Level.Name}}{{if .Level.Description}}<div class="meta">{{.Level.Description}}</div>{{end}}</td>
      <td class="num">{{.PRISM}}</td>
      <td>{{if .Current}}<span class="marker">◄ current practice</span> — {{range $i, $c := .Capabilities}}{{if $i}}, {{end}}{{$c}}{{end}}{{end}}</td>
    </tr>
    {{end}}
  </table>
</section>
{{end}}
{{end}}

{{if .Groups}}
<section class="card">
  {{range .Groups}}
  <div class="level-cap">
    <b>{{.Name}}</b>{{if .Maturity}} <span class="chip" title="Capability maturity: the lowest rung across its laddered metrics (weakest link)">{{.Maturity}}</span>{{end}}
    <table>
      <tr><th>Metric</th><th>Aspect</th><th>Value</th><th>Target</th><th>Attainment</th><th>Owner</th></tr>
      {{range .Rows}}
      <tr>
        <td>{{.Name}}{{if .Maturity}} <span class="chip" title="Maturity per this metric's ladder">{{.Maturity}}</span>{{end}}{{if .Note}}<div class="meta">{{.Note}}</div>{{end}}</td>
        <td><span class="chip" title="{{.AspectName}}{{if .Kind}} ({{.Kind}}){{end}}">{{.AspectLetter}}{{if .Kind}} · {{kindLabel .Kind}}{{end}}</span></td>
        <td class="num">{{.Value}}</td>
        <td class="num">{{.Target}}</td>
        <td>{{if .HasAttain}}<span class="minitrack"><span class="minifill" style="width: {{pct .Attain}}%"></span></span> <span class="barval">{{pct .Attain}}%</span>{{end}}</td>
        <td class="meta">{{.Owner}}</td>
      </tr>
      {{end}}
    </table>
  </div>
  {{end}}
</section>
{{end}}
{{end}}

{{if or .Missing .Excluded}}
<h2>Coverage</h2>
<section class="card">
  {{if .Missing}}<p class="warn"><span class="icon">⚠</span> <b>Not measured this period</b> (rollup-eligible, no observation): {{range $i, $m := .Missing}}{{if $i}}, {{end}}{{$m}}{{end}}</p>{{end}}
  {{if .Excluded}}<p class="warn"><span class="icon">⚠</span> <b>Tracked but excluded from rollups</b> (missing target or owner): {{range $i, $m := .Excluded}}{{if $i}}, {{end}}{{$m}}{{end}}</p>{{end}}
</section>
{{end}}

<footer>
Aspect scores are the mean attainment (value vs. target, clamped to 0–1) of targeted, owned metrics; overall scores average domain scores so metric-heavy domains don't dominate. Rollup semantics are defined by the SCALE specification. External model levels belong to their publishers; PRISM crosswalks are explicit SCALE mappings.
</footer>
</main>
</body>
</html>
`

const modelTemplate = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Model.Name}} — {{.Period}}</title>
<style>` + pageCSS + `</style>
</head>
<body>
<main>
<h1>{{.Model.Name}}</h1>
<p class="sub">{{.Model.Description}}</p>
<p class="meta">Publisher {{.Model.Publisher}} · period {{.Period}}{{if .PrevPeriod}} · compared with {{.PrevPeriod}}{{end}}{{if .GeneratedAt}} · generated {{.GeneratedAt}}{{end}}{{if .Model.SourceURL}} · <a href="{{.Model.SourceURL}}">source</a>{{end}}{{if .Model.Interpretation}} · <span class="chip">{{.Model.Interpretation}}</span>{{end}}{{if .Model.RetrievedAt}} · retrieved {{.Model.RetrievedAt}}{{end}}</p>

<section class="card">
  <p><b>{{.Position}}</b></p>
  {{if .Next}}<p class="sub">Next: <b>{{.Next.Name}}</b>{{if .Next.Description}} — {{.Next.Description}}{{end}}</p>{{end}}
</section>

<h2>The ladder</h2>
{{range .Levels}}
<section class="card{{if .Current}} current{{end}}">
  <h3>{{.Level.Ordinal}}. {{.Level.Name}}
    {{- if .PRISM}} <span class="chip" title="Explicit SCALE crosswalk to PRISM maturity">PRISM {{.PRISM}}</span>{{end}}
    {{- if .Current}} <span class="marker">► current practice</span>{{end}}
    {{- if .IsNext}} <span class="chip">next up</span>{{end}}</h3>
  {{if .Level.Description}}<p class="sub">{{.Level.Description}}</p>{{end}}
  {{range .Capabilities}}
  <div class="level-cap">
    <b>{{.Name}}</b>{{if .Maturity}} <span class="chip" title="Capability maturity: the lowest rung across its laddered metrics (weakest link)">{{.Maturity}}</span>{{end}}{{if .Why}} — <span class="meta">{{.Why}}</span>{{end}}
    {{if .Metrics}}
    <table>
      <tr><th>Metric</th><th>Aspect</th><th>Value</th><th>Target</th><th>Attainment</th></tr>
      {{range .Metrics}}
      <tr>
        <td>{{.Name}}{{if .Maturity}} <span class="chip" title="Maturity per this metric's ladder">{{.Maturity}}</span>{{end}}{{if .Note}}<div class="meta">{{.Note}}</div>{{end}}</td>
        <td><span class="chip" title="{{.AspectName}}{{if .Kind}} ({{.Kind}}){{end}}">{{.AspectLetter}}{{if .Kind}} · {{kindLabel .Kind}}{{end}}</span></td>
        <td class="num">{{.Value}}</td>
        <td class="num">{{.Target}}</td>
        <td>{{if .HasAttain}}<span class="minitrack"><span class="minifill" style="width: {{pct .Attain}}%"></span></span> <span class="barval">{{pct .Attain}}%</span>{{end}}</td>
      </tr>
      {{end}}
    </table>
    {{end}}
  </div>
  {{end}}
</section>
{{end}}

{{if .Bars}}
<h2>SCALE evidence — {{.DomainName}}</h2>
<section class="card">
  {{range .Bars}}
  <div class="barrow" title="{{.Name}}: {{pct .Score}}/100">
    <span class="name"><b>{{.Letter}}</b> {{.Name}}</span>
    <span class="track"><span class="fill" style="width: {{pct .Score}}%"></span></span>
    <span class="barval">{{pct .Score}}
      {{- if .HasDelta}}
        {{- if gt .Delta 0.5}} <span class="delta up">▲{{pts .Delta}}</span>
        {{- else if lt .Delta -0.5}} <span class="delta down">▼{{pts .Delta}}</span>
        {{- end}}
      {{- end}}</span>
  </div>
  {{end}}
</section>
{{end}}

{{range .Journey}}
<section class="card narrative">
  {{if .Title}}<h4>{{.Title}}{{if eq .Kind "outlook"}} (outlook){{end}}</h4>{{end}}
  <p>{{.Body}}</p>
  {{if .Initiatives}}<div class="meta">initiatives: {{range $i, $init := .Initiatives}}{{if $i}}, {{end}}{{$init}}{{end}}</div>{{end}}
</section>
{{end}}

{{if or .Missing .Excluded}}
<h2>Coverage</h2>
<section class="card">
  {{if .Missing}}<p class="warn"><span class="icon">⚠</span> <b>Not measured this period</b> (rollup-eligible, no observation): {{range $i, $m := .Missing}}{{if $i}}, {{end}}{{$m}}{{end}}</p>{{end}}
  {{if .Excluded}}<p class="warn"><span class="icon">⚠</span> <b>Tracked but excluded from rollups</b> (missing target or owner): {{range $i, $m := .Excluded}}{{if $i}}, {{end}}{{$m}}{{end}}</p>{{end}}</section>
{{end}}

<footer>
Levels belong to {{.Model.Publisher}}; the PRISM crosswalk and the current-practice placement are SCALE mappings authored in the catalog, not a vendor assessment. Evidence boundary: this report shows conformance of catalogued capabilities and metrics, not a complete organizational maturity assessment. Attainment is value vs. target, clamped to 0–1, per the SCALE specification.
</footer>
</main>
</body>
</html>
`
