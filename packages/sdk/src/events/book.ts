import type { BookContent, BookId } from "@ketab/core";
import type { BaseNostrEvent } from "@ketab/core";
import { KIND_BOOK } from "@ketab/core";
import { get_public_key } from "../signing/index.js";

/**
 * Build a Book Event (Kind 38891)
 */
export interface BuildBookEventOptions {
  /** Author's secret key (Uint8Array, 32 bytes) */
  secret_key: Uint8Array;
  /** Book ID (slug) */
  book_id: string;
  /** Book event content */
  content: BookContent;
  /** Unix timestamp (defaults to now) */
  created_at?: number;
}

/**
 * Build a Book Event
 */
export function build_book_event(
  options: BuildBookEventOptions
): Omit<BaseNostrEvent, "id" | "sig"> {
  const { secret_key, book_id, content, created_at = Math.floor(Date.now() / 1000) } = options;

  const book_identifier: BookId = book_id;

  const pubkey = get_public_key(secret_key);

  const tags: string[][] = [
    ["d", book_identifier],
  ];

  // Add chapter addresses for relay indexing
  for (const chapter_address of content.chapters) {
    tags.push(["a", chapter_address]);
  }

  // Add author pubkey
  tags.push(["p", pubkey]);

  // Add topic tags (genre, subject, etc.) - these would need to be passed separately
  // For now, we'll leave this to the caller to add if needed

  // Add relay URLs - these would need to be passed separately
  // For now, we'll leave this to the caller to add if needed

  return {
    kind: KIND_BOOK,
    pubkey,
    created_at,
    tags,
    content: JSON.stringify(content),
  };
}

/**
 * Build a book address tag (a tag) for referencing a book
 */
export function build_book_address(
  author_pubkey: string,
  book_id: string
): string {
  return `${KIND_BOOK}:${author_pubkey}:${book_id}`;
}
