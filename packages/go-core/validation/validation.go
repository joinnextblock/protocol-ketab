// Package validation provides event validators for Ketab Protocol.
package validation

import (
	"encoding/json"
	"fmt"
	"strings"

	core "github.com/joinnextblock/ketab-protocol/go-core"
	"github.com/nbd-wtf/go-nostr"
)

// ValidationResult represents the result of validating an event.
type ValidationResult struct {
	Valid   bool
	Message string
}

// ValidateEvent validates a Ketab Protocol event based on its kind.
// Returns a ValidationResult indicating if the event is valid.
func ValidateEvent(event *nostr.Event) ValidationResult {
	switch event.Kind {
	case core.KindLibrary:
		return ValidateLibraryEvent(event)
	case core.KindBook:
		return ValidateBookEvent(event)
	case core.KindLibraryEntry:
		return ValidateLibraryEntryEvent(event)
	default:
		return ValidationResult{Valid: false, Message: fmt.Sprintf("Unknown Ketab Protocol kind: %d", event.Kind)}
	}
}

// ValidateLibraryEvent validates Library events (kind 38890) per LIBRARY-01 specification.
//
// Required tags:
//   - d: Library identifier (any non-empty string, scoped to pubkey)
//   - a: Block coordinate (format: 38808:<clock_pubkey>:...)
//   - p: Library pubkey and clock pubkey (at least 2 p tags)
//
// Required content fields:
//   - name, description, founder_pubkey, protocol_version
//   - ref_library_pubkey, ref_library_id, ref_clock_pubkey, ref_block_id
//   - book_count, reader_count, chapter_count (can be 0)
func ValidateLibraryEvent(event *nostr.Event) ValidationResult {
	if event.Kind != core.KindLibrary {
		return ValidationResult{Valid: false, Message: fmt.Sprintf("Expected kind %d for Library event, got %d", core.KindLibrary, event.Kind)}
	}

	// Must have d tag (library identifier)
	d_tag := get_tag_value(event, "d")
	if d_tag == "" {
		return ValidationResult{Valid: false, Message: "Missing 'd' tag (library identifier)"}
	}

	// Must have block coordinate a tag (format: 38808:clock_pubkey:org.cityprotocol:block:<height>:<hash>)
	a_tags := get_tag_values(event, "a")
	has_block_coord := false
	for _, a_tag := range a_tags {
		if strings.HasPrefix(a_tag, "38808:") {
			has_block_coord = true
			break
		}
	}
	if !has_block_coord {
		return ValidationResult{Valid: false, Message: "Missing block coordinate 'a' tag (format: 38808:clock_pubkey:org.cityprotocol:block:<height>:<hash>)"}
	}

	// Must have p tags (library pubkey and clock pubkey)
	p_tags := get_tag_values(event, "p")
	if len(p_tags) < 2 {
		return ValidationResult{Valid: false, Message: "Missing required 'p' tags (need library pubkey and clock pubkey)"}
	}

	// Content must be valid JSON
	var content_data map[string]interface{}
	if err := json.Unmarshal([]byte(event.Content), &content_data); err != nil {
		return ValidationResult{Valid: false, Message: "Content must be valid JSON"}
	}

	// Check required string fields
	required_string_fields := []string{
		"name", "description", "founder_pubkey", "protocol_version",
		"ref_library_pubkey", "ref_library_id", "ref_clock_pubkey", "ref_block_id",
	}
	for _, field := range required_string_fields {
		val, ok := content_data[field]
		if !ok {
			return ValidationResult{Valid: false, Message: fmt.Sprintf("Content must include '%s'", field)}
		}
		if str, ok := val.(string); !ok || str == "" {
			return ValidationResult{Valid: false, Message: fmt.Sprintf("Content field '%s' must be a non-empty string", field)}
		}
	}

	// Check required count fields (can be 0)
	required_count_fields := []string{"book_count", "reader_count", "chapter_count"}
	for _, field := range required_count_fields {
		if _, ok := content_data[field]; !ok {
			return ValidationResult{Valid: false, Message: fmt.Sprintf("Content must include '%s'", field)}
		}
	}

	return ValidationResult{Valid: true, Message: "Valid Library event"}
}

