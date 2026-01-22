# Task 01 - Search Index And UI

## Goal

Generate `search.json` and ship a nav-integrated search UI that consumes it.

## Scope

- Build `search.json` with: title, url, summary, section, tags.
- Inject search icon button in default nav that opens an overlay modal.
- Add Cmd/Ctrl+K keyboard shortcut to open search.
- Add inline JS/CSS for fuzzy matching with scoring.
- Config flag `search.enabled` (default true) to disable.

## search.json Schema

```json
[
  {
    "url": "/blog/hello/",
    "title": "Hello World",
    "section": "blog",
    "tags": ["intro", "welcome"],
    "summary": "A brief introduction to the site..."
  }
]
```

## Config Extension

```go
type SearchConfig struct {
    Enabled bool `json:"enabled"` // default: true
}
```

## Files to Modify

| File | Change |
|------|--------|
| `internal/core/types.go` | Add `SearchConfig` struct and `Search` field to `Config` |
| `internal/config/loader.go` | Set `cfg.Search.Enabled = true` in defaults |
| `internal/build/build.go` | Add `renderSearchIndex()` when enabled |
| `internal/template/engine.go` | Modify `defaultBaseLayout` to inject search UI |

## Search UI Behavior

| Action | Trigger |
|--------|---------|
| Open modal | Click search icon OR Cmd/Ctrl+K |
| Close modal | Esc OR click backdrop |
| Navigate results | Arrow keys |
| Select result | Enter |
| Filter | Type in input (debounced) |

## Fuzzy Scoring Algorithm

Score each entry by:

1. Title match bonus: +100 for matches in title
2. Consecutive chars: +10 per consecutive match
3. Word boundary bonus: +5 when match follows space/punctuation
4. Position penalty: -1 per character from start

Results sorted by descending score, capped at 10 displayed.
