# sptzxgo

`sptzxgo` is a modular WhatsApp bot written in Go and powered by `whatsmeow`. The project is designed around an event-driven architecture, category-based command modules, and a growing internal framework for normalization, dispatching, and policy enforcement.

## Highlights

- `whatsmeow`-based WhatsApp Web integration
- Modular command registry with category-driven organization
- Interactive login flow for QR or pairing code
- SQLite-backed session and runtime persistence
- Per-user economy: daily limit, extra limit, balance, premium, XP
- Category-based quota enforcement for AI, downloader, tools, maker, sticker, primbon, random, and stalk commands
- Game system with clue progression, surrender flow, and reward balancing
- Event normalization, dispatching, edit tracking, anti-delete, and poll state tracking
- Group-level feature toggles such as welcome, goodbye, anti-spam, and anti-delete

## Requirements

- Go 1.21+
- SQLite 3
- A valid `.env` file

## Quick Start

```bash
git clone https://github.com/siputzx/sptzxgo.git
cd sptzxgo
go mod tidy
cp .env.example .env
go run main.go
```

## Configuration

The bot is configured through environment variables.

| Variable | Description | Required |
|---|---|---|
| `BOT_OWNERS` | Comma-separated owner phone numbers without spaces | Yes |
| `BOT_PREFIX` | Comma-separated command prefixes | No |
| `SESSION_DB` | SQLite DSN for WhatsApp session and runtime storage | No |
| `LOGIN_METHOD` | `qr` or `paircode` | No |
| `PAIRING_PHONE` | Pairing phone number when using pairing mode | No |
| `LOG_LEVEL` | Bot log level | No |
| `TIMEZONE` | Runtime timezone used by utility features | No |
| `STICKER_PACK_NAME` | Sticker pack name metadata | No |
| `STICKER_AUTHOR` | Sticker author metadata | No |
| `ANTISPAM_MAX_PER_SECOND` | Anti-spam threshold per second | No |
| `ANTISPAM_MAX_PER_MINUTE` | Anti-spam threshold per minute | No |
| `ANTISPAM_BAN_DURATION_SECS` | Temporary anti-spam ban duration | No |
| `SIPUTZX_ENABLED` | Enable external API-backed features | No |
| `SIPUTZX_BASE_URL` | Base URL for the SiputzX API | No |
| `GEMINI_COOKIE` | Cookie used by Gemini-related functionality | No |

## Login Flow

When there is no stored session, the bot starts an interactive CLI login setup.

- `QR` mode displays the login QR directly in the terminal.
- `Pairing Code` mode prompts for an international phone number and validates the format.

If a valid WhatsApp session already exists in the session database, the bot reconnects automatically without asking for login again.

## Runtime Architecture

### Core Layers

- `src/config` — environment loading and validation
- `src/core` — command system, context, event normalization, stores, and shared runtime primitives
- `src/handler` — event intake, dispatching, message processing, revoke/edit/reaction/poll handlers
- `src/middleware` — anti-spam and command limiter middleware
- `src/serialize` — WhatsApp message building, media sending, downloads, and protocol helpers
- `src/commands` — modular feature packages grouped by category

### Event Pipeline

`whatsmeow event -> normalizer -> dispatcher -> specialized handler`

The bot now separates event intake from event processing so that message, reaction, edit, revoke, poll, presence, receipt, and call handling can evolve independently.

## Limit and Economy System

### Daily Limit

- Free users receive `100` daily limit.
- Premium users receive `500` daily limit.
- Daily limit resets at `00:00 UTC`.

### Extra Limit

Purchased or owner-added limit is stored separately as extra limit and is not removed by the daily reset.

### Command-Level Policies

The bot supports two different policies per command:

- `Quota` — consumes user limit when the command is used
- `Limit` — rate-limits repeated usage inside a time window

Example:

```go
core.Use(&core.Command{
    Name:     "menu",
    Category: "general",
    Quota:    core.PerUserQuota(1),
    Limit:    core.PerUserLimit(30, time.Minute),
    Handler:  func(ptz *core.Ptz) error { ... },
})
```

### Economy Features

- `mylimit` — view total limit, daily limit, extra limit, balance, and premium status
- `balance` — view balance, XP, and owner top-up contact
- `buylimit <amount>` — buy extra limit using balance
- `leaderboard` — top XP and top balance users with masked numbers
- `addlimit`, `addsaldo`, `addprem` — owner commands for managing user economy

## Group Features

Group-level toggles are available through `.enable` and `.disable`, including:

- welcome
- goodbye
- antidelete
- announce
- locked
- approval
- restrict
- ephemeral

## Anti-Delete and Message History

The bot keeps an in-memory message snapshot store for supported message types. When anti-delete is enabled in a group, the bot can restore or resend deleted content for:

- text
- image
- video
- audio / voice note
- document
- sticker
- location / live location
- contact

## Supported Feature Categories

- general
- group
- downloader
- sticker
- maker
- tools
- primbon
- random
- stalk
- ai
- games
- owner
- search
- info

## Project Layout

```text
sptzxgo/
├── main.go
├── README.md
├── .env.example
├── go.mod
└── src/
    ├── api/
    ├── commands/
    ├── config/
    ├── core/
    ├── handler/
    ├── middleware/
    ├── serialize/
    └── utils/
```

## Development Notes

- Run `go test ./...` to validate the full module.
- Run `go build ./...` to ensure the project still compiles.
- Keep command behavior category-consistent when applying quota or economy rules.
- Prefer extending the normalized event layer before adding new event-driven features.

## License

This project is licensed under the Apache 2.0 License. See `LICENSE` for details.
