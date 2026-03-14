# Gator

CLI RSS feed aggregator built in Go. Boot.dev backend track — [course link](https://www.boot.dev/courses/build-blog-aggregator-golang).

## Requirements

- [Go 1.22+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

```
go install github.com/ibrais/gator/cmd@latest
```

This installs the `gator` binary to your `$GOPATH/bin`. Make sure that directory is on your `$PATH`.

## Database setup

Create the database:

```sql
CREATE DATABASE gator;
```

Run migrations using [goose](https://github.com/pressly/goose):

```
goose -dir ./sql/schema postgres "postgres://<user>:<password>@localhost:5432/gator" up
```

## Config

Create `~/.gatorconfig.json`:

```json
{ "db_url": "postgres://<user>:<password>@localhost:5432/gator" }
```

Replace `<user>` and `<password>` with your PostgreSQL credentials.

## Usage

```
gator <command> [args]
```

Quick start:

```
gator register alice
gator addfeed "Hacker News" https://news.ycombinator.com/rss
gator agg 1m
```

See [docs/usage.md](docs/usage.md) for the full command reference.

## Development

Requires [just](https://github.com/casey/just) and [sqlc](https://sqlc.dev) for development tasks.

`just run <command>` builds and runs the binary in one step — use this during development instead of installing.

| Recipe | Description |
|---|---|
| `just run <args>` | Build and run gator |
| `just dbup` | Apply all pending migrations |
| `just dbdown` | Roll back one migration |
| `just gen` | Regenerate sqlc Go code |
| `just test` | Run all tests |

## Future Ideas

- [ ] `digest` command — filter stored posts by topic and summarize with AI (e.g. daily Go digest, weekly news summary)
- [ ] Daily/weekly digest scheduler
- [ ] `browse` command to page through stored posts
- [ ] Read/unread tracking per user
- [ ] Filter posts by feed or keyword
- [ ] Export digest to email or file
