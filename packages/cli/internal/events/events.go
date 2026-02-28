// Package events provides functions for building Nostr events from book data.
package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/joinnextblock/ketab-protocol/cli/internal/book"
	"github.com/joinnextblock/ketab-protocol/cli/internal/types"
	core "github.com/joinnextblock/ketab-protocol/go-core"
	"github.com/nbd-wtf/go-nostr"
)

// Event kind constants.
const (
	KindKetab   = 38893 // Individual scene/passage
	KindChapter = 30023 // NIP-23 long-form content
	KindBook    = core.KindBook
	KindLibrary = core.KindLibrary
)

// Builder builds Nostr events for a book.
type Builder struct {
	pubkey     string
	relay_hint string
}

// NewBuilder creates a new event builder.
func NewBuilder(pubkey string, relay_hint string) *Builder {
	return &Builder{
		pubkey:     pubkey,
		relay_hint: relay_hint,
	}
}

// BuildKetab builds a ketab event (kind 38893).
func (b *Builder) BuildKetab(ch *book.Chapter, ketab book.Ketab) nostr.Event {
	content := types.KetabContent{
		Title: ketab.Item.Title,
		Index: ketab.Item.Number - 1, // 0-based
		Ord:   ketab.Item.Number,     // 1-based
		Body:  ketab.Body,
	}

	content_json, _ := json.Marshal(content)

	return nostr.Event{
		Kind:      KindKetab,
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Tags: nostr.Tags{
			{"d", ketab.Item.UUID},
		},
		Content: string(content_json),
	}
}

// BuildChapter builds a chapter event (kind 30023).
func (b *Builder) BuildChapter(bk *book.Book, ch *book.Chapter) nostr.Event {
	// Compile full chapter body from all ketabs
	body := ch.CompileChapterBody()

	// Build tags
	tags := nostr.Tags{
		{"d", ch.Metadata.ChapterUUID},
		{"title", fmt.Sprintf("Chapter %s: %s", ch.Metadata.ChapterNumber, ch.Metadata.ChapterTitle)},
		{"published_at", fmt.Sprintf("%d", time.Now().Unix())},
		// Reference parent book
		{"a", fmt.Sprintf("%d:%s:%s", KindBook, b.pubkey, bk.Metadata.BookUUID), b.relay_hint},
	}

	// Reference ketabs
	for _, ketab := range ch.Ketabs {
		tags = append(tags, nostr.Tag{
			"a", fmt.Sprintf("%d:%s:%s", KindKetab, b.pubkey, ketab.Item.UUID), b.relay_hint,
		})
	}

	return nostr.Event{
		Kind:      KindChapter,
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Tags:      tags,
		Content:   body,
	}
}

// BookContent is the JSON content for a book event.
type BookContent struct {
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Author        string       `json:"author"`
	CoverImageURL string       `json:"cover_image_url,omitempty"`
	PublishedAt   int64        `json:"published_at"`
	Shape         any          `json:"shape,omitempty"` // From book-shape.json (optional)
	Acts          []PublishAct `json:"acts"`            // 3-level hierarchy: acts → chapters → ketabs
	RefBookPubkey string       `json:"ref_book_pubkey"`
	RefBookID     string       `json:"ref_book_id"`
}

// PublishAct represents an act in the published book event.
type PublishAct struct {
	Title    string           `json:"title"`
	Chapters []PublishChapter `json:"chapters"` // Always an array, empty if no published chapters
}

// PublishChapter represents a chapter in the published book event.
type PublishChapter struct {
	Number string         `json:"number"`
	Title  string         `json:"title"`
	UUID   string         `json:"uuid"`
	Ketabs []PublishKetab `json:"ketabs"` // Always an array, empty if no ketabs
}

// PublishKetab represents a ketab reference in the published book event.
type PublishKetab struct {
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}

