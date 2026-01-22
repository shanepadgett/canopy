# Canopy Status

## Completed

- Core data types (`Site`, `Page`, `Section`, `Config`), front matter parsing/validation/defaults.
- Content discovery, slug + URL computation, section + tag indexing.
- Markdown rendering (headings, paragraphs, links, lists, code blocks, inline code, emphasis, blockquotes, horizontal rules).
- Template engine with base/page/list/home layouts and default templates.
- Build pipeline (config -> content -> markdown -> templates -> output).
- Static asset copying into `public/`.
- Tag pages (`/tags/<tag>/`) and tags index (`/tags/`).
- Machine-readable outputs: `rss.xml`, `sitemap.xml`, `robots.txt`.
- Sample site nav updated with Tags link.

## Next Up

- Search index (`search.json`) and minimal search UI.
