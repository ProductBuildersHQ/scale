# `<scale-report>` Web Component

A [Lit](https://lit.dev) web component that renders the SCALE platform story report from the framework JSON IR and assessment JSON — the same output as `scale report`, computed client-side, injectable into any div on a React, MkDocs, or plain HTML site.

- `src/compute.js` — dependency-free mirror of the Go reference implementation's rollup and maturity-ladder semantics. Pinned to Go-verified values by `npm test`. If this file and the Go implementation disagree, Go wins.
- `src/scale-report.js` — the `<scale-report>` element (Shadow DOM, light/dark via `prefers-color-scheme`, themeable via `--scale-*` custom properties).

## Producing the JSON IR

A browser can't assemble `catalog/domains/*.json`, so export the catalog to a single file:

```bash
scale export -catalog catalog -o framework.json
```

Assessments are already single JSON files and are consumed as-is.

## Plain HTML / MkDocs

```html
<script type="importmap">
  { "imports": { "lit": "https://esm.run/lit@3" } }
</script>
<script type="module" src="/js/scale-report.js"></script>

<div id="report-container">
  <scale-report
    src-framework="/data/framework.json"
    src-assessment="/data/assessments/2026-q3.json"
    src-prev="/data/assessments/2026-q2.json"></scale-report>
</div>
```

For MkDocs, put the snippet in a markdown file (MkDocs passes raw HTML through), copy `src/*.js` into `docs/js/`, and add the `<script>` tags via `extra_javascript` or the snippet itself. Note the import map must appear before the module script.

## React

```bash
npm install lit
```

```jsx
import { useEffect, useRef } from 'react';
import '@productbuildershq/scale-report'; // or a relative path to src/scale-report.js

export function ScaleReportView({ framework, assessment, prev }) {
  const ref = useRef(null);
  useEffect(() => {
    if (!ref.current) return;
    ref.current.framework = framework;
    ref.current.assessment = assessment;
    ref.current.prev = prev ?? null;
  }, [framework, assessment, prev]);
  return <scale-report ref={ref} />;
}
```

Setting the object properties directly (as above) skips fetching; the `src-*` attributes are the no-build alternative.

## Theming

Colors are exposed as `--scale-*` custom properties and follow `prefers-color-scheme` by default. Override from page CSS:

```css
scale-report {
  --scale-series: #7c3aed;   /* bars, links, markers */
  --scale-page: transparent; /* blend into the host page */
}
```

Available slots: `--scale-page`, `--scale-surface`, `--scale-ink-1`, `--scale-ink-2`, `--scale-muted`, `--scale-grid`, `--scale-baseline`, `--scale-series`, `--scale-good`, `--scale-bad`, `--scale-border`.

## Demo

```bash
scale export -catalog catalog -o examples/framework-ir.json   # once, or after catalog changes
python3 -m http.server 8000                                    # from the repo root
# open http://localhost:8000/web/demo/
```

## Tests

```bash
cd web && npm test
```
