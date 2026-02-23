// Package types defines data structures for book directory parsing.
package types

// BookMetadata represents the book-metadata.json file structure.
type BookMetadata struct {
	BookTitle   string            `json:"book_title"`
	BookSlug    string            `json:"book_slug"`
	Author      string            `json:"author"`
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Thumb       string            `json:"thumb,omitempty"`
	Signer      string            `json:"signer,omitempty"`
	BookUUID    string            `json:"book_uuid"`
	Acts        []ActRef          `json:"acts"`
}

// GetAllChapters returns a flattened list of all chapters across all acts.
func (b *BookMetadata) GetAllChapters() []ChapterRef {
	var chapters []ChapterRef
	for _, act := range b.Acts {
		chapters = append(chapters, act.Chapters...)
	}
	return chapters
}

// ChapterRef is a reference to a chapter in book-metadata.json.
type ChapterRef struct {
	ChapterNumber string `json:"chapter_number"`
	ChapterTitle  string `json:"chapter_title"`
	ChapterUUID   string `json:"chapter_uuid"`
	Coordinate    string `json:"coordinate,omitempty"`
}

// ActRef is a reference to an act in book-metadata.json.
type ActRef struct {
	Title    string       `json:"title"`
	Chapters []ChapterRef `json:"chapters"`
}

// BookShape represents the book-shape.json file structure.
type BookShape struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Image       string           `json:"image,omitempty"`
	Shape       [][]ShapeChapter `json:"shape"`
}

// ShapeChapter is a chapter in book-shape.json.
type ShapeChapter struct {
	Title  string       `json:"title"`
	DTag   string       `json:"d_tag"`
	Ketabs []ShapeKetab `json:"ketabs"`
}

// ShapeKetab is a ketab reference in book-shape.json.
type ShapeKetab struct {
	Title string `json:"title"`
	DTag  string `json:"d_tag"`
}

// ChapterMetadata represents the chapter-metadata.json file structure.
// Supports both "ketabs" (old format) and "scenes" (new format) arrays.
type ChapterMetadata struct {
	ChapterTitle  string       `json:"chapter_title"`
	ChapterNumber string       `json:"chapter_number"`
	ChapterUUID   string       `json:"chapter_uuid"`
	Act           int          `json:"act,omitempty"`
	PublishedAt   int64        `json:"published_at,omitempty"`
	Ketabs        []KetabRef   `json:"ketabs,omitempty"`  // Old format
	Scenes        []SceneRef   `json:"scenes,omitempty"`  // New format
	Tags          [][]string   `json:"tags,omitempty"`
}

// KetabRef is a ketab reference in chapter-metadata.json (old format).
type KetabRef struct {
	KetabNumber int    `json:"ketab_number"`
	KetabFile   string `json:"ketab_file"`
	KetabTitle  string `json:"ketab_title"`
	KetabUUID   string `json:"ketab_uuid"`
}

// SceneRef is a scene reference in chapter-metadata.json (new format).
type SceneRef struct {
	SceneNumber int    `json:"scene_number"`
	SceneFile   string `json:"scene_file"`
	SceneTitle  string `json:"scene_title"`
	KetabUUID   string `json:"ketab_uuid"`
}

// GetKetabs returns a unified list of ketab items from either format.
func (c *ChapterMetadata) GetKetabs() []KetabItem {
	var items []KetabItem

	// Check scenes first (new format)
	if len(c.Scenes) > 0 {
		for _, s := range c.Scenes {
			items = append(items, KetabItem{
				Number: s.SceneNumber,
				File:   s.SceneFile,
				Title:  s.SceneTitle,
				UUID:   s.KetabUUID,
			})
		}
		return items
	}

	// Fall back to ketabs (old format)
	for _, k := range c.Ketabs {
		items = append(items, KetabItem{
			Number: k.KetabNumber,
			File:   k.KetabFile,
			Title:  k.KetabTitle,
			UUID:   k.KetabUUID,
		})
	}
	return items
}

// KetabItem is a unified ketab/scene item.
type KetabItem struct {
	Number int
	File   string
	Title  string
	UUID   string
}

// KetabContent is the JSON content for a ketab event (kind 38893).
type KetabContent struct {
	Title string `json:"title"`
	Index int    `json:"index"` // 0-based
	Ord   int    `json:"ord"`   // 1-based
	Body  string `json:"body"`
}

// PublishConfig holds publishing configuration.
type PublishConfig struct {
	SecretKey  string
	Relays     []string
	DryRun     bool
	ChapterNum string // Empty = all chapters
}

// DefaultRelays are the default relays for publishing.
var DefaultRelays = []string{
	"wss://relay.nextblock.city",
	"wss://relay.primal.net",
	"wss://nos.lol",
	"wss://relay.damus.io",
}
