# Ketab Protocol

Composable interactive stories on Nostr.

## Why

Nostr publishing works at the article level. You can publish a chapter. You can zap a chapter. But you can't address, cite, or engage with a single passage inside it.

Ketab Protocol adds the layer below. A **ketab** (kind 38893) is an atomic content unit — one thought, individually signed, addressable by `naddr`, and engageable on its own. Books are composed from ketabs. Ketabs can be recomposed across books. Every ketab carries its own engagement (zaps, comments, highlights) independent of the chapter it belongs to.

This makes content **composable at the passage level**. A librarian can curate a reading list that pulls ketab 4 from one book, ketab 12 from another, and ketab 1 from a third — each one a verified Nostr event from its original author. No copy-paste, no broken attribution, no aggregator in the middle.

## Four Layers

| Kind | Name | Description |
|------|------|-------------|
| **38890** | Library | Collection of books, curated by a librarian |
| **38891** | Book | Metadata + ordered chapter references |
| **38893** | Ketab | Atomic content unit — one card, one thought |
| 30023 | Chapter | NIP-23 long-form content (Nostr-native, compiles ketabs) |

Ketabs (38893) are the source of truth. Chapters (30023) are compiled views — the same content assembled for long-form readers. Both must produce identical markdown.

## Three Actors

- **Authors** publish books and ketabs. They create the content.
- **Librarians** curate books into libraries with personal metadata (notes, ratings, tags).
- **Readers** engage per-ketab: zaps, comments, highlights, reposts.

A single pubkey can play all three roles.

## Ketab Event (Kind 38893)

The ketab is the core innovation. Each ketab is:
- A standalone, individually addressable Nostr event
- One card in the reading experience
- Citable via `naddr` with its own URL, OG tags, and engagement
- Ordered by `index` field in content (0-based)

```json
{
  "kind": 38893,
  "tags": [
    ["d", "<uuid>"],
    ["a", "30023:<author-pubkey>:<chapter-uuid>", "<relay>"],
    ["t", "ketab"],
    ["t", "nonfiction"]
  ],
  "content": "{\"title\":\"The Slave Markets of Seville\",\"index\":4,\"body\":\"By 1495, Columbus had shipped...\\n\\n---\\n\\n**Sources**\\n\\n1. Las Casas, *Historia de las Indias*, Book I, Ch. 88\"}"
}
```

### Content JSON

```typescript
{
  title: string;    // Ketab title
  index: number;    // 0-based position within chapter
  body: string;     // Markdown body + optional footnotes after ---
}
```

### Hard Rules

- **Tags are single-letter only.** `d`, `a`, `p`, `e`, `t` — for relay indexing. No multi-letter tags (`title`, `summary`, `description`, `published_at`). Ever.
- **Content is JSON only.** All metadata lives in the content field as structured JSON. This is the inverse of how most NIPs work, and it's intentional.
- **Tag order NEVER determines position.** Ordering comes from `index` field in content JSON.
- **`index` is 0-based.** Ketabs are arrays of events. Display adds +1 for humans.
- **Each ketab references its parent chapter** via `a` tag: `['a', '30023:<pubkey>:<chapter-uuid>', '<relay>']`

## Book Event (Kind 38891)

```json
{
  "kind": 38891,
  "tags": [
    ["d", "<book-uuid>"],
    ["a", "30023:<pubkey>:<ch1-uuid>"],
    ["a", "30023:<pubkey>:<ch2-uuid>"],
    ["p", "<author-pubkey>"],
    ["t", "book"]
  ],
  "content": "{\"title\":\"The Copper Islands\",\"subtitle\":\"...\",\"author\":\"NextBlock\",\"summary\":\"...\",\"description\":\"...\",\"image\":\"\"}"
}
```

Chapter ordering in book events follows `a` tag order (books reference chapters, not the other way around).

## Library Event (Kind 38890)

```json
{
  "kind": 38890,
  "tags": [
    ["d", "<library-uuid>"],
    ["a", "38891:<pubkey>:<book1-uuid>"],
    ["a", "38891:<pubkey>:<book2-uuid>"],
    ["t", "library"]
  ],
  "content": "{\"name\":\"the library\",\"description\":\"Books published on Nostr. Read by citizens.\"}"
}
```

## Engagement

Ketab-level engagement is the core requirement. Standard Nostr kinds:

| Kind | Purpose |
|------|---------|
| 1111 | Comments (NIP-22, threaded) |
| 9802 | Highlights (NIP-84) |
| 9735 | Zap receipts (NIP-57) |
| 6 | Reposts (NIP-18) |

All reference the ketab's `a` coordinate: `38893:<pubkey>:<ketab-uuid>`

## Client Discovery

1. Fetch book event (38891) by `naddr`
2. Extract chapter `a` tags → fetch chapters (30023)
3. For each chapter, fetch sibling ketabs (38893) via parent `a` tag
4. Sort ketabs by `index` field in content JSON
5. Each ketab is individually addressable via `naddr`

## City Protocol Integration

Ketab events can reference City Protocol block events for timestamps:

```
Block coordinate: 38808:<clock_pubkey>:org.cityprotocol:block:<height>:<hash>
```

Books published on Bitcoin time.

## Packages

- `@ketab/core` — TypeScript types and constants
- `@ketab/sdk` — Event construction, signing, publishing

## Related

- [City Protocol](https://github.com/joinnextblock/protocol-city) — Block-aware domains
- [ATTN Protocol](https://github.com/joinnextblock/protocol-attn) — Attention marketplace
- [Dynasty Protocol](https://github.com/joinnextblock/protocol-dynasty) — Sovereign genealogy

## License

MIT
