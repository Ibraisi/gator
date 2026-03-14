# Usage

After installing, run commands as:

```
gator <command> [args]
```

During development, `just run <command> [args]` is equivalent.

## Commands

### User management

| Command | Args | Description |
|---|---|---|
| `register` | `<name>` | Create a new user and set as current |
| `login` | `<name>` | Switch to an existing user |
| `users` | | List all users, marks current with `(current)` |
| `reset` | | Delete all users (dev/debug use) |

### Feed management

Commands marked with * require a logged-in user.

| Command | Args | Description |
|---|---|---|
| `addfeed` * | `<name> <url>` | Add a new RSS feed and auto-follow it |
| `feeds` | | List all feeds with their URL and creator |
| `follow` * | `<url>` | Follow an existing feed by URL |
| `following` * | | List feeds the current user follows |
| `unfollow` * | `<url>` | Unfollow a feed by URL |

### Aggregation

| Command | Args | Description |
|---|---|---|
| `agg` * | `<time_between_reqs>` | Continuously fetch all feeds on a timer |

`agg` is a long-running loop — run it in a separate terminal while using other commands elsewhere. It always fetches the least recently updated feed first.

The interval uses Go duration syntax:

| Unit | Suffix | Example |
|---|---|---|
| Milliseconds | `ms` | `500ms` |
| Seconds | `s` | `30s` |
| Minutes | `m` | `5m` |
| Hours | `h` | `1h` |

Units can be combined: `1m30s`, `2h30m`.

## Examples

```
gator register alice
gator login alice
gator addfeed "Hacker News" https://news.ycombinator.com/rss
gator addfeed "Boot.dev Blog" https://www.boot.dev/blog/index.xml
gator follow https://techcrunch.com/feed/
gator following
gator feeds
gator agg 1m
```
