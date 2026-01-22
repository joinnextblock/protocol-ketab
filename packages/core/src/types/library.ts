/**
 * Library Event (Kind 38890) Types
 */

/** Library event content structure */
export interface LibraryContent {
  /** Library name */
  name: string;
  /** Library description */
  description: string;
  /** Library website URL (optional) */
  website_url?: string;
  /** Library relay URL (optional) */
  relay_url?: string;
  /** Founder pubkey (librarian) */
  founder_pubkey: string;
  /** Protocol version */
  protocol_version: string;
  /** Reference to library pubkey */
  ref_library_pubkey: string;
  /** Reference to library ID */
  ref_library_id: string;
  /** Reference to clock pubkey (City Protocol) */
  ref_clock_pubkey: string;
  /** Reference to block event identifier */
  ref_block_id: string;
  /** Total books (can be 0) */
  book_count: number;
  /** Total unique readers (can be 0) */
  reader_count: number;
  /** Total chapters across all books (can be 0) */
  chapter_count: number;
}

/** Library identifier format: org.ketab-protocol:library:<library_id> */
export type LibraryId = `org.ketab-protocol:library:${string}`;

/** Library address format: 38890:<founder_pubkey>:org.ketab-protocol:library:<library_id> */
export type LibraryAddress = `${number}:${string}:org.ketab-protocol:library:${string}`;

/** Parsed Library event */
export interface LibraryEvent {
  /** Library identifier from d tag */
  library_id: LibraryId;
  /** Founder pubkey (librarian) */
  founder_pubkey: string;
  /** Event content */
  content: LibraryContent;
  /** Nostr event ID */
  event_id: string;
  /** Event creation timestamp */
  created_at: number;
}
