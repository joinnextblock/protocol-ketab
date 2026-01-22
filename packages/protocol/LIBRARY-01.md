# LIBRARY-01: Ketab Protocol Specification

This document provides the complete technical specification for Ketab Protocol. For an overview and quick start guide, see [README.md](./README.md).

## Event Kinds

| Kind | Name | Description |
|------|------|-------------|
| 38890 | Library Event | Book curation container (replaceable) |
| 38891 | Book Event | Book metadata and chapter organization (replaceable) |
| 38892 | Library Entry | Library-specific metadata about a curated book (replaceable) |

### Standard Nostr Kinds for Engagement

Chapters (NIP-23 events) support engagement through standard Nostr kinds:

| Kind | Name | Description |
|------|------|-------------|
| 1 | Text Note | Reader comments on chapters - NIP-01 |
| 9802 | Highlights | Reader snippets/quotes from chapters - NIP-84 |
| 9735 | Zap Receipt | Zap reactions - NIP-57 |
| 6 | Repost | Sharing books/chapters - NIP-18 |

---

## LIBRARY Event (kind 38890)

**Purpose**: Curate book events. Library events are solely for organizing and curating Book events (kind 38891). Libraries can curate books published by the librarian and books from other authors.

**Published By**: Librarians (founder_pubkey = librarian's pubkey)

**Curation Model**: Library events are curation containers for Book events. Books are discovered via queries (see below), and Library events provide the organizational context for those books.

**Schema**:

```typescript
interface LibraryContent {
  // Library-specific fields
  name: string;
  description: string;
  website_url?: string;
  relay_url?: string;
  founder_pubkey: string;
  protocol_version: string;

  // Reference fields (ref_ prefix)
  ref_library_pubkey: string;
  ref_library_id: string;
  ref_clock_pubkey: string;  // Clock pubkey this library listens to
  ref_block_id: string;  // Current block event identifier (org.cityprotocol:block:<height>:<hash>)

  // Metrics fields (required, can be 0)
  book_count: number;  // Total books (can be 0)
  reader_count: number;  // Total unique readers (can be 0)
  chapter_count: number;  // Total chapters across all books (can be 0)
}
```

**Tags**:

```typescript
[
  ["d", "org.ketab-protocol:library:<library_id>"],
  ["a", "38808:<clock_pubkey>:org.cityprotocol:block:<height>:<hash>"],  // City Protocol block event coordinate
  ["p", "<library_pubkey>"],
  ["p", "<clock_pubkey>"],  // City clock reference
  ["r", "<relay_url>"],  // Multiple
  ["u", "<website_url>"]  // Optional
]
```

**Note:** Libraries do not list books in the Library event. Books are curated using:
- **Books published by librarian**: Query for Book events where `pubkey` matches the librarian's pubkey
- **Curated books from other authors**: Query for Library Entry events (kind 38892) where `pubkey` matches the librarian

This keeps Library events lightweight and scalable, even for libraries with thousands of books.

**Example**:

```json
{
  "kind": 38890,
  "pubkey": "library_pubkey_hex",
  "created_at": 1234567890,
  "tags": [
    ["d", "org.ketab-protocol:library:nextblock-library"],
    ["a", "38808:clock_pubkey_hex:org.cityprotocol:block:862626:00000000000000000001a7c..."],
    ["p", "author_pubkey_hex"],  // Librarian (author)
    ["p", "clock_pubkey_hex"],
    ["r", "wss://relay.nextblock.city"],
    ["u", "https://library.nextblock.city"]
  ],
  "content": "{\"name\":\"NextBlock Library\",\"description\":\"Decentralized library of books\",\"website_url\":\"https://library.nextblock.city\",\"relay_url\":\"wss://relay.nextblock.city\",\"founder_pubkey\":\"author_pubkey_hex\",\"protocol_version\":\"0.1.0\",\"ref_library_pubkey\":\"author_pubkey_hex\",\"ref_library_id\":\"nextblock-library\",\"ref_clock_pubkey\":\"clock_pubkey_hex\",\"ref_block_id\":\"org.cityprotocol:block:862626:00000000000000000001a7c...\",\"book_count\":42,\"reader_count\":128,\"chapter_count\":350}"
}
```

**Note:** In a personal library model, `founder_pubkey` and `ref_library_pubkey` are the same (the author's pubkey). The library belongs to the author who publishes it.

---

## BOOK Event (kind 38891)

**Purpose**: Define book metadata and organize NIP-23 chapters into a book. Books are independent entities that can be referenced by multiple libraries.

**Published By**: Authors (books are published independently, not tied to a specific library)

**Author Identity:** The author is fundamentally the pubkey owner. The event's `pubkey` field identifies the author. The `author` field in BookContent is a display name/metadata string, but the author's identity is determined by the event's `pubkey`. The `ref_book_pubkey` must match the event's `pubkey`.

**Native Publishing Requirement:** Books must be published natively on Nostr using NIP-23 chapters. The protocol does not support importing existing books from external sources (e.g., PDFs, EPUBs, or other formats). All book content must be created and published directly as Nostr events.

**Note:** Books are collections of NIP-23 chapters. A book can exist in multiple libraries - libraries can curate books from any author. The `ref_library_pubkey` and `ref_library_id` fields in BookContent are optional and represent the author's primary library, but books can be referenced by other libraries as well.

**Schema**:

```typescript
interface BookContent {
  // Book-specific fields
  title: string;
  subtitle?: string;
  description: string;
  dedication?: string;
  author: string;  // Display name/metadata - author identity is the event's pubkey
  cover_image_url?: string;
  published_at: number;  // Unix timestamp
  chapter_count: number;
  chapters: string[];  // Ordered array of chapter addresses: ["30023:<author_pubkey>:<chapter_d-tag>", ...]

  // Reference fields (ref_ prefix)
  ref_book_pubkey: string;  // Must match event's pubkey (author's pubkey)
  ref_book_id: string;
  ref_library_pubkey?: string;  // Optional - author's primary library
  ref_library_id?: string;  // Optional - author's primary library
  ref_block_id: string;  // Block event identifier when book was published/updated
}
```

**Tags**:

```typescript
[
  ["d", "org.ketab-protocol:book:<book_slug>"],
  ["a", "38890:<library_pubkey>:org.ketab-protocol:library:<library_id>"],  // Optional - author's primary library (can be referenced by multiple libraries)
  ["a", "30023:<author_pubkey>:<chapter_d-tag>"],  // Chapter addresses (multiple, for relay indexing only)
  // ... more chapters
  ["p", "<author_pubkey>"],
  ["p", "<library_pubkey>"],  // If library reference exists
  ["t", "<topic_tag>"],  // Multiple - genre, subject, etc.
  ["r", "<relay_url>"],  // Multiple
  ["u", "<book_url>"]  // Optional
]
```

**Note:** Tags are for relay indexing only. Chapter order is determined by the `chapters` array in the content field.

**Example**:

```json
{
  "kind": 38891,
  "pubkey": "author_pubkey_hex",
  "created_at": 1234567890,
  "tags": [
    ["d", "org.ketab-protocol:book:my-book-title"],
    ["a", "38890:author_pubkey_hex:org.ketab-protocol:library:nextblock-library"],  // Author's library
    ["a", "30023:author_pubkey_hex:chapter-1"],
    ["a", "30023:author_pubkey_hex:chapter-2"],
    ["p", "author_pubkey_hex"],
    ["t", "fiction"],
    ["t", "sci-fi"],
    ["r", "wss://relay.nextblock.city"],
    ["u", "https://library.nextblock.city/books/my-book-title"]
  ],
  "content": "{\"title\":\"My Book Title\",\"subtitle\":\"A Science Fiction Novel\",\"description\":\"A story about the future\",\"author\":\"Author Name\",\"cover_image_url\":\"https://example.com/cover.jpg\",\"published_at\":1234567890,\"chapter_count\":2,\"chapters\":[\"30023:author_pubkey_hex:chapter-1\",\"30023:author_pubkey_hex:chapter-2\"],\"ref_book_pubkey\":\"author_pubkey_hex\",\"ref_book_id\":\"my-book-title\",\"ref_library_pubkey\":\"author_pubkey_hex\",\"ref_library_id\":\"nextblock-library\",\"ref_block_id\":\"org.cityprotocol:block:862626:00000000000000000001a7c...\"}"
}
```

---

## Chapter Engagement

Chapters (NIP-23 long-form events) support two primary forms of engagement:

1. **Comments** (kind 1) - Reader discussions and feedback on chapters
2. **Snippets** (kind 9802, NIP-84) - Reader quotes and highlights from chapters

Both reference chapters via `a` tags using the chapter coordinate format: `30023:<author_pubkey>:<chapter_d-tag>`

---

## Comments (Standard Nostr Kind 1 - NIP-01)

**Published By**: Readers

**Schema**: Uses standard NIP-01 text note format

**Tags**:

```typescript
[
  ["e", "<parent_comment_event_id>"],  // Parent comment event ID (for threading, not for chapters)
  ["a", "30023:<author_pubkey>:<chapter_d-tag>"],  // Chapter address (if commenting on chapter)
  ["a", "38891:<author_pubkey>:org.ketab-protocol:book:<book_slug>"],  // Book reference
  ["p", "<author_pubkey>"],  // Book/chapter author
  ["p", "<mentioned_pubkey>"]  // Multiple - mentioned users
]
```

**Content**: Plain text or Markdown comment

**Example**:

```json
{
  "kind": 1,
  "pubkey": "commenter_pubkey_hex",
  "created_at": 1234567890,
  "tags": [
    ["a", "30023:author_pubkey_hex:chapter-5"],
    ["a", "38891:author_pubkey_hex:org.ketab-protocol:book:my-book-title"],
    ["e", "parent_comment_event_id_hex"],  // If replying to another comment
    ["p", "author_pubkey_hex"],
    ["p", "mentioned_user_pubkey_hex"]
  ],
  "content": "Great chapter! I loved the character development."
}
```

---

## Snippets (NIP-84 Highlights - Kind 9802)

**Purpose**: Reader quotes and highlights from chapters. Snippets create standalone quotes that reference specific chapters.

**Published By**: Readers

**Schema**: Uses standard NIP-84 highlight format

**Tags**:

```typescript
[
  ["a", "30023:<author_pubkey>:<chapter_d-tag>"],  // Chapter address
  ["a", "38891:<author_pubkey>:org.ketab-protocol:book:<book_slug>"],  // Book reference (optional)
  ["p", "<author_pubkey>"],  // Chapter author
  ["context", "<surrounding_text>"]  // Optional - surrounding text for context
]
```

**Content**: The quoted text snippet from the chapter

**Example**:

```json
{
  "kind": 9802,
  "pubkey": "reader_pubkey_hex",
  "created_at": 1234567890,
  "tags": [
    ["a", "30023:author_pubkey_hex:chapter-5"],
    ["a", "38891:author_pubkey_hex:org.ketab-protocol:book:my-book-title"],
    ["p", "author_pubkey_hex"],
    ["context", "The surrounding paragraph text for context..."]
  ],
  "content": "This is the exact quoted text from the chapter that the reader highlighted."
}
```

**Note:** Snippets are self-contained quotes. The quoted text is stored in the event's content field, making it a standalone reference even if the chapter content changes later.

---

## Library Curation

Libraries curate books using Library Entry Events (kind 38892).

## Library Entry Events (Kind 38892)

**Purpose**: Curate books with library-specific metadata (notes, tags, rating, added_at, etc.).

**Published By**: Librarians

**Schema**:

```typescript
interface LibraryEntryContent {
  // Library-specific metadata about the book
  notes?: string;  // Personal notes about this book
  rating?: number;  // User's rating (1-5, etc.)
  tags?: string[];  // Personal tags for organization
  added_at: number;  // Unix timestamp when added to library
  read_status?: string;  // "unread", "reading", "completed", etc.

  // Reference fields (ref_ prefix)
  ref_library_owner_pubkey: string;
  ref_library_id: string;
  ref_book_coordinate: string;  // 38891:<author_pubkey>:org.ketab-protocol:book:<book_slug>
  ref_book_pubkey: string;
  ref_book_id: string;
  ref_block_id: string;  // Block event identifier when entry was created/updated
}
```

**Tags**:

```typescript
[
  ["d", "org.ketab-protocol:entry:<library_owner_pubkey>:<book_slug>"],
  ["a", "38891:<author_pubkey>:org.ketab-protocol:book:<book_slug>"],  // Book coordinate
  ["a", "38890:<library_owner_pubkey>:org.ketab-protocol:library:<library_id>"],  // Library coordinate
  ["p", "<library_owner_pubkey>"],
  ["p", "<book_author_pubkey>"]
]
```

**Example**:

```json
{
  "kind": 38892,
  "pubkey": "library_owner_pubkey_hex",
  "created_at": 1234567890,
  "tags": [
    ["d", "org.ketab-protocol:entry:library_owner_pubkey_hex:my-book-title"],
    ["a", "38891:author_pubkey_hex:org.ketab-protocol:book:my-book-title"],
    ["a", "38890:library_owner_pubkey_hex:org.ketab-protocol:library:my-library"],
    ["p", "library_owner_pubkey_hex"],
    ["p", "author_pubkey_hex"]
  ],
  "content": "{\"notes\":\"Great sci-fi novel, loved the ending\",\"rating\":5,\"tags\":[\"favorite\",\"sci-fi\"],\"added_at\":1234567890,\"read_status\":\"completed\",\"ref_library_owner_pubkey\":\"library_owner_pubkey_hex\",\"ref_library_id\":\"my-library\",\"ref_book_coordinate\":\"38891:author_pubkey_hex:org.ketab-protocol:book:my-book-title\",\"ref_book_pubkey\":\"author_pubkey_hex\",\"ref_book_id\":\"my-book-title\",\"ref_block_id\":\"org.cityprotocol:block:862626:00000000000000000001a7c...\"}"
}
```

**Note:** Library Entry events add metadata *about* the book in the context of that library. The book data itself (title, chapters, etc.) comes from the Book event (38891) published by the author. No data duplication - just library-specific annotations. To discover all books curated by a librarian, query for all Library Entry events (kind 38892) where `pubkey` matches the librarian's pubkey.

---

## Coordinate Format

Ketab Protocol uses NIP-33 parameterized replaceable events with coordinates:
- Library: `38890:<founder_pubkey>:org.ketab-protocol:library:<library_id>`
- Book: `38891:<author_pubkey>:org.ketab-protocol:book:<book_slug>`
- Library Entry: `38892:<library_owner_pubkey>:org.ketab-protocol:entry:<library_owner_pubkey>:<book_slug>`
- Chapter (NIP-23): `30023:<author_pubkey>:<chapter_d-tag>`

## Tag Conventions

**Ketab Protocol Events (38890-38891):**
- `d` tag: Namespace identifier for replaceable events
- `a` tag: Addressable event references (libraries, books, chapters, blocks)
- `p` tag: Pubkey references (founder, author, reader)
- `t` tag: Topic tags (genre, subject, etc.)
- `r` tag: Relay URLs
- `u` tag: URLs (website, book URL)

**Standard Nostr Engagement Events (Kinds 1, 9802, 6, 9735):**
- `e` tag: Event references (parent comment event ID for threading only, not for chapters)
- `a` tag: Addressable references (book/chapter coordinates: `38891:author:book-slug` or `30023:author:chapter-slug`)
- `p` tag: Pubkey references (mentioned users, chapter author)
- `context` tag: Surrounding text context (NIP-84 snippets)
- `t` tag: Topic tags
- `zap` tag: Zap receipt (NIP-57, in zap receipt events)

**Library Entry Events (Kind 38892):**
- `d` tag: Entry identifier (`org.ketab-protocol:entry:<library_owner_pubkey>:<book_slug>`)
- `a` tag: Book coordinate (`38891:<author_pubkey>:org.ketab-protocol:book:<book_slug>`)
- `a` tag: Library coordinate (`38890:<library_owner_pubkey>:org.ketab-protocol:library:<library_id>`)
- `p` tag: Librarian pubkey
- `p` tag: Book author pubkey
