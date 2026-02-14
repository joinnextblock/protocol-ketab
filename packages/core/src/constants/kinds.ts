/**
 * Ketab Protocol Event Kinds
 *
 * Reserved range: 38890-38892 for Ketab Protocol events
 */

/** Library Event - Book curation container (replaceable) */
export const KIND_LIBRARY = 38890;

/** Book Event - Book metadata and chapter organization (replaceable) */
export const KIND_BOOK = 38891;

/** Library Entry Event - Library-specific metadata about a curated book (replaceable) */
export const KIND_LIBRARY_ENTRY = 38892;

/** All Ketab Protocol event kinds */
export const LIBRARY_PROTOCOL_KINDS = [
  KIND_LIBRARY,
  KIND_BOOK,
  KIND_LIBRARY_ENTRY,
] as const;

/** Chapter identifier prefix (NIP-23) */
export const CHAPTER_ID_PREFIX: string = "30023:";
