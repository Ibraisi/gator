# Usage

All commands are run via:
```
just run <command> [args]
```

## Commands

| Command | Args | Description |
|---|---|---|
| `register` | `<name>` | Create a new user and set as current |
| `login` | `<name>` | Switch to an existing user |
| `users` | | List all users, marks current with `(current)` |
| `addfeed` | `<name> <url>` | Add an RSS feed tied to the current user |
| `feeds` | | List all feeds with their creator |
| `agg` | | Fetch and print a feed |
| `reset` | | Delete all users |

## Examples

```
just run register alice
just run login alice
just run addfeed "Wagslane" https://www.wagslane.dev/index.xml
just run feeds
just run agg
```

## just recipes

| Recipe | Description |
|---|---|
| `just dbup` | Apply all pending migrations |
| `just dbdown` | Roll back one migration |
| `just dbfulldown` | Roll back all migrations |
| `just gen` | Regenerate sqlc Go code |
| `just run <args>` | Build and run gator |
| `just test` | Run all tests |
