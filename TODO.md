# TODO — protocol-ketab

## High Priority

- [ ] **Remove unused deps from @ketab/core** — `nostr-tools` and `zod` are declared but never imported. Core is types-only.
- [ ] **Add input validation to SDK builders** — `build_*_event()` functions don't validate required content fields. Mirror Go's validation logic or add Zod schemas.
- [ ] **Add unit tests** — No tests exist in any package. Priority: Go `validation/validation.go`, SDK event builders, round-trip sign/verify.

## Medium Priority

- [ ] **Standardize import extensions** — SDK uses `.js`, core uses `.ts`. Pick one.
- [ ] **Fix `BookContent.Validate()` error** — Returns `ErrMissingName` for empty `Author`. Add `ErrMissingAuthor`.
- [ ] **Refactor Go validation to use typed structs** — `validation.go` unmarshals to `map[string]interface{}` instead of using `core.*Content` structs + their `Validate()` methods.

## Low Priority

- [ ] **Add go-core README** — Only package without a README.
- [ ] **Document missing chapter builder** — SDK has no `build_chapter_event`. If intentional (NIP-23 is standard), document why.
- [ ] **Clarify `CHAPTER_ID_PREFIX` naming** — It's a kind prefix (`30023:`), not a namespace prefix like the others.
