# Ketab Protocol

Composable interactive stories on Nostr.

## The Problem

Nostr publishing stops at the article. You sign a chapter, publish it, readers zap the whole thing. But the ideas inside — the claims, the sourced passages, the moments that shift how someone sees the world — they're trapped in the body text. No address. No independent engagement. No way to pull one out and set it next to a passage from another author without copy-pasting and killing the attribution chain.

Articles are containers. The ideas inside have no identity.

## The Ketab

A **ketab** is kind 38893. One signed event. One thought. One passage with sourced footnotes.

Every ketab has its own `naddr`. Its own engagement — zaps, comments, highlights, reposts. Its own cryptographic attribution. You can verify who wrote it without trusting anyone.

And they compose. A librarian pulls ketab 4 from one book and ketab 12 from another into a curated reading list. No aggregator. No copy-paste. No broken signatures. The original author's event, referenced in place.

Books are made of ketabs. Ketabs can be recomposed across books. The atomic unit of a book is the thought, not the chapter.

## Why This Way

Tags route. Content describes. Single-letter tags only — `d`, `a`, `p`, `e`, `t` — for relay filtering. All metadata lives in the content field as structured JSON.

Most Nostr kinds do the opposite. They put title, summary, description in tags so relays can filter on them. That turns relays into queryable databases. We keep relays dumb. One JSON parse gives the client everything. No reassembling metadata from scattered tags.

Exception: kind 30023 chapters follow NIP-23. Interop trumps preferences.

Position is explicit. Ketab ordering comes from the `index` field in content JSON. Not tag order. Not timestamps. Not relay response order. If you can't trust the position, you can't trust the story.

Engagement belongs to the thought. When someone zaps a ketab, the zap targets that passage's coordinate. Not the chapter. Not the book. Authors see which ideas resonate. Readers see what they're paying for.

## Structure

Five event kinds:

| Kind | Name | What It Does |
|------|------|-------------|
| 38890 | Library | A curated collection of books |
| 38891 | Book | Metadata + ordered chapter references |
| 38892 | Library Entry | A book added to someone's personal library |
| 38893 | Ketab | One thought. The atomic unit. |
| 30023 | Chapter | NIP-23 long-form content (compiled view of ketabs) |

Three roles:

- **Authors** create books and ketabs
- **Librarians** curate books into libraries
- **Readers** engage per-ketab

One pubkey can play all three.

## Specification

Event schemas, tag definitions, content JSON formats, discovery flows: **[PROTOCOL.md](./PROTOCOL.md)**

## Packages

- `@ketab/core` — TypeScript types and constants
- `@ketab/sdk` — Event construction, signing, publishing

## Related

- [City Protocol](https://github.com/joinnextblock/protocol-city) — Block-aware domains on Nostr
- [ATTN Protocol](https://github.com/joinnextblock/protocol-attn) — Decentralized attention marketplace
- [Dynasty Protocol](https://github.com/joinnextblock/protocol-dynasty) — Sovereign genealogy on Nostr

## License

MIT
