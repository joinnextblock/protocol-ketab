# protocol-ketab

Ketab Protocol - book publishing on Nostr (kind 38890-38892).

## Stack

- TypeScript + Go monorepo
- npm workspaces

## Commands

```bash
npm install               # Install dependencies
npm run build             # Build all packages
npm run typecheck         # Type check all packages
npm run test              # Run tests
```

## Packages

| Package | Language | Description |
|---------|----------|-------------|
| `core` | TS | Core types and constants |
| `sdk` | TS | Event creation SDK |
| `go-core` | Go | Core types for Go |

## Event Kinds

- `38890` - Book metadata
- `38891` - Chapter content
- `38892` - Book index