// ValidateBookEvent validates Book events (kind 38891) per LIBRARY-01 specification.
//
// Required tags:
//   - d: Book identifier (any non-empty string, scoped to pubkey)
//   - p: Author pubkey
//
// Required content fields:
//   - title, description, author, published_at, chapter_count, chapters
//   - ref_book_pubkey, ref_book_id, ref_block_id
//
// Validation: ref_book_pubkey must match event's pubkey
func ValidateBookEvent(event *nostr.Event) ValidationResult {
	if event.Kind != core.KindBook {
		return ValidationResult{Valid: false, Message: fmt.Sprintf("Expected kind %d for Book event, got %d", core.KindBook, event.Kind)}
	}

	// Must have d tag (book identifier)
	d_tag := get_tag_value(event, "d")
	if d_tag == "" {
		return ValidationResult{Valid: false, Message: "Missing 'd' tag (book identifier)"}
	}

	// Must have p tag (author pubkey)
	p_tags := get_tag_values(event, "p")
	if len(p_tags) == 0 {
		return ValidationResult{Valid: false, Message: "Missing required 'p' tag (author pubkey)"}
	}

	// Content must be valid JSON
	var content_data map[string]interface{}
	if err := json.Unmarshal([]byte(event.Content), &content_data); err != nil {
		return ValidationResult{Valid: false, Message: "Content must be valid JSON"}
	}

	// Check required string fields
	required_string_fields := []string{
		"title", "description", "author",
		"ref_book_pubkey", "ref_book_id", "ref_block_id",
	}
	for _, field := range required_string_fields {
		val, ok := content_data[field]
		if !ok {
			return ValidationResult{Valid: false, Message: fmt.Sprintf("Content must include '%s'", field)}
		}
		if str, ok := val.(string); !ok || str == "" {
			return ValidationResult{Valid: false, Message: fmt.Sprintf("Content field '%s' must be a non-empty string", field)}
		}
	}

	// Check published_at (required, numeric)
	if _, ok := content_data["published_at"]; !ok {
		return ValidationResult{Valid: false, Message: "Content must include 'published_at'"}
	}

	// Check chapter_count (required, numeric)
	if _, ok := content_data["chapter_count"]; !ok {
		return ValidationResult{Valid: false, Message: "Content must include 'chapter_count'"}
	}

	// Check chapters array (required)
	chapters, ok := content_data["chapters"]
	if !ok {
		return ValidationResult{Valid: false, Message: "Content must include 'chapters' array"}
	}
	if _, ok := chapters.([]interface{}); !ok {
		return ValidationResult{Valid: false, Message: "Content field 'chapters' must be an array"}
	}

	// Validate ref_book_pubkey matches event's pubkey
	ref_book_pubkey, _ := content_data["ref_book_pubkey"].(string)
	if ref_book_pubkey != event.PubKey {
		return ValidationResult{Valid: false, Message: "ref_book_pubkey must match event's pubkey (author identity)"}
	}

	return ValidationResult{Valid: true, Message: "Valid Book event"}
}

// ValidateLibraryEntryEvent validates Library Entry events (kind 38892) per LIBRARY-01 specification.
//
// Required tags:
//   - d: Entry identifier (any non-empty string, scoped to pubkey)
//   - a: Book coordinate and library coordinate (at least 2 a tags)
//   - p: Library owner pubkey and book author pubkey (at least 2 p tags)
//
// Required content fields:
//   - added_at
//   - ref_library_owner_pubkey, ref_library_id, ref_book_coordinate, ref_book_pubkey, ref_book_id, ref_block_id
func ValidateLibraryEntryEvent(event *nostr.Event) ValidationResult {
	if event.Kind != core.KindLibraryEntry {
		return ValidationResult{Valid: false, Message: fmt.Sprintf("Expected kind %d for Library Entry event, got %d", core.KindLibraryEntry, event.Kind)}
	}

	// Must have d tag (entry identifier)
	d_tag := get_tag_value(event, "d")
	if d_tag == "" {
		return ValidationResult{Valid: false, Message: "Missing 'd' tag (entry identifier)"}
	}

	// Must have a tags (book coordinate and library coordinate)
	a_tags := get_tag_values(event, "a")
	has_book_coord := false
	has_library_coord := false
	for _, a_tag := range a_tags {
		if strings.HasPrefix(a_tag, "38891:") {
			has_book_coord = true
		}
		if strings.HasPrefix(a_tag, "38890:") {
			has_library_coord = true
		}
	}
	if !has_book_coord {
		return ValidationResult{Valid: false, Message: "Missing book coordinate 'a' tag (format: 38891:<author_pubkey>:<book_id>)"}
	}
	if !has_library_coord {
		return ValidationResult{Valid: false, Message: "Missing library coordinate 'a' tag (format: 38890:<library_owner_pubkey>:<library_id>)"}
	}

	// Must have p tags (library owner pubkey and book author pubkey)
	p_tags := get_tag_values(event, "p")
	if len(p_tags) < 2 {
		return ValidationResult{Valid: false, Message: "Missing required 'p' tags (need library owner pubkey and book author pubkey)"}
	}

	// Content must be valid JSON
	var content_data map[string]interface{}
	if err := json.Unmarshal([]byte(event.Content), &content_data); err != nil {
		return ValidationResult{Valid: false, Message: "Content must be valid JSON"}
	}

	// Check required string fields
	required_string_fields := []string{
		"ref_library_owner_pubkey", "ref_library_id", "ref_book_coordinate",
		"ref_book_pubkey", "ref_book_id", "ref_block_id",
	}
	for _, field := range required_string_fields {
		val, ok := content_data[field]
		if !ok {
			return ValidationResult{Valid: false, Message: fmt.Sprintf("Content must include '%s'", field)}
		}
		if str, ok := val.(string); !ok || str == "" {
			return ValidationResult{Valid: false, Message: fmt.Sprintf("Content field '%s' must be a non-empty string", field)}
		}
	}

	// Check added_at (required, numeric)
	if _, ok := content_data["added_at"]; !ok {
		return ValidationResult{Valid: false, Message: "Content must include 'added_at'"}
	}

	return ValidationResult{Valid: true, Message: "Valid Library Entry event"}
}

// Helper functions

// get_tag_value returns the first value of a tag with the given name.
func get_tag_value(event *nostr.Event, tag_name string) string {
	for _, tag := range event.Tags {
		if len(tag) >= 2 && tag[0] == tag_name {
			return tag[1]
		}
	}
	return ""
}

// get_tag_values returns all values of tags with the given name.
func get_tag_values(event *nostr.Event, tag_name string) []string {
	var values []string
	for _, tag := range event.Tags {
		if len(tag) >= 2 && tag[0] == tag_name {
			values = append(values, tag[1])
		}
	}
	return values
}
