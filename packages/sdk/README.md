# @ketab/sdk

TypeScript SDK for creating and publishing Ketab Protocol events.

## Installation

```bash
npm install @ketab/sdk
```

## Usage

```typescript
import { build_library_event, build_book_event, build_library_entry_event } from "@ketab/sdk";
import { sign_event } from "@ketab/sdk";
import type { LibraryContent, BookContent, LibraryEntryContent } from "@ketab/core";

// Build and sign a Library event
const library_content: LibraryContent = {
  name: "My Library",
  description: "A personal library",
  founder_pubkey: "abc123...",
  protocol_version: "0.1.0",
  ref_library_pubkey: "abc123...",
  ref_library_id: "my-library",
  ref_clock_pubkey: "clock_pubkey...",
  ref_block_id: "org.cityprotocol:block:862626:...",
  book_count: 0,
  reader_count: 0,
  chapter_count: 0,
};

const library_event_template = build_library_event({
  secret_key: my_secret_key,
  library_id: "my-library",
  content: library_content,
});

const signed_library_event = sign_event(library_event_template, my_secret_key);

// Build and sign a Book event
const book_content: BookContent = {
  title: "My Book",
  description: "A great book",
  author: "Author Name",
  published_at: Math.floor(Date.now() / 1000),
  chapter_count: 2,
  chapters: [
    "30023:author_pubkey:chapter-1",
    "30023:author_pubkey:chapter-2",
  ],
  ref_book_pubkey: "author_pubkey...",
  ref_book_id: "my-book",
  ref_block_id: "org.cityprotocol:block:862626:...",
};

const book_event_template = build_book_event({
  secret_key: author_secret_key,
  book_id: "my-book",
  content: book_content,
});

const signed_book_event = sign_event(book_event_template, author_secret_key);
```

## Exports

### Event Builders

- `build_library_event` - Build a Library event (kind 38890)
- `build_book_event` - Build a Book event (kind 38891)
- `build_library_entry_event` - Build a Library Entry event (kind 38892)
- `build_library_address` - Build a library address tag
- `build_book_address` - Build a book address tag

### Signing Utilities

- `sign_event` - Sign a Nostr event with a secret key
- `get_public_key` - Get public key from secret key
- `verify_event` - Verify an event signature
- `hex_to_secret_key` - Convert hex string to Uint8Array
- `secret_key_to_hex` - Convert Uint8Array to hex string
