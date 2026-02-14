import type { LibraryEntryContent, LibraryEntryId } from "@ketab/core";
import type { BaseNostrEvent } from "@ketab/core";
import { KIND_LIBRARY_ENTRY, KIND_BOOK, KIND_LIBRARY } from "@ketab/core";
import { get_public_key } from "../signing/index.js";

/**
 * Build a Library Entry Event (Kind 38892)
 */
export interface BuildLibraryEntryEventOptions {
  /** Librarian's secret key (Uint8Array, 32 bytes) */
  secret_key: Uint8Array;
  /** Library owner pubkey */
  library_owner_pubkey: string;
  /** Book slug */
  book_slug: string;
  /** Library Entry event content */
  content: LibraryEntryContent;
  /** Book author pubkey */
  book_author_pubkey: string;
  /** Unix timestamp (defaults to now) */
  created_at?: number;
}

/**
 * Build a Library Entry Event
 */
export function build_library_entry_event(
  options: BuildLibraryEntryEventOptions
): Omit<BaseNostrEvent, "id" | "sig"> {
  const { secret_key, library_owner_pubkey, book_slug, content, book_author_pubkey, created_at = Math.floor(Date.now() / 1000) } = options;

  const entry_identifier: LibraryEntryId = `${library_owner_pubkey}:${book_slug}`;

  const pubkey = get_public_key(secret_key);

  // Build book coordinate
  const book_coordinate = `${KIND_BOOK}:${book_author_pubkey}:${book_slug}`;

  // Build library coordinate
  const library_coordinate = `${KIND_LIBRARY}:${library_owner_pubkey}:${content.ref_library_id}`;

  const tags: string[][] = [
    ["d", entry_identifier],
    ["a", book_coordinate], // Book coordinate
    ["a", library_coordinate], // Library coordinate
    ["p", library_owner_pubkey], // Librarian
    ["p", book_author_pubkey], // Book author
  ];

  return {
    kind: KIND_LIBRARY_ENTRY,
    pubkey,
    created_at,
    tags,
    content: JSON.stringify(content),
  };
}
