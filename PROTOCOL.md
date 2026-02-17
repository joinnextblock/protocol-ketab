# Ketab Protocol Specification (LIBRARY-01)

> Composable interactive stories on Nostr.

## Event Kinds

| Kind | Name | Type | Description |
|------|------|------|-------------|
| 38890 | Library | Replaceable | Collection of books curated by a librarian |
| 38891 | Book | Replaceable | Book metadata + ordered chapter references |
| 38893 | Ketab | Replaceable | Atomic content unit — one thought, one card |
| 30023 | Chapter | Replaceable | NIP-23 long-form content (compiled from ketabs) |

### Protocol Kinds (38890, 38891, 38893)

- **Tags**: Single-letter only (`d`, `a`, `p`, `e`, `t`). For relay indexing. No multi-letter tags.
- **Content**: JSON string. All metadata lives here.

### Nostr-Native Kind (30023)

- Follows NIP-23 spec. `title`, `summary`, `published_at` tags are permitted.
- Content is markdown.

---

## Kind 38893 — Ketab

The atomic unit. One signed event. One passage. Individually addressable, citable, and engageable.

### Tags

| Tag | Required | Description |
|-----|----------|-------------|
| `d` | Yes | Unique identifier (UUID) |
| `a` | Yes | Parent chapter coordinate: `30023:<pubkey>:<chapter-d-tag>` |
| `t` | No | Topic tags for discovery |

### Content (JSON)

```json
{
  "title": "The Slave Markets of Seville",
  "index": 4,
  "body": "By 1495, Columbus had shipped over 500 enslaved...\n\n---\n\n**Sources**\n\n1. Las Casas, *Historia de las Indias*, Book I, Ch. 88"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Ketab title |
| `index` | number | 0-based position within chapter |
| `body` | string | Markdown. Optional footnotes after `---` separator |

### Rules

- `index` is 0-based. Display adds +1 for humans.
- Tag order never determines position. `index` is authoritative.
- Body may contain footnote superscripts (`[1]`, `[2]`) with corresponding sources after `---`.

---

## Kind 38891 — Book

### Tags

| Tag | Required | Description |
|-----|----------|-------------|
| `d` | Yes | Unique identifier (UUID) |
| `a` | Yes (repeated) | Chapter coordinates in reading order: `30023:<pubkey>:<chapter-d-tag>` |
| `p` | Yes | Author pubkey |
| `t` | No | Topic tags |

### Content (JSON)

```json
{
  "title": "The Copper Islands",
  "subtitle": "Cape Verde, Columbus, and the real slave trade.",
  "author": "NextBlock",
  "summary": "Short summary for previews.",
  "description": "Full description with details.",
  "image": "https://..."
}
```

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Book title |
| `subtitle` | string | Optional subtitle |
| `author` | string | Display name |
| `summary` | string | Short preview text |
| `description` | string | Full description |
| `image` | string | Cover image URL |

Chapter ordering follows `a` tag order in the book event.

---

## Kind 38890 — Library

### Tags

| Tag | Required | Description |
|-----|----------|-------------|
| `d` | Yes | Unique identifier (UUID) |
| `a` | Yes (repeated) | Book coordinates: `38891:<pubkey>:<book-d-tag>` |
| `t` | No | Topic tags |

### Content (JSON)

```json
{
  "name": "the library",
  "description": "Books published on Nostr. Read by citizens."
}
```

---

## Kind 30023 — Chapter (NIP-23)

Chapters are compiled views of ketabs. Follows NIP-23 spec.

### Tags

| Tag | Required | Description |
|-----|----------|-------------|
| `d` | Yes | Unique identifier (UUID) |
| `title` | Yes | Chapter title (NIP-23) |
| `a` | No (repeated) | Ketab coordinates: `38893:<pubkey>:<ketab-d-tag>` |
| `published_at` | No | Unix timestamp (NIP-23) |

### Content

Markdown. Must produce identical output to the concatenation of its ketabs' body fields (separated by `\n\n---\n\n`).

---

## Coordinates

All addressable events use the coordinate format:

```
<kind>:<pubkey>:<d-tag>
```

Examples:
- Library: `38890:<pubkey>:<library-uuid>`
- Book: `38891:<pubkey>:<book-uuid>`
- Chapter: `30023:<pubkey>:<chapter-uuid>`
- Ketab: `38893:<pubkey>:<ketab-uuid>`

Coordinates are used in `a` tags for cross-referencing and in engagement events for targeting.

---

## Engagement

Per-ketab engagement via standard Nostr kinds:

| Kind | NIP | Purpose |
|------|-----|---------|
| 1111 | NIP-22 | Threaded comments |
| 9802 | NIP-84 | Highlights |
| 9735 | NIP-57 | Zap receipts |
| 6 | NIP-18 | Reposts |

All reference the target's `a` coordinate.

---

## Discovery

### From a book `naddr`:

1. Decode `naddr` → fetch kind 38891 event
2. Read `a` tags → list of chapter coordinates
3. Fetch chapters (kind 30023) by `#d` filter
4. For each chapter, fetch ketabs (kind 38893) with `#a` filter matching `30023:<pubkey>:<chapter-d-tag>`
5. Sort ketabs by `index` in content JSON

### From a ketab `naddr`:

1. Decode `naddr` → fetch kind 38893 event
2. Read parent `a` tag → chapter coordinate
3. Fetch sibling ketabs with same parent `a` tag
4. Sort by `index` for prev/next navigation

---

## City Protocol Integration

Ketab events may reference City Protocol block events (kind 38808) for block-time timestamps. This is optional and does not affect the core protocol.

```
38808:<clock_pubkey>:org.cityprotocol:block:<height>:<hash>
```
