# translator

A Telegram bot that translates text into Georgian using GPT. It's a personal
project I built to help a Russian- and English-speaking contact of mine
communicate in Georgian: send it a message and it replies with the Georgian
translation.

It also supports **Latin transliteration**: if you type Georgian phonetically
in Latin letters (e.g. `gamarjoba`), the bot detects that, asks you to confirm,
and converts it to the Mkhedruli script before translating.

## How it works

```
Telegram update
      │
      ▼
handler.TelegramHandler   (telebot.v3 handlers: /start, /setLang, plain text, inline buttons)
      │
      ├─ service.UserService              user registration & language preference
      ├─ service.LanguageService           supported languages
      ├─ service.MessageTemplateService    per-language UI strings
      ├─ service.ButtonService             inline keyboard definitions
      └─ service.TranslationService
                │
                ├─ service.OpenAiService   calls the OpenAI chat completions API
                └─ util.ToGeorgian         Latin → Mkhedruli transliteration
      │
      ▼
repository.*Postgres   (pgx/v5, one interface + Postgres implementation per aggregate)
      │
      ▼
PostgreSQL
```

Each repository is defined as an interface (`repository.UserRepository`, etc.)
with a single Postgres-backed implementation, so the service layer isn't
coupled to the database driver. Users, languages, message templates and
button/menu definitions are all rows in Postgres rather than hardcoded
strings, so adding a language or changing bot copy doesn't require a
deployment.

## Running it

Requires Go 1.26+, Docker (for Postgres), the [golang-migrate](https://github.com/golang-migrate/migrate)
CLI, and an OpenAI API key.

```bash
brew install golang-migrate   # or: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

cp .env.example .env   # fill in TELEGRAM_BOT_TOKEN, OPENAI_API_KEY, DATABASE_URL, etc.
docker-compose up -d   # starts Postgres
make migrate-up        # runs migrations
go run ./cmd/server
```

Environment variables (see `cmd/server/main.go` and
`internal/service/open_ai.go`):

| Variable                       | Purpose                                        |
|---------------------------------|-------------------------------------------------|
| `DATABASE_URL`                  | Postgres connection string                     |
| `TELEGRAM_BOT_TOKEN`            | Bot token from @BotFather                      |
| `OPENAI_API_KEY`                | OpenAI API key                                 |
| `GPT_MODEL_VERSION`             | Model name, e.g. `gpt-5.4-nano`                 |
| `TRANSLATION_TIMEOUT_SECONDS`   | Per-request timeout for the OpenAI call        |
| `POSTGRES_USER/PASSWORD/DB`     | Used by `docker-compose.yml` to init Postgres  |

See `.env.example` for the full list.

## Testing

```bash
go test ./...
```

Unit tests cover the pure logic: the Latin→Georgian transliteration map and
the translation-response resolution logic. The Telegram/Postgres/OpenAI
integration itself isn't covered by automated tests.

## Project layout

```
cmd/server/           entrypoint: config, DB pool, bot wiring
internal/handler/      Telegram update handlers
internal/service/      business logic (users, languages, templates, buttons, translation)
internal/repository/   Postgres-backed persistence, one interface + impl per aggregate
internal/model/        domain types and sentinel errors
internal/util/         env helpers, Latin→Georgian transliteration
migrations/            golang-migrate SQL migrations (schema + seed data)
```
