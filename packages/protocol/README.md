# Ketab Protocol

A Nostr protocol for organizing native books and library curation.

## Protocol Overview

**Problem:** NIP-23 (long-form content) chapters need organization into books, and libraries need rich metadata for curation—existing Nostr mechanisms don't provide both book organization and per-book library metadata.

**The Three-Actor Model:** Ketab Protocol is built on three distinct roles (note: a single pubkey owner can play multiple roles):
- **Authors** publish books (Book events, kind 38891) organizing NIP-23 chapters into books
- **Librarians** publish library entries (Library Entry events, kind 38892) to curate books with personal metadata (notes, rating, tags, read status)
- **Readers** publish comments (kind 1) and snippets (kind 9802) to engage with chapters

This three-actor model is the protocol's identity—each role has a clear purpose and publishes specific event types. Anyone can be a librarian by creating a library and curating books.

**Native Publishing Only:** Ketab Protocol only supports native Nostr book publishing. Books must be created and published directly on Nostr using NIP-23 chapters. The protocol does not support importing or bringing existing books from external sources (PDFs, EPUBs, etc.) onto Nostr. This is a core principle—all book content must be native to Nostr.

## Quick Start

The simplest flow to get started:

1. **Publish a Book**: An author publishes a Book event (kind 38891) organizing NIP-23 chapters into a book
2. **Create a Library**: Anyone can publish a Library event (kind 38890) to establish their personal library
3. **Curate a Book**: A librarian publishes a Library Entry event (kind 38892) to add the book to their library with personal metadata (notes, rating, tags, read status)

That's it. Books are independent entities that can be curated by multiple libraries, each with their own metadata. An author can also be a librarian (and vice versa)—roles are defined by what events you publish, not by identity.

## Core Rules

Three rules govern the entire system:

1. **Native Publishing Only**: All book content must be published natively on Nostr as NIP-23 chapters—no imports from external sources
2. **Author is Pubkey Owner**: The author's identity is determined by the event's `pubkey` field, not by metadata—the `author` field in BookContent is display-only
3. **Books are Independent**: Books are independent entities (collections of NIP-23 chapters) that can be curated by multiple libraries—libraries don't own books, they curate them

These rules create a system where books are published once by authors, libraries add curation metadata, and the protocol scales without data duplication.

## Protocol Hierarchy

Ketab Protocol organizes content in a clear hierarchy:

1. **Library** (kind 38890) - Curation container for Book events. Library events are solely for organizing and curating books.
2. **Book** (kind 38891) - Independent collections of NIP-23 chapters. Books can be curated by multiple libraries.
3. **Chapter** (NIP-23 kind 30023) - Long-form content within books.
4. **Engagement** - Chapters support comments (kind 1) and snippets (NIP-84 kind 9802).

**Key Design:** Books are independent entities (collections of 30023 events). Library events are curation containers that organize Book events. A single book can be curated by multiple libraries.

## Library Model

Ketab Protocol uses a **personal library model**:

- **Personal Libraries**: Anyone can create and maintain their own library (Library event kind 38890). Libraries are personal curation spaces—each pubkey owner can have their own library.
- **Book Publishing**: Authors publish books (Book event kind 38891). Books are independent collections of NIP-23 chapters. The author is the pubkey owner who publishes the book event. Books must be published natively on Nostr—no imports from external sources.
- **Library Curation**: Librarians curate books using Library Entry Events (kind 38892) with per-book metadata (notes, rating, tags, read status, etc.). When adding a book to a library, the app creates a Library Entry event. This avoids data duplication—books are published once by authors, and libraries add metadata about the book in their context. A single book can be curated by multiple libraries, each with their own metadata.
- **Chapter Organization**: Books organize NIP-23 long-form events as chapters
- **Chapter Engagement**: Chapters support comments (kind 1) and snippets (NIP-84 kind 9802) from readers
- **Public Reading**: Anyone can read books and chapters from any library (public access, no authentication required)
- **Decentralized Storage**: Events can be stored on any relay that supports Ketab Protocol events (kinds 38890, 38891, 38892)

### Access Control

**Publishing:**
- Anyone can publish books (Book events, kind 38891)
- Anyone can create a library (Library event, kind 38890) - the library belongs to the pubkey owner who publishes it
- Anyone can curate books by publishing Library Entry events (kind 38892)
- On NextBlock City relay: Only authenticated citizens can publish (NIP-42 authentication)
- On other relays: Access control is determined by the relay operator

**Reading:**
- Public access - anyone can read books and chapters from any library
- No authentication required for reading

**Relay Support:**
- Ketab Protocol events can be stored on any relay that supports the protocol
- NextBlock City relay supports Ketab Protocol events
- The protocol is open and decentralized, not tied to specific relay infrastructure

## Documentation

For complete technical specifications, see [LIBRARY-01.md](./LIBRARY-01.md).
