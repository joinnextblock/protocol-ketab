# Ketab Protocol

Ketab Protocol (LIBRARY-01) is a Nostr-based protocol for organizing native books and library curation.

## Core Concept

**Three-Actor Model**: Ketab Protocol is built on three distinct roles (a single pubkey owner can play multiple roles):
- **Authors** publish books (Book events, kind 38891) organizing NIP-23 chapters into books
- **Librarians** publish library entries (Library Entry events, kind 38892) to curate books with personal metadata (notes, rating, tags, read status)
- **Readers** publish comments (kind 1) and snippets (kind 9802) to engage with chapters

**Native Publishing Only**: Ketab Protocol only supports native Nostr book publishing. Books must be created and published directly on Nostr using NIP-23 chapters. The protocol does not support importing or bringing existing books from external sources (PDFs, EPUBs, etc.) onto Nostr.

## Event Kinds

| Kind | Name | Description |
|------|------|-------------|
| 38890 | Library Event | Book curation container (replaceable) |
| 38891 | Book Event | Book metadata and chapter organization (replaceable) |
| 38892 | Library Entry | Library-specific metadata about a curated book (replaceable) |

### Standard Nostr Kinds for Engagement

Chapters (NIP-23 events) support engagement through standard Nostr kinds:
- **Kind 1** - Text Note (reader comments)
- **Kind 9802** - Highlights (reader snippets/quotes, NIP-84)
- **Kind 9735** - Zap Receipt (zap reactions, NIP-57)
- **Kind 6** - Repost (sharing books/chapters, NIP-18)

## Protocol Integration

### NextBlock City Integration

Ketab Protocol is integrated into NextBlock City through the `KetabAdapter`. When enabled, the city automatically subscribes to Ketab Protocol events (kinds 38890-38892) and provides hooks for library, book, and library entry events.

**Configuration:**
```typescript
const city = new NextBlockCity({
  signer: window.nostr,
  protocols: {
    ketab: true,  // Enable Ketab Protocol hooks (default: true)
  },
});

await city.enter();

// Access Ketab adapter
const ketab_adapter = city.ketab_adapter;  // If exposed

// Register handlers
ketab_adapter.on_library_event(async (ctx) => {
  console.log('Library:', ctx.content.name);
});

ketab_adapter.on_book_event(async (ctx) => {
  console.log('Book:', ctx.content.title);
});

ketab_adapter.on_library_entry_event(async (ctx) => {
  console.log('Library Entry:', ctx.content.notes);
});

ketab_adapter.start();
```

**Integration Details:**
- Adapter: `KetabAdapter` in `@nextblock/city/packages/city/src/adapters/ketab.ts` (to be implemented)
- Event Kinds: 38890 (Library), 38891 (Book), 38892 (Library Entry)
- Hooks: `on_library_event()`, `on_book_event()`, `on_library_entry_event()`

**Integration Status:** Ketab Protocol integration is planned. The adapter will be implemented following the pattern established by ATTN, City, and Dynasty protocols.

See [NextBlock City Protocol Integration Guide](../nextblock-city/PROTOCOL_INTEGRATION.md) for details on how protocols integrate.

### City Protocol Integration

Ketab Protocol integrates with City Protocol for block synchronization. Library events reference City block events for timing:

```
Block Event Coordinate: 38808:<clock_pubkey>:org.cityprotocol:block:<height>:<hash>
```

This allows libraries to operate on Bitcoin time without needing their own block event infrastructure.

## Packages

This monorepo contains the core Ketab Protocol packages:

- `@ketab/protocol` - LIBRARY-01 specification (markdown docs only)
- `@ketab/core` - TypeScript types and constants (no runtime code)
- `@ketab/sdk` - Event construction utilities, Nostr signing/publishing

## Development

This project uses npm workspaces for local development.

### Setup

```bash
npm install
```

### Type Check

```bash
npm run typecheck
```

### Working with Workspaces

- Install dependency in specific package: `npm install <package> --workspace=@ketab/core`
- Run script in specific package: `npm run <script> --workspace=@ketab/core`
- Run script in all packages: `npm run <script>` (from root)

## Related Projects

- [City Protocol](https://github.com/joinnextblock/city-protocol) - Block-aware domain protocol that Ketab Protocol references for timing
- [ATTN Protocol](https://github.com/joinnextblock/attn-protocol) - Decentralized attention marketplace
- [Nostr Protocol](https://github.com/nostr-protocol/nips)

## License

MIT
