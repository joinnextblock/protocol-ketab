// Package book provides functions for loading and parsing book directories.
package book

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/joinnextblock/ketab-protocol/cli/internal/types"
)

// Book represents a loaded book directory.
type Book struct {
	Dir      string
	Metadata *types.BookMetadata
	Shape    *types.BookShape
	Chapters map[string]*Chapter // key = chapter number (e.g., "00", "01")
}

// Chapter represents a loaded chapter.
type Chapter struct {
	Dir      string
	Number   string
	Metadata *types.ChapterMetadata
	Ketabs   []Ketab
}

// Ketab represents a loaded ketab/scene.
type Ketab struct {
	Item types.KetabItem
	Body string // Markdown content with scene headers stripped
}

// Load loads a book from a directory.
func Load(dir string) (*Book, error) {
	abs_dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("invalid directory: %w", err)
	}

	book := &Book{
		Dir:      abs_dir,
		Chapters: make(map[string]*Chapter),
	}

	// Load book-metadata.json
	meta_path := filepath.Join(abs_dir, "book-metadata.json")
	meta_data, err := os.ReadFile(meta_path)
	if err != nil {
		return nil, fmt.Errorf("failed to read book-metadata.json: %w", err)
	}
	book.Metadata = &types.BookMetadata{}
	if err := json.Unmarshal(meta_data, book.Metadata); err != nil {
		return nil, fmt.Errorf("failed to parse book-metadata.json: %w", err)
	}

	// Load book-shape.json (optional)
	shape_path := filepath.Join(abs_dir, "book-shape.json")
	if shape_data, err := os.ReadFile(shape_path); err == nil {
		book.Shape = &types.BookShape{}
		if err := json.Unmarshal(shape_data, book.Shape); err != nil {
			return nil, fmt.Errorf("failed to parse book-shape.json: %w", err)
		}
	}

	// Load chapters from acts
	for _, ch_ref := range book.Metadata.GetAllChapters() {
		ch, err := load_chapter(abs_dir, ch_ref.ChapterNumber)
		if err != nil {
			// Skip chapters that don't exist on disk
			continue
		}
		book.Chapters[ch_ref.ChapterNumber] = ch
	}

	return book, nil
}

// load_chapter loads a single chapter from disk.
func load_chapter(book_dir, chapter_num string) (*Chapter, error) {
	ch_dir := filepath.Join(book_dir, chapter_num)
	meta_path := filepath.Join(ch_dir, "chapter-metadata.json")

	meta_data, err := os.ReadFile(meta_path)
	if err != nil {
		return nil, fmt.Errorf("failed to read chapter-metadata.json: %w", err)
	}

	meta := &types.ChapterMetadata{}
	if err := json.Unmarshal(meta_data, meta); err != nil {
		return nil, fmt.Errorf("failed to parse chapter-metadata.json: %w", err)
	}

	ch := &Chapter{
		Dir:      ch_dir,
		Number:   chapter_num,
		Metadata: meta,
	}

	// Load ketabs
	for _, item := range meta.GetKetabs() {
		ketab_path := filepath.Join(ch_dir, item.File)
		body, err := os.ReadFile(ketab_path)
		if err != nil {
			return nil, fmt.Errorf("failed to read ketab file %s: %w", item.File, err)
		}

		// Strip scene headers (e.g., "# Scene 1\n")
		cleaned := strip_scene_header(string(body))

		ch.Ketabs = append(ch.Ketabs, Ketab{
			Item: item,
			Body: cleaned,
		})
	}

	return ch, nil
}

// strip_scene_header removes the scene header line from ketab content.
var scene_header_re = regexp.MustCompile(`^#\s*Scene\s*\d+[^\n]*\n+`)

func strip_scene_header(body string) string {
	body = scene_header_re.ReplaceAllString(body, "")
	return strings.TrimSpace(body)
}

// GetChapterNumbers returns all chapter numbers in order.
func (b *Book) GetChapterNumbers() []string {
	var nums []string
	for num := range b.Chapters {
		nums = append(nums, num)
	}
	sort.Strings(nums)
	return nums
}

// GetChapter returns a specific chapter by number.
func (b *Book) GetChapter(num string) (*Chapter, bool) {
	ch, ok := b.Chapters[num]
	return ch, ok
}

// CompileChapterBody combines all ketabs into a single chapter body.
func (c *Chapter) CompileChapterBody() string {
	var bodies []string
	for _, k := range c.Ketabs {
		bodies = append(bodies, k.Body)
	}
	return strings.Join(bodies, "\n\n---\n\n")
}

