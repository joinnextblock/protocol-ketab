package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/joinnextblock/ketab-protocol/cli/internal/book"
	"github.com/joinnextblock/ketab-protocol/cli/internal/events"
	"github.com/joinnextblock/ketab-protocol/cli/internal/types"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
)

var (
	flag_nsec     string
	flag_chapters string
	flag_dry_run  bool
	flag_relays   string
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
		Args:  cobra.ExactArgs(1),
		RunE:  run_publish,
	}
	publish_cmd.Flags().StringVar(&flag_nsec, "nsec", "", "Author nsec (or set ARCHITECT_NSEC env)")
	publish_cmd.Flags().StringVar(&flag_chapters, "chapters", "", "Comma-separated chapter numbers (default: all)")
	publish_cmd.Flags().BoolVar(&flag_dry_run, "dry-run", false, "Generate events without publishing")
	publish_cmd.Flags().StringVar(&flag_relays, "relays", strings.Join(types.DefaultRelays, ","), "Comma-separated relay URLs")

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

	root.AddCommand(publish_cmd, validate_cmd, status_cmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func resolve_nsec(book_dir string) (string, error) {
	// Flag first
	if flag_nsec != "" {
		return flag_nsec, nil
	}

	// Env
	if v := os.Getenv("ARCHITECT_NSEC"); v != "" {
		return v, nil
	}

	// .env in book dir
	env_path := filepath.Join(book_dir, ".env")
	if _, err := os.Stat(env_path); err == nil {
		godotenv.Load(env_path)
		if v := os.Getenv("ARCHITECT_NSEC"); v != "" {
			return v, nil
		}
		if v := os.Getenv("AUTHOR_NSEC"); v != "" {
			return v, nil
		}
	}

	return "", fmt.Errorf("no nsec found — use --nsec, ARCHITECT_NSEC env, or .env in book dir")
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

	nsec_str, err := resolve_nsec(book_dir)
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

	// 2. Chapters
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

	_ = time.Now() // keep import
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
