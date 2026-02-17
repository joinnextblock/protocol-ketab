---
name: nostr-ketab
description: Ketab Protocol (LIBRARY-01) for Nostr-native book publishing and library curation. Use when building book reader apps, creating book/library/library-entry events (kinds 38890-38892), organizing NIP-23 chapters into books, or implementing reader engagement features. Covers the three-actor model (Authors, Librarians, Readers) and City Protocol block-time integration.
---

# Ketab Protocol

Nostr-native book publishing and library curation.

## Core Concept

**Three-Actor Model** (a single pubkey can play multiple roles):

| Actor | Role | Events |
|-------|------|--------|
| **Authors** | Publish books organizing NIP-23 chapters | Book (38891) |
| **Librarians** | Curate books with personal metadata | Library (38890), Library Entry (38892) |
| **Readers** | Engage with chapters | Kind 1 (comments), 9802 (highlights), 9735 (zaps) |

**Native Publishing Only**: Books must be created on Nostr using NIP-23 chapters. No importing PDFs/EPUBs.

## Event Kinds

| Kind | Name | Description |
|------|------|-------------|
| 38890 | Library | Book curation container (replaceable) |
| 38891 | Book | Book metadata and chapter organization (replaceable) |
| 38892 | Library Entry | Library-specific metadata about a curated book (replaceable) |

All events are **parameterized replaceable events** (NIP-33).

## Kind 38891 - Book Event

Published by Authors. Organizes NIP-23 long-form chapters into a book.

```json
{
  "kind": 38891,
  "pubkey": "<author_pubkey>",
  "tags": [
    ["d", "my-book-slug"],
    ["title", "My Book Title"],
    ["summary", "A brief description of the book"],
    ["image", "https://example.com/cover.jpg"],
    ["a", "30023:<author_pubkey>:chapter-1", "relay-url", "Chapter 1 Title"],
    ["a", "30023:<author_pubkey>:chapter-2", "relay-url", "Chapter 2 Title"],
    ["a", "30023:<author_pubkey>:chapter-3", "relay-url", "Chapter 3 Title"]
  ],
  "content": ""
}
```

**Tags:**
- `d`: Unique book identifier (slug)
- `title`: Book title
- `summary`: Book description
- `image`: Cover image URL
- `a`: Ordered chapter references (NIP-23 addressable events, kind 30023)

## Kind 38890 - Library Event

Published by Librarians. A container for curated book collections.

```json
{
  "kind": 38890,
  "pubkey": "<librarian_pubkey>",
  "tags": [
    ["d", "my-library"],
    ["name", "My Reading Collection"],
    ["description", "Books I've curated"]
  ],
  "content": ""
}
```

## Kind 38892 - Library Entry Event

Published by Librarians. Personal metadata about a book in their library.

```json
{
  "kind": 38892,
  "pubkey": "<librarian_pubkey>",
  "tags": [
    ["d", "<library_d_tag>:<book_coordinate>"],
    ["a", "38891:<author_pubkey>:<book_d_tag>"],
    ["a", "38890:<librarian_pubkey>:<library_d_tag>"]
  ],
  "content": "{\"notes\":\"My thoughts on this book\",\"rating\":5,\"status\":\"reading\",\"tags\":[\"bitcoin\",\"economics\"]}"
}
```

**Content JSON:**
- `notes`: Personal notes about the book
- `rating`: 1-5 rating
- `status`: `want-to-read`, `reading`, `finished`, `abandoned`
- `tags`: Personal categorization tags

## Reader Engagement (Standard Nostr Kinds)

Chapters (NIP-23, kind 30023) support:

| Kind | Name | Use |
|------|------|-----|
| 1 | Text Note | Reader comments on chapters |
| 9802 | Highlight | Reader snippets/quotes (NIP-84) |
| 9735 | Zap Receipt | Zap reactions (NIP-57) |
| 6 | Repost | Sharing books/chapters (NIP-18) |

## Queries

```javascript
// Find all books by an author
{ kinds: [38891], authors: [author_pubkey] }

// Find a specific book
{ kinds: [38891], "#d": ["book-slug"], authors: [author_pubkey] }

// Find all libraries by a user
{ kinds: [38890], authors: [librarian_pubkey] }

// Find all library entries for a book
{ kinds: [38892], "#a": ["38891:<author_pubkey>:<book_d_tag>"] }

// Find highlights on a chapter
{ kinds: [9802], "#a": ["30023:<author_pubkey>:<chapter_d_tag>"] }
```

## City Protocol Integration

Library events can reference City block events for timing:

```
Block Event Coordinate: 38808:<clock_pubkey>:org.cityprotocol:block:<height>:<hash>
```

## Product Context

**Library** is the reader app for Ketab Protocol. **Publisher** is the authoring tool that guides users through creating books — especially interview-based family histories that link to Dynasty Protocol members.

Citizens bring their reading history and book collections to any app that supports Ketab Protocol.

## References

- [Ketab Protocol GitHub](https://github.com/joinnextblock/protocol-ketab)
- [NIP-23: Long-form Content](https://github.com/nostr-protocol/nips/blob/master/23.md)
- [NIP-33: Parameterized Replaceable Events](https://github.com/nostr-protocol/nips/blob/master/33.md)
- [NIP-84: Highlights](https://github.com/nostr-protocol/nips/blob/master/84.md)