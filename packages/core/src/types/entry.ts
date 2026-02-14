/**
 * Library Entry Event (Kind 38892) Types
 */

/** Library Entry event content structure */
export interface LibraryEntryContent {
  /** Personal notes about this book (optional) */
  notes?: string;
  /** User's rating (1-5, etc.) (optional) */
  rating?: number;
  /** Personal tags for organization (optional) */
  tags?: string[];
  /** Unix timestamp when added to library */
  added_at: number;
  /** Read status (optional) */
  read_status?: string;
  /** Reference to library owner pubkey */
  ref_library_owner_pubkey: string;
  /** Reference to library ID */
  ref_library_id: string;
  /** Reference to book coordinate */
  ref_book_coordinate: string;
  /** Reference to book pubkey */
  ref_book_pubkey: string;
  /** Reference to book ID */
  ref_book_id: string;
  /** Reference to block event identifier */
  ref_block_id: string;
}

/** Library Entry identifier (d-tag value, scoped to librarian pubkey) */
export type LibraryEntryId = string;

/** Parsed Library Entry event */
export interface LibraryEntryEvent {
  /** Library Entry identifier from d tag */
  entry_id: LibraryEntryId;
  /** Library owner pubkey (librarian) */
  library_owner_pubkey: string;
  /** Event content */
  content: LibraryEntryContent;
  /** Nostr event ID */
  event_id: string;
  /** Event creation timestamp */
  created_at: number;
}
