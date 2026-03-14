default:
    just --list

set dotenv-load := true
set quiet := true

root := justfile_directory()
migrations := "./sql/schema"

# apply all pending migrations
dbup:
    #!/usr/bin/env nu
    goose -dir {{migrations}} postgres $"postgres://postgres:($env.DB_PASSWORD)@localhost:5432/gator" up

# roll back one migration
dbdown:
    #!/usr/bin/env nu
    goose -dir {{migrations}} postgres $"postgres://postgres:($env.DB_PASSWORD)@localhost:5432/gator" down

# roll back all migrations
dbfulldown:
    #!/usr/bin/env nu
    goose -dir {{migrations}} postgres $"postgres://postgres:($env.DB_PASSWORD)@localhost:5432/gator" down-to 0

# regenerate sqlc Go code
gen:
    #!/usr/bin/env nu
    sqlc generate

# build and run gator
run *args:
    #!/usr/bin/env nu
    go build -o {{root}}/gator {{root}}/cmd; if $env.LAST_EXIT_CODE == 0 { {{root}}/gator {{args}} }

# run all tests
test:
    #!/usr/bin/env nu
    go test ./... -v

# build release binary
release-build:
    #!/usr/bin/env nu
    go build -o gator-linux-amd64 ./cmd