// BuildBook builds a book event (kind 38891).
func (b *Builder) BuildBook(bk *book.Book, chapter_nums []string) nostr.Event {
	now := time.Now().Unix()

	// Build the 3-level acts hierarchy from book metadata
	published_chapters := make(map[string]bool)
	for _, ch_num := range chapter_nums {
		published_chapters[ch_num] = true
	}

	var acts []PublishAct
	for _, act_ref := range bk.Metadata.Acts {
		publishAct := PublishAct{
			Title:    act_ref.Title,
			Chapters: []PublishChapter{}, // Initialize as empty array, not nil
		}

		// Add chapters from this act
		for _, ch_ref := range act_ref.Chapters {
			publishChapter := PublishChapter{
				Number: ch_ref.ChapterNumber,
				Title:  ch_ref.ChapterTitle,
				UUID:   ch_ref.ChapterUUID,
				Ketabs: []PublishKetab{}, // Initialize as empty array, not nil
			}

			// Only add ketabs if this chapter is being published
			if published_chapters[ch_ref.ChapterNumber] {
				if ch, ok := bk.GetChapter(ch_ref.ChapterNumber); ok {
					for _, ketab := range ch.Ketabs {
						publishChapter.Ketabs = append(publishChapter.Ketabs, PublishKetab{
							Title: ketab.Item.Title,
							UUID:  ketab.Item.UUID,
						})
					}
				}
			}
			
			publishAct.Chapters = append(publishAct.Chapters, publishChapter)
		}

		acts = append(acts, publishAct)
	}

	// Build content
	content := BookContent{
		Title:         bk.Metadata.BookTitle,
		Description:   bk.Metadata.Description,
		Author:        bk.Metadata.Author,
		CoverImageURL: bk.Metadata.Image,
		PublishedAt:   now,
		Acts:          acts, // Use the constructed 3-level hierarchy
		RefBookPubkey: b.pubkey,
		RefBookID:     bk.Metadata.BookUUID,
	}

	// Optionally include book-shape.json data if available
	if bk.Shape != nil {
		content.Shape = bk.Shape.Shape
	}

	content_json, _ := json.Marshal(content)

	// Build tags
	tags := nostr.Tags{
		{"d", bk.Metadata.BookUUID},
		{"title", bk.Metadata.BookTitle},
		{"p", b.pubkey},
	}

	if bk.Metadata.Image != "" {
		tags = append(tags, nostr.Tag{"image", bk.Metadata.Image})
	}
	if bk.Metadata.Thumb != "" {
		tags = append(tags, nostr.Tag{"thumb", bk.Metadata.Thumb})
	}
	if bk.Metadata.Summary != "" {
		tags = append(tags, nostr.Tag{"summary", bk.Metadata.Summary})
	}

	// Add chapter references
	for _, ch_num := range chapter_nums {
		if ch, ok := bk.GetChapter(ch_num); ok {
			addr := fmt.Sprintf("%d:%s:%s", KindChapter, b.pubkey, ch.Metadata.ChapterUUID)
			tags = append(tags, nostr.Tag{"a", addr, b.relay_hint})
		}
	}

	return nostr.Event{
		Kind:      KindBook,
		CreatedAt: nostr.Timestamp(now),
		Tags:      tags,
		Content:   string(content_json),
	}
}

// LibraryContent is the JSON content for a library event.
type LibraryContent struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Books       []string `json:"books"`
}

// BuildLibrary builds a library event (kind 38890).
func (b *Builder) BuildLibrary(bk *book.Book, library_id string, library_name string) nostr.Event {
	book_coord := fmt.Sprintf("%d:%s:%s", KindBook, b.pubkey, bk.Metadata.BookUUID)

	content := LibraryContent{
		Name:        library_name,
		Description: "Books published on Nostr. Read by citizens.",
		Books:       []string{book_coord},
	}

	content_json, _ := json.Marshal(content)

	tags := nostr.Tags{
		{"d", library_id},
		{"title", library_name},
		{"a", book_coord, b.relay_hint},
	}

	return nostr.Event{
		Kind:      KindLibrary,
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Tags:      tags,
		Content:   string(content_json),
	}
}

// SignEvent signs an event with the given secret key.
func SignEvent(event *nostr.Event, sk string) error {
	return event.Sign(sk)
}
