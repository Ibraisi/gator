set shell := ["nu", "-c"]
set dotenv-load := true

root := justfile_directory()
db_url := "postgres://postgres:" + env_var("DB_PASSWORD") + "@localhost:5432/gator"
migrations := "./sql/schema"

dbup:
    goose -dir {{migrations}} postgres "{{db_url}}" up

dbdown:
    goose -dir {{migrations}} postgres "{{db_url}}" down

dbfulldown:
    goose -dir {{migrations}} postgres "{{db_url}}" down-to 0

gen:
    sqlc generate

run *args:
    go build -o {{root}}/gator {{root}}/cmd; if $env.LAST_EXIT_CODE == 0 { {{root}}/gator {{args}} }

test:
    go test ./...
