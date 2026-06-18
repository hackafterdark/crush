# Tree-Sitter Structural Search Implementation

## Summary

Added tree-sitter-based structural search as a new tool alongside `grep` for finding Go code by syntax structure. The agent uses `structural_search` as the "sniper rifle" (precise, syntax-aware) and falls back to `grep` as the "nuclear option" for unknown text patterns.

## What Was Implemented

### Phase 1: Dependencies &amp; CGO Configuration

- Added `github.com/tree-sitter/go-tree-sitter` and `github.com/tree-sitter/tree-sitter-go` to `go.mod`
- Updated `CGO_ENABLED: 0` → `1` in:
  - `Taskfile.yaml` (line 12)
  - `.goreleaser.yml` (line 46)
  - `.github/workflows/build.yml` (added `CGO_ENABLED: 1` env to build and test steps)

### Phase 2: Structural Search Tool

- `internal/agent/parser/parser.go` — tree-sitter parsing and querying API:
  - `Parse(code []byte) *sitter.Node` — parses Go source into AST
  - `Query(root *sitter.Node, querySExpr string) ([]Match, error)` — runs S-expression queries
  - Helper functions: `FindCaptures`, `DeduplicateByPosition`, `NodeChildren`, `NodeDescendants`
- `internal/agent/tools/structural_search.go` — tool implementation:
  - `StructuralSearchParams` — template_name, path, include, max_results
  - `NewStructuralSearchTool(workingDir string) fantasy.AgentTool`
  - Multi-file search with glob pattern support
  - AST caching for repeated file parsing
- `internal/agent/tools/structural_search.md.tpl` — tool description

### Phase 3: Query Template Registry

- `internal/agent/parser/templates.go` — 7 Go query templates:
  - `find_functions` — function and method declarations
  - `find_structs` — struct type definitions with fields
  - `find_variables` — var declarations with/without values
  - `find_interfaces` — interface types with method signatures
  - `find_calls` — function calls and method calls
  - `find_imports` — import declarations
  - `find_comments` — comment nodes

### Phase 4: Integration

- Registered in `internal/agent/coordinator.go` (line 698)
- Registered in `internal/agent/common_test.go` (line 177)
- Added `<tool_funnel>` section to `internal/agent/templates/coder.md.tpl` with protocol:
  1. Try `structural_search` first for syntax-aware code search
  2. Fallback to `grep` for complex/unknown text patterns
  3. Use LSP tools for cross-file symbol resolution

## How to Build (Windows)

### Prerequisites

1. **Go 1.26.4+** (via GOTOOLCHAIN=auto or manual)
2. **MSYS2 UCRT64 GCC** — required for CGO:
  - Install MinGW-w64: Run `F:\msys64\mingw64.exe` (or `F:\msys64\ucrt64.exe`)
  - Install gcc: `pacman -S mingw-w64-ucrt-x86_64-gcc`
  - GCC location: `F:\msys64\ucrt64\bin\gcc.exe`

### Build Commands

```
$env:CGO_ENABLED="1"; $env:GOTOOLCHAIN="auto"; $env:PATH="F:/msys64/ucrt64/bin;" + $env:Path; go build -o crush-sitter.exe .
```

### CI/CD (GitHub Actions)

The `.github/workflows/build.yml` workflow now includes `CGO_ENABLED: 1` in the build and test steps. The UCRT64 MinGW64 toolchain must be available on the runner (install via `mingw-w64` package on Windows runners).

## Key Design Decisions

1. **Go-only initially** — Only the Go grammar (`tree-sitter-go`) is bundled. Adding other languages requires adding their tree-sitter packages.
2. **Multi-file search** — Supports glob patterns (`*.go`, `internal/agent//*.go`) to search across directories, unlike single-file tools.
3. **Parallel with grep** — Not a replacement; the agent chooses based on the search pattern type.
4. **AST caching** — Parsed ASTs are cached per-file to avoid re-parsing the same files.
5. **CGO required** — Tree-sitter is a C library, so CGO must be enabled. This excludes Android 32-bit and other CGO-disabled targets.

## File Locations

```
internal/
  agent/
    parser/
      parser.go          # Tree-sitter parsing/querying API
      templates.go       # Query template registry (7 Go patterns)
    tools/
      structural_search.go      # Tool implementation
      structural_search.md.tpl  # Tool description template
    coordinator.go              # Tool registration (line 698)
    common_test.go              # Test setup registration (line 177)
templates/
  coder.md.tpl                  # Tool funnel protocol section
```

## Tree-Sitter API Notes

- `go-tree-sitter` v0.25.0 API used
- `Node.Utf8Text(source []byte)` — get node text from source
- `Query.CaptureNames()` — get all capture names by index
- `QueryCursor.Matches(query, node, text)` — execute query, returns `QueryMatches`
- `QueryMatches.Next()` — returns `*QueryMatch` (nil when exhausted)
- `QueryCapture.Node` — struct value (not pointer), use `&cap.Node` for methods

