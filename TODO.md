# TODO — protocol-ketab

*Last reviewed: 2026-02-22*

## High Priority

- [ ] **Add unit tests** — No tests exist in any package. Priority: Go `validation/validation.go`, SDK event builders, round-trip sign/verify.
- [ ] **Remove unused deps from @ketab/core** — `nostr-tools` and `zod` are declared but never imported. Core is types-only.
- [ ] **Add input validation to SDK builders** — `build_*_event()` functions don't validate required content fields. Mirror Go's validation logic or add Zod schemas.
- [ ] **Make library_id configurable in CLI** — Currently hardcoded UUID in `main.go:219`. Add flag/env/metadata source.

## Medium Priority

- [ ] **Refactor Go validation to use typed structs** — `validation.go` unmarshals to `map[string]interface{}` instead of using `core.*Content` structs + their `Validate()` methods.
- [ ] **Fix `BookContent.Validate()` error** — Returns `ErrMissingName` for empty `Author`. Add `ErrMissingAuthor`.
- [ ] **Standardize import extensions** — SDK uses `.js`, core uses `.ts`. Pick one.
- [ ] **Remove dead code** — `_ = time.Now()` in CLI main.go (time is already used).
- [ ] **Batch relay connections in CLI** — `publish_event()` opens/closes per relay per event. Reuse connections.

## Low Priority

- [ ] **Add go-core README** — Only package without a README.
- [ ] **Document missing chapter builder** — SDK has no `build_chapter_event`. If intentional (NIP-23 is standard), document why.
- [ ] **Update `interface{}` to `any`** — 6 occurrences in Go code.
- [ ] **Remove unused `godotenv` dependency** — In CLI go.mod but not imported.
