# Repository Guidelines

This Go module converts iTunes Timed Text assets into TTML or WebVTT and ships both a library and CLI. The toolchain is verified with `go version go1.25.3`, `node v24.11.1`, and `ruby 3.2.3`; Go is the only required runtime for builds.

## Project Structure & Module Organization
- `cmd/ittconv` houses the CLI entrypoint; build outputs land in the repo root unless you pass `-o`.
- `internal/parser`, `internal/timecode`, `internal/ttml`, and `internal/vtt` isolate parsing, rational time handling, and format emitters. Keep new packages scoped under `internal` unless they are public API.
- `ittconv.go` exposes the high-level Go API, while `samples/` and `testdata/` contain representative .itt and golden subtitle fixtures.
- `docs/` holds contributor references such as `GUIDE.md`, checklists, and acceptance criteria used during reviews.

## Build, Test, and Development Commands
- `go build ./cmd/ittconv` – Compile the CLI; for a quick run use `go run ./cmd/ittconv --help`.
- `go test ./...` – Execute unit and property tests across every package; tack on `-cover` to inspect coverage.
- `go test ./... -run Integration -v` – Focus on the slower end-to-end cases guarded by the `Integration` suffix.
- `go vet ./...` – Static analysis pass; run it before opening a PR. Format new code with `gofmt -w <files>`.

## Coding Style & Naming Conventions
- Follow idiomatic Go: tabs for indentation, CamelCase exported symbols, and `lowercase` package names that mirror their directory.
- Keep filenames aligned with the role they play (`converter.go`, `writer.go`). Tests live beside sources and end with `_test.go`.
- Default logging uses the structured helpers already wired in the CLI; prefer extending existing options instead of introducing new flags without need.

## Testing Guidelines
- Write table-driven tests and name them `Test<Function><Scenario>` for clarity (`TestTimecodeFromFrameRate`).
- Golden TTML/VTT outputs belong in `testdata/` and should mirror the directory hierarchy of the code under test.
- Aim for full coverage of parsing edge cases (multi-byte glyphs, overlapping cues). Property tests live alongside the units; keep seeds deterministic when practical.

## Commit & Pull Request Guidelines
- Recent commits (`re-packaging`, `Cleanup`, `More tests`) show the expected style: concise, imperative, under ~60 characters, no prefixes unless linking to an issue (use `ABC-123: describe fix` when needed).
- PRs should describe the problem, the approach, and validation commands (`go test ./...`). Attach diffs or sample subtitle outputs when behavior changes, and reference any relevant checklist items from `docs/`.
