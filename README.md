# Ketab Protocol

**Your stories. Your keys. Your sovereignty.**

The attention economy is broken. Platforms capture your content, monetize your audience, and own your reach. Writers publish into walled gardens where algorithms decide who sees their work.

Ketab Protocol fixes this. Every paragraph, scene, and insight becomes an independent, addressable asset on Nostr. Readers pay you directly in Bitcoin. No platform can delete your work or ban your audience.

## The 80/20 Revolution

AI writes the first draft. You provide the 20% that matters: judgment, taste, voice. Every insight you refine becomes a **ketab** — an atomic unit of content with its own identity, its own economics.

While others ship raw AI output, you curate. While platforms extract, you own. While algorithms decide reach, your readers choose directly.

This isn't about replacing human creativity. It's about amplifying it with tools that serve creators, not platforms.

## How It Works

**Write → Publish → Earn**

1. **Write your story** in scenes and chapters, like always
2. **Publish to Nostr** with one command — each scene becomes an addressable ketab
3. **Readers find and pay you directly** in Bitcoin, no platform taking a cut

Your content gets a permanent address (`naddr`) that no platform can delete. Readers can reference specific scenes, pay for individual insights, and you keep 100% of what they pay.

No algorithms deciding your reach. No terms of service changing your deal. Just you, your readers, and honest exchange.

## The Platform Problem

**Medium owns your audience.** Change their algorithm, lose your reach. Violate their terms, lose your archive.

**Substack owns your payments.** They decide the fee structure. They control the subscriber relationship.

**Twitter owns your thoughts.** 280 characters, their algorithm, their ads next to your content.

**Ketab gives you sovereignty.** Your keys, your content, your readers, your payments. Publish once, own forever. No platform can deactivate your account or change your deal.

When readers pay 100 sats for your scene about bitcoin custody, those sats go directly to your wallet. No platform fee. No payment processor. No monthly subscription to keep your content accessible.

## How Stories Live on Nostr

**Every scene becomes a ketab** — an atomic unit of content with its own permanent address. Readers can reference, zap, and discuss individual scenes without losing the connection to your larger work.

**Books organize your ketabs** into chapters and stories. The same ketabs can appear in multiple books, creating composable narratives.

**Libraries curate the best work** across authors. Like a bookstore, but owned by curators who earn Bitcoin when readers discover books through their collections.

**Three roles, one protocol:**
- **Authors** create and publish ketabs
- **Librarians** curate collections and earn from discovery
- **Readers** pay directly and own their engagement data

You can play all three roles with the same Nostr keypair. Write your own stories, curate others' work, build your library.

## Publishing Your Work

**One command publishes your entire book** to Nostr. The `ketab` CLI handles all the protocol complexity so you focus on writing.

```bash
# Check your book is ready
ketab validate ./my-story

# Publish to the network
ketab publish ./my-story

# Your readers now have permanent addresses for every scene
```

**Keep it simple:**
```
my-story/
├── book-metadata.json       # Title, description, cover image
├── 01/                      # Chapter 1
│   ├── scene-1.md          # Individual scenes
│   └── scene-2.md
├── 02/                      # Chapter 2
│   └── scene-1.md
```

Write in markdown. Organize in folders. Publish with one command. Your story becomes permanently addressable content that readers can engage with at the scene level.

## Getting Started

**1. Install the CLI**
```bash
# Download for your system or build from source
go install github.com/joinnextblock/protocol-ketab/packages/cli/ketab@latest
```

**2. Set up your first book**
```bash
mkdir my-first-book && cd my-first-book
```

**3. Write your story** in markdown files organized by chapters

**4. Publish to Nostr**
```bash
ketab publish . --nsec your_nostr_private_key
```

Your scenes are now permanently addressable. Readers can find them, pay for them, and reference them forever.

**Need your Nostr key?** Any Nostr client can generate one. Or use: `ketab keygen`

## Why This Matters

**For Writers:**
- Own your audience relationship, no platform intermediary
- Earn Bitcoin directly from readers, keep 100%
- Content lives forever, immune to deplatforming
- Granular monetization — scenes, not just whole books

**For Readers:** 
- Pay creators directly, no subscription overhead
- Reference and quote specific scenes permanently  
- Discover new authors through librarian curation
- Own your reading data, not trapped in a platform

**For Curators:**
- Build valuable libraries and earn from discovery
- Surface the best content across the network
- Become a trusted taste-maker with economic incentives

**The Network Effect:**
As more creators publish ketabs, the citation web grows stronger. Your best insights get referenced across books, driving new readers to your work. Quality content rises through economic signals, not algorithmic manipulation.

## Beyond Publishing

**Bitcoin-Native Timing** — Ketab integrates with City Protocol to timestamp your work to Bitcoin blocks. Your publication becomes part of the Bitcoin timeline, not just calendar time.

**Attention Marketplace** — Through ATTN Protocol, readers can pay for your attention directly. Write scenes, earn sats, respond to the market signal of what your audience values most.

**Sovereign Infrastructure** — Run your own Nostr relay. Control your distribution. No dependency on corporate infrastructure that can change the rules or shut you down.

This isn't just about publishing books. It's about building an economy where creators and readers interact directly, honestly, with no extractive intermediaries.

## The Ecosystem

**library.nextblock.city** — Read and discover ketab-based books with scene-level engagement and Bitcoin payments.

**ketab CLI** — Publish your books from the command line. No GUI, no platform signup, just your content going directly to the network.

**Your own tools** — The protocol is open. Build whatever reading or publishing experience you want. Readers can access your ketabs from any compatible client.

## The Compounding Effect

Every ketab you publish becomes a discovery engine for your other work. When someone quotes your scene about custody, new readers find your book about sovereignty. When a curator adds your work to their library, their audience discovers you.

Quality content gets cited. Citations drive discovery. Discovery generates revenue. Revenue incentivizes quality.

It's how markets work when creators own the infrastructure instead of renting it from platforms.

## The Stack

- **ketab CLI** (Go) — Publish your books to Nostr
- **@ketab/sdk** (TypeScript) — Build reading and publishing apps  
- **Protocol spec** — Open standard, no vendor lock-in

Everything is open source. Fork it, modify it, run your own version. No company controls the protocol.

## NextBlock Protocol Suite

Ketab works alongside other sovereignty protocols:

- **City Protocol** — Bitcoin block time as the universal clock
- **ATTN Protocol** — Direct creator payments in the attention economy
- **Dynasty Protocol** — Sovereign identity and reputation

Together, they form the infrastructure for an economy where creators own their relationships with readers.

## Start Publishing

```bash
# Install the CLI
go install github.com/joinnextblock/protocol-ketab/packages/cli/ketab@latest

# Publish your first book
ketab publish ./my-story --nsec your_private_key
```

Your work becomes permanently addressable content that readers can discover, engage with, and pay for directly.

No platform permission required. No terms of service to agree to. Just you, your readers, and the network.

---

**The attention economy is broken. We're building the alternative.**

Your stories. Your keys. Your sovereignty.