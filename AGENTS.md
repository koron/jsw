# jsw — Jekyll Serve Watcher

Single-main-package Go module (`github.com/koron/jsw`, Go 1.25). The entrypoint is
`jsw.go` at the root. Three internal packages under `internal/`.

## Commands

```sh
make build       # go build -gcflags '-e' ./...
make test        # go test ./...
make race        # go test -race ./...
make checkall    # go vet ./... && staticcheck ./... (runs both)
make cover       # test with coverage HTML report in tmp/
```

## Architecture

| Directory | Role |
|-----------|------|
| `internal/jekyll/` | Spawns `jekyll serve` at start and `jekyll build` on file change |
| `internal/timebuf/` | Debounce timer (200 ms window) — coalesces rapid change events |
| `internal/watcher/` | Recursive `fswatcher` wrapper; excludes `_site/` and `.git/` |

## Testing

- `internal/watcher/watcher_test.go` is skipped (`t.Skip`).
- `internal/timebuf/buffer_test.go` is a basic integration-style test (no fixtures needed).
- Tests require no external services or credentials.
