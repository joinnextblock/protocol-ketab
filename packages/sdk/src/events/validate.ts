/**
 * Input validation for SDK event builders.
 * Mirrors Go validation logic — fail fast before building events.
 */

export class KetabValidationError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "KetabValidationError";
  }
}

function require_string(value: unknown, field: string): asserts value is string {
  if (typeof value !== "string" || value.trim().length === 0) {
    throw new KetabValidationError(`${field} is required and must be a non-empty string`);
  }
}

function require_number(value: unknown, field: string): asserts value is number {
  if (typeof value !== "number" || isNaN(value)) {
    throw new KetabValidationError(`${field} is required and must be a number`);
  }
}

function require_pubkey(value: unknown, field: string): void {
  require_string(value, field);
  if ((value as string).length !== 64 || !/^[0-9a-f]+$/.test(value as string)) {
    throw new KetabValidationError(`${field} must be a 64-character lowercase hex pubkey`);
  }
}

function require_secret_key(value: unknown): void {
  if (!(value instanceof Uint8Array) || value.length !== 32) {
    throw new KetabValidationError("secret_key must be a 32-byte Uint8Array");
  }
}

/** Validate library event builder inputs */
export function validate_library_inputs(options: {
  secret_key: Uint8Array;
  library_id: string;
  content: Record<string, unknown>;
}): void {
  require_secret_key(options.secret_key);
  require_string(options.library_id, "library_id");

  const c = options.content;
  require_string(c.name, "content.name");
  require_string(c.description, "content.description");
  require_string(c.founder_pubkey, "content.founder_pubkey");
  require_pubkey(c.founder_pubkey, "content.founder_pubkey");
  require_string(c.protocol_version, "content.protocol_version");
  require_string(c.ref_library_pubkey, "content.ref_library_pubkey");
  require_string(c.ref_library_id, "content.ref_library_id");
  require_string(c.ref_clock_pubkey, "content.ref_clock_pubkey");
}

/** Validate book event builder inputs */
export function validate_book_inputs(options: {
  secret_key: Uint8Array;
  book_id: string;
  content: Record<string, unknown>;
}): void {
  require_secret_key(options.secret_key);
  require_string(options.book_id, "book_id");

  const c = options.content;
  require_string(c.title, "content.title");
  require_string(c.description, "content.description");
  require_string(c.author, "content.author");
  require_number(c.published_at, "content.published_at");
  if (!Array.isArray(c.shape) || c.shape.length === 0) {
    throw new KetabValidationError("content.shape must be a non-empty array of acts");
  }

  require_string(c.ref_book_pubkey, "content.ref_book_pubkey");
  require_pubkey(c.ref_book_pubkey, "content.ref_book_pubkey");
  require_string(c.ref_book_id, "content.ref_book_id");
}

/** Validate library entry event builder inputs */
export function validate_entry_inputs(options: {
  secret_key: Uint8Array;
  library_owner_pubkey: string;
  book_slug: string;
  book_author_pubkey: string;
  content: Record<string, unknown>;
}): void {
  require_secret_key(options.secret_key);
  require_pubkey(options.library_owner_pubkey, "library_owner_pubkey");
  require_string(options.book_slug, "book_slug");
  require_pubkey(options.book_author_pubkey, "book_author_pubkey");

  const c = options.content;
  require_number(c.added_at, "content.added_at");
  require_string(c.ref_library_owner_pubkey, "content.ref_library_owner_pubkey");
  require_string(c.ref_library_id, "content.ref_library_id");
  require_string(c.ref_book_coordinate, "content.ref_book_coordinate");
  require_string(c.ref_book_pubkey, "content.ref_book_pubkey");
  require_string(c.ref_book_id, "content.ref_book_id");
}
