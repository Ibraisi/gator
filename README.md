# Gator

CLI RSS feed aggregator built in Go. Boot.dev backend track — [course link](https://www.boot.dev/courses/build-blog-aggregator-golang).

## Requirements

- Go 1.22+
- PostgreSQL
- [goose](https://github.com/pressly/goose)
- [sqlc](https://sqlc.dev)

## Setup

1. Copy and configure the env file:
   ```
   cp .env.example .env
   ```
   Set `DB_PASSWORD` to your PostgreSQL password.

2. Run migrations (using [just](https://github.com/casey/just), optional but recommended):
   ```
   just dbup
   ```
   Or directly with goose:
   ```
   goose -dir ./sql/schema postgres "postgres://postgres:<password>@localhost:5432/gator" up
   ```

3. Create `~/.gatorconfig.json`:
   ```json
   { "db_url": "postgres://postgres:<password>@localhost:5432/gator" }
   ```

4. Build and run:
   ```
   just run <command>
   ```
   Or without just:
   ```
   go build -o gator ./cmd && ./gator <command>
   ```

See [docs/usage.md](docs/usage.md) for all available commands.
