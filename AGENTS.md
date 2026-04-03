# AGENTS.md

Guidelines for agentic coding assistants working in this repository.

## What this repo is

A personal landing page for Tom Vendolsky, built with
[Hugo](https://gohugo.io/) and the
[Blowfish theme](https://github.com/nunocoracao/blowfish) (git submodule).
The only deployable artifact is the static site Hugo emits into `public/`.
There is no application code, test suite, or package manager.

## Repository structure

```text
content/            → Page content (Markdown + TOML front matter)
  _index.md         → Homepage — the only content file
layouts/            → Custom overrides on top of the Blowfish theme
  shortcodes/       → accordionItem.html, timelineItem.html
  partials/         → favicons.html, home/ (empty — reserved)
assets/img/         → Images referenced from content or CSS
static/img/         → Favicons and web-manifest (served verbatim)
themes/blowfish/    → Theme submodule — do NOT edit files here
hugo.toml           → Site configuration
.pre-commit-config.yaml → Pre-commit hook that runs hugo --minify
.github/workflows/  → CI/CD pipelines (test + deploy)
```

## Build and preview

Hugo version in use: **0.159.1 extended**.

```bash
# Live-reload development server
hugo server

# Production build (same as CI)
hugo --minify
```

Output lands in `public/` — this directory is gitignored.

There is no test framework. "Tests" are a successful `hugo --minify`
build: a broken template, bad shortcode call, or missing asset causes
Hugo to exit non-zero and the CI job fails.

To verify a change before committing:

```bash
hugo --minify          # must exit 0 with no ERRORs in output
hugo server            # visual inspection in the browser
```

The pre-commit hook runs `hugo --minify` automatically on every commit.
Install it once with:

```bash
pre-commit install
```

## Content conventions

### Front matter

Homepage (`content/_index.md`) uses **TOML** delimiters (`+++`):

```toml
+++
title = "Hi, I'm Tom"
+++
```

Do not switch to YAML (`---`) front matter — stay consistent with the
existing file.

### Shortcodes

Custom shortcodes live in `layouts/shortcodes/`. Always use named
parameters (never positional):

```md
{{< timelineItem
  icon="worktree"
  header="Job Title"
  company="Company Name"
  badge="Month YYYY - Month YYYY"
  subheader="Location (Mode)"
  md="true"
>}}
Markdown body here.
{{< /timelineItem >}}
```

```md
{{< accordionItem
  title="Section Title"
  icon="icon-name"
  items="Item A|Item B|Item C"
>}}
{{< /accordionItem >}}
```

Use `|` as the delimiter for pipe-separated `items` lists.
Use `md="true"` only when the inner content contains Markdown.

### Dates in badges

Format: `Mon YYYY - Mon YYYY` or `Mon YYYY - Present`.
Example: `Feb 2025 - Present`.

## Hugo template conventions

Templates live in `layouts/`. They override or extend Blowfish — never
modify `themes/blowfish/` directly.

- Use `{{ with }}` to guard optional values instead of `{{ if }}` +
  separate `{{ . }}`.
- Whitespace trimming: use `{{- -}}` only around block-level content
  inside `{{ if $md }}` branches to avoid extra `<p>` tags.
- CSS: Tailwind utility classes only. Do not introduce custom CSS files.
- Use `relURL` (not `absURL`) for all internal asset references.

## Commit style

Follow Conventional Commits. Observed prefixes in this repo:

| Prefix | Use for |
| ------ | ------- |
| `feat` | New content, features, or dependency bumps |
| `fix`  | Corrections to existing content or templates |
| `chore` | Tooling, CI, dependency updates (non-content) |

Scope is optional but encouraged for targeted changes:
`feat(content): add new timeline entry`.

## CI/CD

Two workflows run on GitHub Actions (Hugo 0.159.1 extended,
`ubuntu-latest`):

- **test.yml** — runs on every push and PR; executes `hugo --minify`
- **deploy.yaml** — runs on push to `main`; builds then deploys to
  GitHub Pages

Both check out with `submodules: recursive` to pull in the Blowfish
theme. A build that exits non-zero blocks merge.

## What NOT to do

- **Do not edit** anything inside `themes/blowfish/` — use layout
  overrides in `layouts/` instead.
- **Do not commit** `public/` or `resources/` — both are gitignored
  generated artifacts.
- **Do not add** a package manager (npm, go.mod, etc.) unless the site
  genuinely requires it; this is intentionally dependency-free.
- **Do not introduce** custom CSS files; use Tailwind classes available
  from Blowfish.

## See also

- [Hugo documentation](https://gohugo.io/documentation/)
- [Blowfish theme docs](https://blowfish.page/docs/)
- [Conventional Commits](https://www.conventionalcommits.org/)
