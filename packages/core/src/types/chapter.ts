/**
 * Chapter Event (Kind 30023 - NIP-23) Types
 *
 * Chapters are NIP-23 long-form content events used by Ketab Protocol.
 * While NIP-23 allows title/published_at tags, Ketab Protocol stores
 * all metadata in content for consistency with other Ketab events.
 */

/** Chapter event content structure */
export interface ChapterContent {
  /** Chapter title */
  title: string;
  /** Published timestamp (Unix) */
  published_at: number;
  /** Chapter body (Markdown content) */
  body: string;
}

/** Chapter identifier format (d-tag value) */
export type ChapterId = string;

/** Chapter address format: 30023:<author_pubkey>:<chapter_d-tag> */
export type ChapterAddress = `30023:${string}:${string}`;

/** Parsed Chapter event */
export interface ChapterEvent {
  /** Chapter identifier from d tag */
  chapter_id: ChapterId;
  /** Author pubkey */
  author_pubkey: string;
  /** Event content */
  content: ChapterContent;
  /** Nostr event ID */
  event_id: string;
  /** Event creation timestamp */
  created_at: number;
}
