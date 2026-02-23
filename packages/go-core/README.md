# @ketab/go-core

Go types, constants, and validation for the Ketab Protocol.

## Package Structure

| File | Purpose |
|------|---------|
| `kind.go` | Event kind constants (38890–38893) |
| `types.go` | Content structs with `Validate()` methods |
| `validation/` | Event-level validation (tags, content, cross-field checks) |

## Usage

```go
import (
    core "github.com/joinnextblock/ketab-protocol/go-core"
    "github.com/joinnextblock/ketab-protocol/go-core/validation"
)

// Validate a Nostr event
result := validation.ValidateBookEvent(event)
if !result.Valid {
    log.Fatal(result.Message)
}
```

## Event Kinds

| Kind | Name | Description |
|------|------|-------------|
| 38890 | Library | Book curation container |
| 38891 | Book | Book metadata and chapter structure |
| 38892 | LibraryEntry | Library-specific book metadata |
| 38893 | Ketab | Individual content unit within a chapter |

## Imported By

- `city-relay` — validates incoming Ketab events
- `ketab` CLI — validates before publishing
