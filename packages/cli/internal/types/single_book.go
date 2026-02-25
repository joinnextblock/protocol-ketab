// Package types defines data structures for single-file book format.
package types

// SingleBookFile represents the new unified book.json file structure.
type SingleBookFile struct {
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	UUID        string    `json:"uuid"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Summary     string    `json:"summary,omitempty"`
	Image       string    `json:"image,omitempty"`
	Thumb       string    `json:"thumb,omitempty"`
	RefBlockID  string    `json:"ref_block_id,omitempty"`
	Acts        []SingleAct `json:"acts"`
}

// SingleAct represents an act in the single file format.
type SingleAct struct {
	Title    string           `json:"title"`
	Chapters []SingleChapter  `json:"chapters"`
}

// SingleChapter represents a chapter with embedded ketabs.
type SingleChapter struct {
	Number string        `json:"number"`
	Title  string        `json:"title"`
	UUID   string        `json:"uuid"`
	Ketabs []SingleKetab `json:"ketabs"`
}

// SingleKetab represents a ketab with file reference.
type SingleKetab struct {
	Title string `json:"title"`
	UUID  string `json:"uuid"`
	File  string `json:"file"`
}

// ToBookMetadata converts SingleBookFile to the legacy BookMetadata format.
func (s *SingleBookFile) ToBookMetadata() *BookMetadata {
	var acts []ActRef
	
	for _, act := range s.Acts {
		var chapters []ChapterRef
		for _, ch := range act.Chapters {
			chapters = append(chapters, ChapterRef{
				ChapterNumber: ch.Number,
				ChapterTitle:  ch.Title,
				ChapterUUID:   ch.UUID,
			})
		}
		acts = append(acts, ActRef{
			Title:    act.Title,
			Chapters: chapters,
		})
	}

	return &BookMetadata{
		BookTitle:   s.Title,
		BookSlug:    s.Slug,
		BookUUID:    s.UUID,
		Author:      s.Author,
		Summary:     s.Summary,
		Description: s.Description,
		Image:       s.Image,
		Thumb:       s.Thumb,
		Acts:        acts,
	}
}

// ToBookShape converts SingleBookFile to the legacy BookShape format.
func (s *SingleBookFile) ToBookShape() *BookShape {
	var shape [][]ShapeChapter
	
	for _, act := range s.Acts {
		var actChapters []ShapeChapter
		for _, ch := range act.Chapters {
			var ketabs []ShapeKetab
			for _, ketab := range ch.Ketabs {
				ketabs = append(ketabs, ShapeKetab{
					Title: ketab.Title,
					DTag:  ketab.UUID,
				})
			}
			actChapters = append(actChapters, ShapeChapter{
				Title:  ch.Title,
				DTag:   ch.UUID,
				Ketabs: ketabs,
			})
		}
		shape = append(shape, actChapters)
	}

	return &BookShape{
		Title:       s.Title,
		Description: s.Description,
		Image:       s.Image,
		Shape:       shape,
	}
}

// GetChapterMetadata returns ChapterMetadata for a specific chapter number.
func (s *SingleBookFile) GetChapterMetadata(chapterNum string) *ChapterMetadata {
	for _, act := range s.Acts {
		for _, ch := range act.Chapters {
			if ch.Number == chapterNum {
				var scenes []SceneRef
				for i, ketab := range ch.Ketabs {
					scenes = append(scenes, SceneRef{
						SceneNumber: i + 1,
						SceneFile:   ketab.File,
						SceneTitle:  ketab.Title,
						KetabUUID:   ketab.UUID,
					})
				}
				
				return &ChapterMetadata{
					ChapterTitle:  ch.Title,
					ChapterNumber: ch.Number,
					ChapterUUID:   ch.UUID,
					Scenes:        scenes,
				}
			}
		}
	}
	return nil
}