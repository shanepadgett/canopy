# Task 02 - Shortcodes

## Goal

Support template-backed shortcodes in Markdown with attributes and optional inner content.

## Scope

- Parse `{{< name attr="value" >}} ... {{< /name >}}`.
- Map shortcode name to `templates/shortcodes/<name>.html`.
- Support inline shortcodes without closing tag.
- Provide built-ins: `callout`, `figure`, `embed`/`youtube`, `toc`.

## Notes

- Shortcodes should render after Markdown parse or during Markdown render.
- TOC shortcode should access `Page.TOC`.
