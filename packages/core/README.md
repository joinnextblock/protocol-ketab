# @ketab/core

Core constants and types for Ketab Protocol.

## Installation

```bash
npm install @ketab/core
```

## Usage

```typescript
import { KIND_LIBRARY, KIND_BOOK, KIND_LIBRARY_ENTRY } from "@ketab/core";
import type { LibraryContent, BookContent, LibraryEntryContent } from "@ketab/core";

// Use constants
const library_kind = KIND_LIBRARY; // 38890

// Use types
const library_content: LibraryContent = {
  name: "My Library",
  description: "A personal library",
  // ...
};
```

## Exports

### Constants

- `KIND_LIBRARY` - Library event kind (38890)
- `KIND_BOOK` - Book event kind (38891)
- `KIND_LIBRARY_ENTRY` - Library Entry event kind (38892)
- `LIBRARY_PROTOCOL_KINDS` - Array of all Ketab Protocol event kinds
- `LIBRARY_NAMESPACE` - Protocol namespace prefix
- `LIBRARY_ID_PREFIX` - Library identifier prefix
- `BOOK_ID_PREFIX` - Book identifier prefix
- `LIBRARY_ENTRY_ID_PREFIX` - Library Entry identifier prefix
- `CHAPTER_ID_PREFIX` - Chapter identifier prefix (NIP-23)

### Types

- `LibraryContent` - Library event content structure
- `LibraryId` - Library identifier type
- `LibraryAddress` - Library address type
- `LibraryEvent` - Parsed Library event
- `BookContent` - Book event content structure
- `BookId` - Book identifier type
- `BookAddress` - Book address type
- `ChapterAddress` - Chapter address type
- `BookEvent` - Parsed Book event
- `LibraryEntryContent` - Library Entry event content structure
- `LibraryEntryId` - Library Entry identifier type
- `LibraryEntryEvent` - Parsed Library Entry event
- `BaseNostrEvent` - Base Nostr event structure
