import type { LibraryContent, LibraryId } from "@ketab/core";
import type { BaseNostrEvent } from "@ketab/core";
import { KIND_LIBRARY, LIBRARY_ID_PREFIX } from "@ketab/core";
import { get_public_key } from "../signing/index.js";

/**
 * Build a Library Event (Kind 38890)
 */
export interface BuildLibraryEventOptions {
  /** Librarian's secret key (Uint8Array, 32 bytes) */
  secret_key: Uint8Array;
  /** Library ID (slug) */
  library_id: string;
  /** Library event content */
  content: LibraryContent;
  /** Unix timestamp (defaults to now) */
  created_at?: number;
}

/**
 * Build a Library Event
 */
export function build_library_event(
  options: BuildLibraryEventOptions
): Omit<BaseNostrEvent, "id" | "sig"> {
  const { secret_key, library_id, content, created_at = Math.floor(Date.now() / 1000) } = options;

  // Validate library ID format
  const library_identifier: LibraryId = `${LIBRARY_ID_PREFIX}${library_id}`;

  const pubkey = get_public_key(secret_key);

  const tags: string[][] = [
    ["d", library_identifier],
    ["a", content.ref_block_id], // City Protocol block event coordinate
    ["p", pubkey], // Librarian
    ["p", content.ref_clock_pubkey], // City clock reference
  ];

  // Add relay URLs if provided
  if (content.relay_url) {
    tags.push(["r", content.relay_url]);
  }

  // Add website URL if provided
  if (content.website_url) {
    tags.push(["u", content.website_url]);
  }

  return {
    kind: KIND_LIBRARY,
    pubkey,
    created_at,
    tags,
    content: JSON.stringify(content),
  };
}

/**
 * Build a library address tag (a tag) for referencing a library
 */
export function build_library_address(
  founder_pubkey: string,
  library_id: string
): string {
  return `${KIND_LIBRARY}:${founder_pubkey}:${LIBRARY_ID_PREFIX}${library_id}`;
}
