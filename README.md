# Ketab Protocol

Composable interactive stories on Nostr.

## The Problem

Publishing on Nostr today stops at the article. You can write a chapter, sign it, and publish it. Readers can zap the whole thing. But the ideas inside — the individual claims, the sourced passages, the moments that change someone's mind — are trapped inside the article body. You can't point to them. You can't engage with them independently. You can't pull one out and place it next to a passage from a different author without copy-pasting and breaking attribution.

Articles are containers. The ideas inside them have no identity of their own.

## The Ketab

A **ketab** is a signed Nostr event (kind 38893) that carries exactly one thought. One claim. One passage. 150–300 words with sourced footnotes.

Every ketab is:

- **Individually addressable** — its own `naddr`, its own URL, its own OG metadata
- **Independently engageable** — zap it, comment on it, highlight it, repost it
- **Cryptographically attributed** — signed by its author, verifiable by anyone
- **Composable** — a librarian can pull ketab 4 from one book and ketab 12 from another into a curated reading list. No copy-paste. No broken attribution. No aggregator in the middle.

Books are composed from ketabs. Ketabs can be recomposed across books. The atomic unit of a book is no longer the chapter — it's the thought.

## The Spirit

Ketab Protocol exists because ideas should be citable at the level they're expressed. A footnoted claim about Columbus shipping enslaved people from Hispaniola in 1495 deserves its own address — not because it's a tweet, but because it's a verifiable unit of knowledge that readers should be able to engage with, challenge, and build on.

The protocol makes no assumptions about how ketabs are rendered. Cards, scrolls, audio, AR — that's the client's job. The protocol's job is to make every thought addressable, every source verifiable, and every engagement attributable to the specific idea that earned it.

## Design Principles

**Single-letter tags only.** Tags exist for relay indexing: `d`, `a`, `p`, `e`, `t`. No `title`, `summary`, or `description` in tags. All metadata lives in the content field as structured JSON. This is the inverse of how most Nostr kinds work, and it's intentional — tags are for machines, content is for meaning.

Exception: kind 30023 chapters follow NIP-23 spec. Interop trumps our preferences.

**Position is explicit, never inferred.** Ketab ordering comes from the `index` field in content JSON, not from tag order, not from `created_at` timestamps, not from relay response order. If you can't trust the position, you can't trust the story.

**Engagement belongs to the thought, not the container.** When someone zaps a ketab, the zap targets that specific passage's coordinate — not the chapter, not the book. Authors see exactly which ideas resonate. Readers see exactly what they're paying for.

## Structure

Four event kinds, four layers:

| Kind | Name | What It Is |
|------|------|-----------|
| 38890 | Library | A librarian's curated collection of books |
| 38891 | Book | A book — metadata + ordered chapter references |
| 38893 | Ketab | One thought. The atomic unit. |
| 30023 | Chapter | NIP-23 long-form content (compiled view of ketabs) |

Three actors:

- **Authors** create books and ketabs
- **Librarians** curate books into libraries
- **Readers** engage per-ketab

A single pubkey can play all three roles.

## Specification

Full event schemas, tag definitions, content JSON formats, and discovery flows: **[PROTOCOL.md](./PROTOCOL.md)**

## Packages

- `@ketab/core` — TypeScript types and constants
- `@ketab/sdk` — Event construction, signing, publishing

## Related

- [City Protocol](https://github.com/joinnextblock/protocol-city) — Block-aware domains on Nostr
- [ATTN Protocol](https://github.com/joinnextblock/protocol-attn) — Decentralized attention marketplace
- [Dynasty Protocol](https://github.com/joinnextblock/protocol-dynasty) — Sovereign genealogy on Nostr

## License

MIT