// Validate checks if a book directory has all required files.
func Validate(dir string) []string {
	var errors []string
	abs_dir, err := filepath.Abs(dir)
	if err != nil {
		return []string{fmt.Sprintf("invalid directory: %v", err)}
	}

	// Check book-metadata.json
	meta_path := filepath.Join(abs_dir, "book-metadata.json")
	if _, err := os.Stat(meta_path); os.IsNotExist(err) {
		errors = append(errors, "missing book-metadata.json")
	}

	// Check book-shape.json (optional but recommended)
	shape_path := filepath.Join(abs_dir, "book-shape.json")
	if _, err := os.Stat(shape_path); os.IsNotExist(err) {
		errors = append(errors, "missing book-shape.json (optional)")
	}

	// Load metadata to check chapters
	meta_data, err := os.ReadFile(meta_path)
	if err != nil {
		return errors
	}

	meta := &types.BookMetadata{}
	if err := json.Unmarshal(meta_data, meta); err != nil {
		errors = append(errors, fmt.Sprintf("invalid book-metadata.json: %v", err))
		return errors
	}

	// Check required fields
	if meta.BookTitle == "" {
		errors = append(errors, "book-metadata.json: missing book_title")
	}
	if meta.BookSlug == "" {
		errors = append(errors, "book-metadata.json: missing book_slug")
	}
	if meta.Author == "" {
		errors = append(errors, "book-metadata.json: missing author")
	}
	if meta.BookUUID == "" {
		errors = append(errors, "book-metadata.json: missing book_uuid")
	}

	// Check each chapter from acts
	for _, ch_ref := range meta.GetAllChapters() {
		ch_dir := filepath.Join(abs_dir, ch_ref.ChapterNumber)
		ch_meta_path := filepath.Join(ch_dir, "chapter-metadata.json")

		if _, err := os.Stat(ch_dir); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("missing chapter directory: %s", ch_ref.ChapterNumber))
			continue
		}

		if _, err := os.Stat(ch_meta_path); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("chapter %s: missing chapter-metadata.json", ch_ref.ChapterNumber))
			continue
		}

		// Load chapter metadata
		ch_meta_data, err := os.ReadFile(ch_meta_path)
		if err != nil {
			errors = append(errors, fmt.Sprintf("chapter %s: cannot read chapter-metadata.json", ch_ref.ChapterNumber))
			continue
		}

		ch_meta := &types.ChapterMetadata{}
		if err := json.Unmarshal(ch_meta_data, ch_meta); err != nil {
			errors = append(errors, fmt.Sprintf("chapter %s: invalid chapter-metadata.json: %v", ch_ref.ChapterNumber, err))
			continue
		}

		if ch_meta.ChapterUUID == "" {
			errors = append(errors, fmt.Sprintf("chapter %s: missing chapter_uuid", ch_ref.ChapterNumber))
		}

		// Check ketab files
		for _, item := range ch_meta.GetKetabs() {
			ketab_path := filepath.Join(ch_dir, item.File)
			if _, err := os.Stat(ketab_path); os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("chapter %s: missing ketab file %s", ch_ref.ChapterNumber, item.File))
			}
			if item.UUID == "" {
				errors = append(errors, fmt.Sprintf("chapter %s: ketab %s missing UUID", ch_ref.ChapterNumber, item.File))
			}
		}
	}

	return errors
}

// Status returns a summary of what exists on disk.
type BookStatus struct {
	BookTitle     string
	BookSlug      string
	BookUUID      string
	Author        string
	ChapterCount  int
	TotalKetabs   int
	HasShape      bool
	Chapters      []ChapterStatus
}

// ChapterStatus is status for a single chapter.
type ChapterStatus struct {
	Number       string
	Title        string
	UUID         string
	KetabCount   int
	HasMetadata  bool
	MissingFiles []string
}

// GetStatus returns the status of a book directory.
func GetStatus(dir string) (*BookStatus, error) {
	abs_dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("invalid directory: %w", err)
	}

	// Load metadata
	meta_path := filepath.Join(abs_dir, "book-metadata.json")
	meta_data, err := os.ReadFile(meta_path)
	if err != nil {
		return nil, fmt.Errorf("cannot read book-metadata.json: %w", err)
	}

	meta := &types.BookMetadata{}
	if err := json.Unmarshal(meta_data, meta); err != nil {
		return nil, fmt.Errorf("invalid book-metadata.json: %w", err)
	}

	all_chapters := meta.GetAllChapters()
	status := &BookStatus{
		BookTitle:    meta.BookTitle,
		BookSlug:     meta.BookSlug,
		BookUUID:     meta.BookUUID,
		Author:       meta.Author,
		ChapterCount: len(all_chapters),
	}

	// Check book-shape.json
	shape_path := filepath.Join(abs_dir, "book-shape.json")
	if _, err := os.Stat(shape_path); err == nil {
		status.HasShape = true
	}

	// Check chapters from acts
	for _, ch_ref := range all_chapters {
		ch_status := ChapterStatus{
			Number: ch_ref.ChapterNumber,
			Title:  ch_ref.ChapterTitle,
			UUID:   ch_ref.ChapterUUID,
		}

		ch_dir := filepath.Join(abs_dir, ch_ref.ChapterNumber)
		ch_meta_path := filepath.Join(ch_dir, "chapter-metadata.json")

		if ch_meta_data, err := os.ReadFile(ch_meta_path); err == nil {
			ch_status.HasMetadata = true
			ch_meta := &types.ChapterMetadata{}
			if err := json.Unmarshal(ch_meta_data, ch_meta); err == nil {
				items := ch_meta.GetKetabs()
				ch_status.KetabCount = len(items)
				status.TotalKetabs += len(items)

				// Check for missing ketab files
				for _, item := range items {
					ketab_path := filepath.Join(ch_dir, item.File)
					if _, err := os.Stat(ketab_path); os.IsNotExist(err) {
						ch_status.MissingFiles = append(ch_status.MissingFiles, item.File)
					}
				}
			}
		}

		status.Chapters = append(status.Chapters, ch_status)
	}

	return status, nil
}
