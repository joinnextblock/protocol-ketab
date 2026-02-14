// Package core provides core types, constants, and utilities for Ketab Protocol.
// Ketab Protocol defines Nostr event kinds 38890-38892 for decentralized book libraries.
package core

// Event kinds for Ketab Protocol (38890-38892)
const (
	// KindLibrary is the event kind for Library events (38890).
	// Library events define book curation containers.
	KindLibrary = 38890

	// KindBook is the event kind for Book events (38891).
	// Book events define book metadata and chapter organization.
	KindBook = 38891

	// KindLibraryEntry is the event kind for Library Entry events (38892).
	// Library Entry events define library-specific metadata about curated books.
	KindLibraryEntry = 38892
)

// Protocol constants
const (
	// Version is the current Ketab Protocol version.
	Version = "0.1.0"

	// ChapterIDPrefix is the prefix for chapter identifiers (NIP-23).
	// Format: 30023:<author_pubkey>:<chapter_d-tag>
	ChapterIDPrefix = "30023:"
)

// KetabProtocolKinds contains all Ketab Protocol event kinds.
var KetabProtocolKinds = map[int]bool{
	KindLibrary:      true,
	KindBook:         true,
	KindLibraryEntry: true,
}

// IsKetabProtocolKind returns true if the kind is a Ketab Protocol event kind.
func IsKetabProtocolKind(kind int) bool {
	return KetabProtocolKinds[kind]
}
