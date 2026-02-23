package core

import (
	"errors"
)

var (
	// ErrInvalidVersion is returned when the protocol version is not supported.
	ErrInvalidVersion = errors.New("unsupported protocol version")

	// ErrMissingName is returned when a required name field is empty.
	ErrMissingName = errors.New("name is required")

	// ErrMissingTitle is returned when a required title field is empty.
	ErrMissingTitle = errors.New("title is required")

	// ErrMissingDescription is returned when a required description field is empty.
	ErrMissingDescription = errors.New("description is required")

	// ErrMissingFounderPubkey is returned when founder_pubkey is missing.
	ErrMissingFounderPubkey = errors.New("founder_pubkey is required")

	// ErrMissingRefLibraryID is returned when ref_library_id is missing.
	ErrMissingRefLibraryID = errors.New("ref_library_id is required")

	// ErrMissingRefBookID is returned when ref_book_id is missing.
	ErrMissingRefBookID = errors.New("ref_book_id is required")

)

// LibraryContent represents the content structure for Library events (kind 38890).
// Per LIBRARY-01 specification.
type LibraryContent struct {
	// Name is the library name.
	Name string `json:"name"`

	// Description is the library description.
	Description string `json:"description"`

	// WebsiteURL is the library website URL (optional).
	WebsiteURL string `json:"website_url,omitempty"`

	// RelayURL is the library relay URL (optional).
	RelayURL string `json:"relay_url,omitempty"`

	// FounderPubkey is the librarian's public key.
	FounderPubkey string `json:"founder_pubkey"`

	// ProtocolVersion is the protocol version.
	ProtocolVersion string `json:"protocol_version"`

	// RefLibraryPubkey is the reference to library pubkey.
	RefLibraryPubkey string `json:"ref_library_pubkey"`

	// RefLibraryID is the reference to library ID.
	RefLibraryID string `json:"ref_library_id"`

	// RefClockPubkey is the reference to City Protocol clock pubkey.
	RefClockPubkey string `json:"ref_clock_pubkey"`

	// BookCount is the total number of books (can be 0).
	BookCount int `json:"book_count"`

	// ReaderCount is the total number of unique readers (can be 0).
	ReaderCount int `json:"reader_count"`
}

// Validate checks if the LibraryContent has required fields per LIBRARY-01.
func (l *LibraryContent) Validate() error {
	if l.Name == "" {
		return ErrMissingName
	}
	if l.Description == "" {
		return ErrMissingDescription
	}
	if l.FounderPubkey == "" {
		return ErrMissingFounderPubkey
	}
	if l.RefLibraryID == "" {
		return ErrMissingRefLibraryID
	}
	return nil
}

// BookShapeKetab represents a ketab reference in the book shape.
type BookShapeKetab struct {
	Title string `json:"title"`
	DTag  string `json:"d_tag"`
}

// BookShapeChapter represents a chapter in the book shape.
type BookShapeChapter struct {
	Title        string           `json:"title"`
	DTag         string           `json:"d_tag"`
	DiscussionID string           `json:"discussion_id,omitempty"`
	Ketabs       []BookShapeKetab `json:"ketabs"`
}

// BookContent represents the content structure for Book events (kind 38891).
// Per LIBRARY-01 specification.
type BookContent struct {
	// Title is the book title.
	Title string `json:"title"`

	// Subtitle is the book subtitle (optional).
	Subtitle string `json:"subtitle,omitempty"`

	// Description is the book description.
	Description string `json:"description"`

	// Dedication is the book dedication (optional).
	Dedication string `json:"dedication,omitempty"`

	// Author is the author display name (metadata only - author identity is event's pubkey).
	Author string `json:"author"`

	// CoverImageURL is the cover image URL (optional).
	CoverImageURL string `json:"cover_image_url,omitempty"`

	// PublishedAt is the published timestamp (Unix).
	PublishedAt int64 `json:"published_at"`

	// Shape is the book structure: an array of acts, each act is an array of chapters.
	Shape [][]BookShapeChapter `json:"shape"`

	// RefBookPubkey is the reference to book pubkey (must match event's pubkey).
	RefBookPubkey string `json:"ref_book_pubkey"`

	// RefBookID is the reference to book ID.
	RefBookID string `json:"ref_book_id"`

	// RefLibraryPubkey is the reference to library pubkey (optional - author's primary library).
	RefLibraryPubkey string `json:"ref_library_pubkey,omitempty"`

	// RefLibraryID is the reference to library ID (optional - author's primary library).
	RefLibraryID string `json:"ref_library_id,omitempty"`

}

// Validate checks if the BookContent has required fields per LIBRARY-01.
func (b *BookContent) Validate() error {
	if b.Title == "" {
		return ErrMissingTitle
	}
	if b.Description == "" {
		return ErrMissingDescription
	}
	if b.Author == "" {
		return ErrMissingName
	}
	if b.RefBookID == "" {
		return ErrMissingRefBookID
	}
	return nil
}

// LibraryEntryContent represents the content structure for Library Entry events (kind 38892).
// Per LIBRARY-01 specification.
type LibraryEntryContent struct {
	// Notes is personal notes about this book (optional).
	Notes string `json:"notes,omitempty"`

	// Rating is the user's rating (1-5, etc.) (optional).
	Rating *int `json:"rating,omitempty"`

	// Tags is personal tags for organization (optional).
	Tags []string `json:"tags,omitempty"`

	// AddedAt is the Unix timestamp when added to library.
	AddedAt int64 `json:"added_at"`

	// ReadStatus is the read status (optional): "unread", "reading", "completed", etc.
	ReadStatus string `json:"read_status,omitempty"`

	// RefLibraryOwnerPubkey is the reference to library owner pubkey.
	RefLibraryOwnerPubkey string `json:"ref_library_owner_pubkey"`

	// RefLibraryID is the reference to library ID.
	RefLibraryID string `json:"ref_library_id"`

	// RefBookCoordinate is the reference to book coordinate.
	// Format: 38891:<author_pubkey>:<book_id>
	RefBookCoordinate string `json:"ref_book_coordinate"`

	// RefBookPubkey is the reference to book pubkey.
	RefBookPubkey string `json:"ref_book_pubkey"`

	// RefBookID is the reference to book ID.
	RefBookID string `json:"ref_book_id"`

}

// Validate checks if the LibraryEntryContent has required fields per LIBRARY-01.
func (e *LibraryEntryContent) Validate() error {
	if e.RefLibraryOwnerPubkey == "" {
		return ErrMissingFounderPubkey
	}
	if e.RefLibraryID == "" {
		return ErrMissingRefLibraryID
	}
	if e.RefBookID == "" {
		return ErrMissingRefBookID
	}
	return nil
}
