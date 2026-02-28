package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joinnextblock/ketab-protocol/cli/internal/book"
	"github.com/joinnextblock/ketab-protocol/cli/internal/events"
	"github.com/joinnextblock/ketab-protocol/cli/internal/types"
	core "github.com/joinnextblock/ketab-protocol/go-core"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
)

var (
	flag_nsec          string
	flag_chapters      string
	flag_dry_run       bool
	flag_relays        string
	flag_ketabs_only   bool
	flag_clean_metadata bool
	// add-to-library flags
	flag_library_id    string
	flag_notes         string
	flag_rating        int
	flag_status        string
	flag_tags          string
)

func main() {
	root := &cobra.Command{
		Use:   "ketab",
		Short: "Ketab Protocol CLI — publish books to Nostr",
	}

	// publish
	publish_cmd := &cobra.Command{
		Use:   "publish <book-dir>",
		Short: "Publish book events (ketabs, chapters, book, library) to relays",
		Long:  "Publish book events to relays. Use --ketabs-only to skip chapters (30023) and only publish ketabs (38893), book (38891), and library (38890).",
		Args:  cobra.ExactArgs(1),
		RunE:  run_publish,
	}
	publish_cmd.Flags().StringVar(&flag_nsec, "nsec", "", "Author nsec (or set KETAB_NSEC env)")
	publish_cmd.Flags().StringVar(&flag_chapters, "chapters", "", "Comma-separated chapter numbers (default: all)")
	publish_cmd.Flags().BoolVar(&flag_dry_run, "dry-run", false, "Generate events without publishing")
	publish_cmd.Flags().StringVar(&flag_relays, "relays", strings.Join(types.DefaultRelays, ","), "Comma-separated relay URLs")
	publish_cmd.Flags().BoolVar(&flag_ketabs_only, "ketabs-only", false, "Publish only ketabs (38893), skip chapters (30023)")

	// validate
	validate_cmd := &cobra.Command{
		Use:   "validate <book-dir>",
		Short: "Validate book directory structure",
		Args:  cobra.ExactArgs(1),
		RunE:  run_validate,
	}

	// status
	status_cmd := &cobra.Command{
		Use:   "status <book-dir>",
		Short: "Show book status",
		Args:  cobra.ExactArgs(1),
		RunE:  run_status,
	}

	// delete-threads
	delete_threads_cmd := &cobra.Command{
		Use:   "delete-threads <book-dir>",
		Short: "Delete discussion threads referenced in book metadata",
		Long:  "Send NIP-09 deletion requests for discussion threads. Use --clean-metadata to also remove thread IDs from book metadata and republish.",
		Args:  cobra.ExactArgs(1),
		RunE:  run_delete_threads,
	}
	delete_threads_cmd.Flags().StringVar(&flag_nsec, "nsec", "", "Author nsec (or set KETAB_NSEC env)")
	delete_threads_cmd.Flags().StringVar(&flag_chapters, "chapters", "", "Comma-separated chapter numbers (default: all)")
	delete_threads_cmd.Flags().BoolVar(&flag_dry_run, "dry-run", false, "Show what would be deleted without sending deletion events")
	delete_threads_cmd.Flags().StringVar(&flag_relays, "relays", strings.Join(types.DefaultRelays, ","), "Comma-separated relay URLs")
	delete_threads_cmd.Flags().BoolVar(&flag_clean_metadata, "clean-metadata", false, "Remove discussion_id fields from metadata and republish book")

	// add-to-library
	add_to_library_cmd := &cobra.Command{
		Use:   "add-to-library <book-naddr>",
		Short: "Add a book to your library (publish Library Entry event kind 38892)",
		Long:  "Creates a Library Entry event that adds the specified book to the librarian's library collection.",
		Args:  cobra.ExactArgs(1),
		RunE:  run_add_to_library,
	}
	add_to_library_cmd.Flags().StringVar(&flag_nsec, "nsec", "", "Librarian's nsec (or set KETAB_NSEC env)")
	add_to_library_cmd.Flags().StringVar(&flag_library_id, "library-id", "a5213b36-5ad4-41c0-93d4-06b2adddcea8", "Library UUID")
	add_to_library_cmd.Flags().StringVar(&flag_notes, "notes", "", "Personal notes about the book")
	add_to_library_cmd.Flags().IntVar(&flag_rating, "rating", 0, "Rating 1-5 (optional)")
	add_to_library_cmd.Flags().StringVar(&flag_status, "status", "reading", "Status: want-to-read, reading, finished, abandoned")
	add_to_library_cmd.Flags().StringVar(&flag_tags, "tags", "", "Comma-separated tags")
	add_to_library_cmd.Flags().StringVar(&flag_relays, "relays", strings.Join(types.DefaultRelays, ","), "Comma-separated relay URLs")
	add_to_library_cmd.Flags().BoolVar(&flag_dry_run, "dry-run", false, "Generate event without publishing")

	root.AddCommand(publish_cmd, validate_cmd, status_cmd, delete_threads_cmd, add_to_library_cmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func resolve_nsec() (string, error) {
	if flag_nsec != "" {
		return flag_nsec, nil
	}
	if v := os.Getenv("KETAB_NSEC"); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("no nsec found — use --nsec flag or KETAB_NSEC env")
}

func decode_nsec(nsec_str string) (string, string, error) {
	prefix, data, err := nip19.Decode(nsec_str)
	if err != nil {
		return "", "", fmt.Errorf("invalid nsec: %w", err)
	}
	if prefix != "nsec" {
		return "", "", fmt.Errorf("expected nsec, got %s", prefix)
	}
	// go-nostr nip19.Decode returns hex string for nsec
	var sk string
	switch v := data.(type) {
	case string:
		sk = v
	case []byte:
		sk = hex.EncodeToString(v)
	default:
		return "", "", fmt.Errorf("unexpected nsec data type: %T", data)
	}
	pk, err := nostr.GetPublicKey(sk)
	if err != nil {
		return "", "", fmt.Errorf("failed to derive pubkey: %w", err)
	}
	return sk, pk, nil
}

func publish_event(ctx context.Context, event *nostr.Event, relays []string) error {
	for _, url := range relays {
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			fmt.Printf("  ⚠️  %s: %v\n", url, err)
			continue
		}
		err = relay.Publish(ctx, *event)
		relay.Close()
		if err != nil {
			fmt.Printf("  ⚠️  %s: %v\n", url, err)
		} else {
			fmt.Printf("  ✅ %s\n", url)
		}
	}
	return nil
}

func run_publish(cmd *cobra.Command, args []string) error {
	book_dir := args[0]
	relays := strings.Split(flag_relays, ",")

	nsec_str, err := resolve_nsec()
	if err != nil {
		return err
	}

	sk, pk, err := decode_nsec(nsec_str)
	if err != nil {
		return err
	}

	fmt.Printf("📐 Pubkey: %s\n", pk)
	fmt.Printf("📖 Loading book from %s\n\n", book_dir)

	bk, err := book.Load(book_dir)
	if err != nil {
		return fmt.Errorf("failed to load book: %w", err)
	}

	// Determine chapters
	var chapter_nums []string
	if flag_chapters != "" {
		for _, c := range strings.Split(flag_chapters, ",") {
			chapter_nums = append(chapter_nums, strings.TrimSpace(c))
		}
	} else {
		chapter_nums = bk.GetChapterNumbers()
	}

	fmt.Printf("Chapters: %s\n", strings.Join(chapter_nums, ", "))
	fmt.Printf("Dry run: %v\n\n", flag_dry_run)

	builder := events.NewBuilder(pk, relays[0])
	ctx := context.Background()

	var success, total int

	// 1. Ketabs
	fmt.Println("═══ KETABS ═══")
	for _, ch_num := range chapter_nums {
		ch, ok := bk.GetChapter(ch_num)
		if !ok {
			fmt.Printf("  ⚠️  Chapter %s not found on disk, skipping\n", ch_num)
			continue
		}
		for _, ketab := range ch.Ketabs {
			event := builder.BuildKetab(ch, ketab)
			event.PubKey = pk
			if err := events.SignEvent(&event, sk); err != nil {
				fmt.Printf("  ❌ Sign failed: %v\n", err)
				continue
			}
			total++
			fmt.Printf("\n📤 Ketab: ch%s #%d \"%s\" (id: %s)\n", ch_num, ketab.Item.Number, ketab.Item.Title, event.ID[:12])
			if !flag_dry_run {
				publish_event(ctx, &event, relays)
			}
			success++
		}
	}

	// 2. Chapters (skip if --ketabs-only)
	if !flag_ketabs_only {
		fmt.Println("\n═══ CHAPTERS ═══")
		for _, ch_num := range chapter_nums {
			ch, ok := bk.GetChapter(ch_num)
			if !ok {
				continue
			}
			event := builder.BuildChapter(bk, ch)
			event.PubKey = pk
			if err := events.SignEvent(&event, sk); err != nil {
				fmt.Printf("  ❌ Sign failed: %v\n", err)
				continue
			}
			total++
			fmt.Printf("\n📤 Chapter %s: \"%s\" (id: %s)\n", ch_num, ch.Metadata.ChapterTitle, event.ID[:12])
			if !flag_dry_run {
				publish_event(ctx, &event, relays)
			}
			success++
		}
	} else {
		fmt.Println("\n═══ CHAPTERS ═══")
		fmt.Println("  🚫 Skipped (--ketabs-only mode)")
	}

	// 3. Book
	fmt.Println("\n═══ BOOK ═══")
	book_event := builder.BuildBook(bk, chapter_nums)
	book_event.PubKey = pk
	if err := events.SignEvent(&book_event, sk); err != nil {
		return fmt.Errorf("sign book failed: %w", err)
	}
	total++
	fmt.Printf("\n📤 Book: \"%s\" (id: %s)\n", bk.Metadata.BookTitle, book_event.ID[:12])
	if !flag_dry_run {
		publish_event(ctx, &book_event, relays)
	}
	success++

	// 4. Library
	fmt.Println("\n═══ LIBRARY ═══")
	library_id := "a5213b36-5ad4-41c0-93d4-06b2adddcea8"
	library_event := builder.BuildLibrary(bk, library_id, "the library")
	library_event.PubKey = pk
	if err := events.SignEvent(&library_event, sk); err != nil {
		return fmt.Errorf("sign library failed: %w", err)
	}
	total++
	fmt.Printf("\n📤 Library (id: %s)\n", library_event.ID[:12])
	if !flag_dry_run {
		publish_event(ctx, &library_event, relays)
	}
	success++

	// Summary
	fmt.Printf("\n🏁 Done: %d/%d events published\n", success, total)

	// Print naddr
	naddr, err := nip19.EncodeEntity(pk, 38891, bk.Metadata.BookUUID, relays)
	if err == nil {
		fmt.Printf("\n📚 naddr: %s\n", naddr)
		fmt.Printf("   /book/%s\n", naddr)
	}

	return nil
}

func run_validate(cmd *cobra.Command, args []string) error {
	errors := book.Validate(args[0])
	if len(errors) == 0 {
		fmt.Println("✅ Book directory is valid")
		return nil
	}
	fmt.Printf("❌ %d issues found:\n", len(errors))
	for _, e := range errors {
		fmt.Printf("  - %s\n", e)
	}
	return fmt.Errorf("%d validation errors", len(errors))
}

func run_status(cmd *cobra.Command, args []string) error {
	status, err := book.GetStatus(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("📖 %s\n", status.BookTitle)
	fmt.Printf("   Slug: %s\n", status.BookSlug)
	fmt.Printf("   UUID: %s\n", status.BookUUID)
	fmt.Printf("   Author: %s\n", status.Author)
	fmt.Printf("   Chapters: %d\n", status.ChapterCount)
	fmt.Printf("   Total ketabs: %d\n", status.TotalKetabs)
	fmt.Printf("   Has shape: %v\n\n", status.HasShape)

	for _, ch := range status.Chapters {
		marker := "✅"
		if !ch.HasMetadata || len(ch.MissingFiles) > 0 {
			marker = "⚠️"
		}
		fmt.Printf("  %s %s: %s (%d ketabs)\n", marker, ch.Number, ch.Title, ch.KetabCount)
		for _, f := range ch.MissingFiles {
			fmt.Printf("       missing: %s\n", f)
		}
	}

	return nil
}

func run_delete_threads(cmd *cobra.Command, args []string) error {
	book_dir := args[0]
	relays := strings.Split(flag_relays, ",")

	nsec_str, err := resolve_nsec()
	if err != nil {
		return err
	}

	sk, pk, err := decode_nsec(nsec_str)
	if err != nil {
		return err
	}

	fmt.Printf("📐 Pubkey: %s\n", pk)
	fmt.Printf("📖 Loading book from %s\n\n", book_dir)

	// Load book metadata
	book_file := filepath.Join(book_dir, "book.json")
	book_data, err := os.ReadFile(book_file)
	if err != nil {
		return fmt.Errorf("failed to read book.json: %w", err)
	}

	var book_metadata map[string]interface{}
	if err := json.Unmarshal(book_data, &book_metadata); err != nil {
		return fmt.Errorf("failed to parse book.json: %w", err)
	}

	// Extract discussion IDs
	discussion_ids := make(map[string]string) // chapter_number -> discussion_id
	
	// Try to extract from "shape" field first (new format)
	if shape, ok := book_metadata["shape"].([]interface{}); ok {
		for _, act_interface := range shape {
			if act, ok := act_interface.([]interface{}); ok {
				for _, chapter_interface := range act {
					if chapter, ok := chapter_interface.(map[string]interface{}); ok {
						if discussion_id, ok := chapter["discussion_id"].(string); ok && discussion_id != "" {
							// Need to derive chapter number from d_tag or title
							if title, ok := chapter["title"].(string); ok {
								// This is a simplified mapping - in a real implementation you'd need better logic
								chapter_num := derive_chapter_number_from_title(title)
								if chapter_num != "" {
									discussion_ids[chapter_num] = discussion_id
								}
							}
						}
					}
				}
			}
		}
	}

	// Also try "acts" field (legacy format)
	if acts, ok := book_metadata["acts"].([]interface{}); ok {
		for _, act_interface := range acts {
			if act, ok := act_interface.(map[string]interface{}); ok {
				if chapters, ok := act["chapters"].([]interface{}); ok {
					for _, chapter_interface := range chapters {
						if chapter, ok := chapter_interface.(map[string]interface{}); ok {
							if discussion_id, ok := chapter["discussion_id"].(string); ok && discussion_id != "" {
								if number, ok := chapter["number"].(string); ok {
									discussion_ids[number] = discussion_id
								}
							}
						}
					}
				}
			}
		}
	}

	if len(discussion_ids) == 0 {
		fmt.Println("No discussion threads found in book metadata")
		return nil
	}

	// Filter by chapters if specified
	var target_chapters []string
	if flag_chapters != "" {
		target_chapters = strings.Split(flag_chapters, ",")
		for i, ch := range target_chapters {
			target_chapters[i] = strings.TrimSpace(ch)
		}
	} else {
		for ch := range discussion_ids {
			target_chapters = append(target_chapters, ch)
		}
	}

	fmt.Printf("Found discussion threads for chapters: %v\n", target_chapters)
	fmt.Printf("Dry run: %v\n\n", flag_dry_run)

	// Delete discussion threads
	fmt.Println("═══ DELETING THREADS ═══")
	ctx := context.Background()
	var deleted_count int

	for _, ch_num := range target_chapters {
		discussion_id, exists := discussion_ids[ch_num]
		if !exists {
			fmt.Printf("⚠️  Chapter %s: no discussion thread found\n", ch_num)
			continue
		}

		fmt.Printf("\n🗑️  Chapter %s: %s\n", ch_num, discussion_id[:12]+"...")
		
		if flag_dry_run {
			fmt.Printf("  [DRY RUN] Would send deletion event\n")
			deleted_count++
			continue
		}

		// Create NIP-09 deletion event
		deletion_event := nostr.Event{
			Kind:      5, // NIP-09 Event Deletion
			CreatedAt: nostr.Now(),
			Tags: nostr.Tags{
				{"e", discussion_id},
			},
			Content: "Deleting discussion thread - superseded by improved version",
		}
		deletion_event.PubKey = pk

		if err := events.SignEvent(&deletion_event, sk); err != nil {
			fmt.Printf("  ❌ Sign failed: %v\n", err)
			continue
		}

		// Publish to relays
		for _, url := range relays {
			relay, err := nostr.RelayConnect(ctx, url)
			if err != nil {
				fmt.Printf("  ⚠️  %s: %v\n", url, err)
				continue
			}
			err = relay.Publish(ctx, deletion_event)
			relay.Close()
			if err != nil {
				fmt.Printf("  ⚠️  %s: %v\n", url, err)
			} else {
				fmt.Printf("  ✅ %s\n", url)
			}
		}
		deleted_count++
	}

	fmt.Printf("\n🏁 Deletion requests sent for %d threads\n", deleted_count)

	// Clean metadata if requested
	if flag_clean_metadata {
		fmt.Println("\n═══ CLEANING METADATA ═══")
		
		// Remove discussion_id fields from the specified chapters
		cleaned := false
		
		// Clean "shape" field
		if shape, ok := book_metadata["shape"].([]interface{}); ok {
			for _, act_interface := range shape {
				if act, ok := act_interface.([]interface{}); ok {
					for _, chapter_interface := range act {
						if chapter, ok := chapter_interface.(map[string]interface{}); ok {
							if title, ok := chapter["title"].(string); ok {
								chapter_num := derive_chapter_number_from_title(title)
								for _, target_ch := range target_chapters {
									if chapter_num == target_ch {
										delete(chapter, "discussion_id")
										cleaned = true
									}
								}
							}
						}
					}
				}
			}
		}

		// Clean "acts" field
		if acts, ok := book_metadata["acts"].([]interface{}); ok {
			for _, act_interface := range acts {
				if act, ok := act_interface.(map[string]interface{}); ok {
					if chapters, ok := act["chapters"].([]interface{}); ok {
						for _, chapter_interface := range chapters {
							if chapter, ok := chapter_interface.(map[string]interface{}); ok {
								if number, ok := chapter["number"].(string); ok {
									for _, target_ch := range target_chapters {
										if number == target_ch {
											delete(chapter, "discussion_id")
											cleaned = true
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if cleaned {
			if flag_dry_run {
				fmt.Println("✅ [DRY RUN] Would clean metadata, removing discussion_id fields")
				fmt.Println("✅ [DRY RUN] Would republish book with cleaned metadata")
			} else {
				// Write updated metadata
				updated_data, err := json.MarshalIndent(book_metadata, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal updated metadata: %w", err)
				}

				if err := os.WriteFile(book_file, updated_data, 0644); err != nil {
					return fmt.Errorf("failed to write updated book.json: %w", err)
				}

				fmt.Println("✅ Metadata cleaned, discussion_id fields removed")

				// Republish book
				fmt.Println("\n═══ REPUBLISHING BOOK ═══")
				
				bk, err := book.Load(book_dir)
				if err != nil {
					return fmt.Errorf("failed to reload book: %w", err)
				}

				builder := events.NewBuilder(pk, relays[0])
				book_event := builder.BuildBook(bk, bk.GetChapterNumbers())
				book_event.PubKey = pk
				if err := events.SignEvent(&book_event, sk); err != nil {
					return fmt.Errorf("sign book failed: %w", err)
				}

				fmt.Printf("\n📤 Book: \"%s\" (id: %s)\n", bk.Metadata.BookTitle, book_event.ID[:12])
				publish_event(ctx, &book_event, relays)

				fmt.Println("✅ Book republished with cleaned metadata")
			}
		}
	}

	fmt.Println("\nNote: Event deletion is not guaranteed - some relays may ignore deletion requests")

	return nil
}

// Helper function to derive chapter number from title
// This is a simplified implementation - you might want to make this more robust
func derive_chapter_number_from_title(title string) string {
	chapter_map := map[string]string{
		"21 Sats":                     "01",
		"Shittier Twitter":            "02", 
		"Permission Granted":          "03",
		"Can't Steal What's Free":     "04",
		"Obviously, I Checked":        "04", // Alternative title
		"Congrats on Fifteen Years":   "05",
		"GM ☕":                       "06",
		"My Muse":                     "07",
		"Tick Tock":                   "08",
	}
	
	if num, ok := chapter_map[title]; ok {
		return num
	}
	return ""
}

func run_add_to_library(cmd *cobra.Command, args []string) error {
	book_naddr := args[0]
	relays := strings.Split(flag_relays, ",")

	// Resolve nsec
	nsec_str, err := resolve_nsec()
	if err != nil {
		return err
	}

	sk, pk, err := decode_nsec(nsec_str)
	if err != nil {
		return err
	}

	fmt.Printf("📐 Librarian pubkey: %s\n", pk)
	fmt.Printf("📚 Adding book to library: %s\n", flag_library_id)
	fmt.Printf("📖 Book naddr: %s\n\n", book_naddr)

	// Parse book naddr
	prefix, data, err := nip19.Decode(book_naddr)
	if err != nil {
		return fmt.Errorf("invalid naddr: %w", err)
	}
	if prefix != "naddr" {
		return fmt.Errorf("expected naddr, got %s", prefix)
	}

	entity_pointer, ok := data.(nostr.EntityPointer)
	if !ok {
		return fmt.Errorf("invalid naddr data")
	}

	// Validate this is a book event
	if entity_pointer.Kind != core.KindBook {
		return fmt.Errorf("expected book event (kind %d), got kind %d", core.KindBook, entity_pointer.Kind)
	}

	book_author_pubkey := entity_pointer.PublicKey
	book_d_tag := entity_pointer.Identifier

	fmt.Printf("Book author: %s\n", book_author_pubkey)
	fmt.Printf("Book d-tag: %s\n", book_d_tag)
	fmt.Printf("Status: %s\n", flag_status)
	fmt.Printf("Dry run: %v\n\n", flag_dry_run)

	// Build Library Entry content
	var rating_ptr *int
	if flag_rating > 0 {
		rating_ptr = &flag_rating
	}

	var tags_array []string
	if flag_tags != "" {
		for _, tag := range strings.Split(flag_tags, ",") {
			tags_array = append(tags_array, strings.TrimSpace(tag))
		}
	}

	book_coordinate := fmt.Sprintf("%d:%s:%s", core.KindBook, book_author_pubkey, book_d_tag)
	library_coordinate := fmt.Sprintf("%d:%s:%s", core.KindLibrary, pk, flag_library_id)
	d_tag := fmt.Sprintf("%s:%s", flag_library_id, book_coordinate)

	content := core.LibraryEntryContent{
		Notes:                 flag_notes,
		Rating:                rating_ptr,
		Tags:                  tags_array,
		ReadStatus:            flag_status,
		AddedAt:               nostr.Now().Time().Unix(),
		RefLibraryOwnerPubkey: pk,
		RefLibraryID:          flag_library_id,
		RefBookCoordinate:     book_coordinate,
		RefBookPubkey:         book_author_pubkey,
		RefBookID:             book_d_tag,
	}

	content_json, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	// Build Library Entry event (kind 38892)
	library_entry_event := nostr.Event{
		Kind:      core.KindLibraryEntry,
		CreatedAt: nostr.Now(),
		Tags: nostr.Tags{
			{"d", d_tag},
			{"a", book_coordinate},
			{"a", library_coordinate},
		},
		Content: string(content_json),
		PubKey:  pk,
	}

	// Sign the event
	if err := events.SignEvent(&library_entry_event, sk); err != nil {
		return fmt.Errorf("sign library entry failed: %w", err)
	}

	fmt.Println("═══ LIBRARY ENTRY ═══")
	fmt.Printf("\n📤 Library Entry: book \"%s\" → library %s (id: %s)\n", 
		book_d_tag[:12]+"...", flag_library_id[:8]+"...", library_entry_event.ID[:12])

	// Publish event
	if !flag_dry_run {
		ctx := context.Background()
		publish_event(ctx, &library_entry_event, relays)
	} else {
		fmt.Println("  [DRY RUN] Event would be published")
	}

	fmt.Printf("\n🏁 Library Entry event created successfully\n")

	// Show event details
	if flag_dry_run {
		fmt.Printf("\n📋 Event JSON:\n%s\n", string(content_json))
	}

	return nil
}
