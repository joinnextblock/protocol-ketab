/**
 * Book Event (Kind 38891) Types
 */

/** Book event content structure */
export interface BookContent {
  /** Book title */
  title: string;
  /** Book subtitle (optional) */
  subtitle?: string;
  /** Book description */
  description: string;
  /** Book dedication (optional) */
  dedication?: string;
  /** Author display name (metadata only - author identity is event's pubkey) */
  author: string;
  /** Cover image URL (optional) */
  cover_image_url?: string;
  /** Published timestamp (Unix) */
  published_at: number;
  /** Number of chapters */
  chapter_count: number;
  /** Ordered array of chapter addresses */
  chapters: string[];
  /** Reference to book pubkey (must match event's pubkey) */
  ref_book_pubkey: string;
  /** Reference to book ID */
  ref_book_id: string;
  /** Reference to library pubkey (optional - author's primary library) */
  ref_library_pubkey?: string;
  /** Reference to library ID (optional - author's primary library) */
  ref_library_id?: string;
  /** Reference to block event identifier */
  ref_block_id: string;
}

/** Book identifier format: org.ketab-protocol:book:<book_slug> */
export type BookId = `org.ketab-protocol:book:${string}`;

/** Book address format: 38891:<author_pubkey>:org.ketab-protocol:book:<book_slug> */
export type BookAddress = `${number}:${string}:org.ketab-protocol:book:${string}`;

/** Chapter address format: 30023:<author_pubkey>:<chapter_d-tag> */
export type ChapterAddress = `30023:${string}:${string}`;

/** Parsed Book event */
export interface BookEvent {
  /** Book identifier from d tag */
  book_id: BookId;
  /** Author pubkey */
  author_pubkey: string;
  /** Event content */
  content: BookContent;
  /** Nostr event ID */
  event_id: string;
  /** Event creation timestamp */
  created_at: number;
}
