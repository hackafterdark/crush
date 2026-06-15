# Tree-Sitter Structural Search

Crush's `structural_search` tool enables AI agents to find code by **syntax structure** rather than plain text. It parses source files into ASTs using tree-sitter and runs S-expression queries to locate functions, structs, calls, imports, and more — across **19 programming languages**.

## How It Works

```
Agent calls structural_search
    │
    ▼
Language auto-detected from file extensions
    │
    ▼
Files discovered by glob pattern, filtered by language
    │
    ▼
S-expression template resolved (e.g., "find_functions")
    │
    ▼
Each file parsed with its language-specific tree-sitter grammar
    │
    ▼
Query executed against AST → captures with positions
    │
    ▼
Results formatted: file, line, column, matched text
```

### Language Detection

The `language` parameter is **optional**. If omitted, Crush auto-detects the language:

1. From the `include` pattern (e.g., `"*.ts"` → TypeScript)
2. From file extensions when walking directories
3. Falls back to Go if nothing matches

Each file is parsed with its own detected grammar, so mixed-language directories work correctly.

### Template Abstraction

Agents use **human-readable template names** — not raw S-expressions:


| Template Name     | What It Finds                         |
| ----------------- | ------------------------------------- |
| `find_functions`  | Function and method declarations      |
| `find_structs`    | Struct, class, and type definitions   |
| `find_variables`  | Variable declarations and assignments |
| `find_interfaces` | Interface and trait definitions       |
| `find_calls`      | Function and method call sites        |
| `find_imports`    | Import and use statements             |
| `find_comments`   | Comment nodes                         |


Each language has its own S-expression set for each template, since tree-sitter node types differ between languages.

## Supported Languages &amp; Grammars


| Language       | Extensions                            | Package                  | Templates Available                          |
| -------------- | ------------------------------------- | ------------------------ | -------------------------------------------- |
| **Go**         | `.go`                                 | `tree-sitter-go`         | All 7                                        |
| **TypeScript** | `.ts`, `.tsx`                         | `tree-sitter-typescript` | All 7                                        |
| **JavaScript** | `.js`, `.jsx`                         | `tree-sitter-javascript` | 6 (no interfaces)                            |
| **Python**     | `.py`                                 | `tree-sitter-python`     | All 7                                        |
| **SQL**        | `.sql`                                | `tree-sitter-sql`        | 5 (no variables, interfaces, calls)          |
| **Rust**       | `.rs`                                 | `tree-sitter-rust`       | All 7                                        |
| **Java**       | `.java`                               | `tree-sitter-java`       | All 7                                        |
| **C#**         | `.cs`                                 | `tree-sitter-c-sharp`    | All 7                                        |
| **PHP**        | `.php`                                | `tree-sitter-php`        | All 7                                        |
| **C++**        | `.cpp`, `.cc`, `.cxx`, `.hpp`, `.hxx` | `tree-sitter-cpp`        | All 7                                        |
| **C**          | `.c`, `.h`                            | `tree-sitter-c`          | 5 (no interfaces, imports)                   |
| **Bash**       | `.sh`                                 | `tree-sitter-bash`       | 4 (no structs, interfaces, calls)            |
| **HCL**        | `.hcl`                                | `tree-sitter-hcl`        | 4 (no functions, interfaces, calls, imports) |
| **Ruby**       | `.rb`                                 | `tree-sitter-ruby`       | 5 (no interfaces, imports)                   |
| **JSON**       | `.json`                               | `tree-sitter-json`       | 3 (no functions, interfaces, calls, imports) |
| **HTML**       | `.html`, `.htm`                       | `tree-sitter-html`       | 3 (no functions, interfaces, calls)          |
| **CSS**        | `.css`                                | `tree-sitter-css`        | 3 (no functions, interfaces, calls)          |
| **TOML**       | `.toml`                               | `tree-sitter-toml`       | 3 (no functions, interfaces, calls, imports) |
| **Scala**      | `.scala`, `.sbt`                      | `tree-sitter-scala`      | All 7                                        |


### Template Availability by Language

```
Go:            ████████████████████  7/7 templates
TypeScript:    ████████████████████  7/7
JavaScript:    ████████████████████  6/7 (no interfaces)
Python:        ████████████████████  7/7
Rust:          ████████████████████  7/7
Java:          ████████████████████  7/7
C#:            ████████████████████  7/7
PHP:           ████████████████████  7/7
C++:           ████████████████████  7/7
C:             ██████████████████    5/7 (no interfaces, imports)
Bash:          ████████████          4/7 (no structs, interfaces, calls)
HCL:           ████████              4/7 (no functions, interfaces, calls, imports)
Ruby:          ██████████████        5/7 (no interfaces, imports)
JSON:          ██████                3/7 (no functions, interfaces, calls, imports)
HTML:          ██████                3/7 (no functions, interfaces, calls)
CSS:           ██████                3/7 (no functions, interfaces, calls)
TOML:          ██████                3/7 (no functions, interfaces, calls, imports)
Scala:         ████████████████████  7/7
SQL:           ████████████          5/7 (no variables, interfaces, calls)
```

## Tool Parameters


| Parameter       | Type    | Required | Description                                                       |
| --------------- | ------- | -------- | ----------------------------------------------------------------- |
| `template_name` | string  | Yes      | Name of the S-expression template to use                          |
| `path`          | string  | No       | Directory to search (default: current directory)                  |
| `include`       | string  | No       | Glob pattern to filter files (e.g., `"*.go"`, `"internal//*.go"`) |
| `max_results`   | integer | No       | Maximum number of results to return (default: 100)                |
| `language`      | string  | No       | Force a specific language; auto-detected if omitted               |


## Tool Funnel Protocol

Crush's agents follow this search priority:

1. `**structural_search` first** — Use for finding functions, structs, variables, interfaces, calls, imports, or comments by syntax structure. This is the "sniper rifle": precise, AST-aware, and language-aware.
2. `**grep` fallback** — Use only if the pattern is too complex for S-expressions, the file is too large to parse, or searching for plain text/regex patterns.
3. **LSP tools** — Use for "Find References", "Go to Definition", or cross-file symbol resolution with type information.

## CGO Requirement

Tree-sitter is a C library. All grammar bindings use CGO to embed the C parser sources, so **a C compiler is required** for building Crush.

### Build Requirements


| Platform    | Compiler         | Installation                                                         |
| ----------- | ---------------- | -------------------------------------------------------------------- |
| **Windows** | MSYS2 UCRT64 GCC | Install MinGW-w64: `pacman -S mingw-w64-ucrt-x86_64-gcc`             |
| **Linux**   | GCC or Clang     | `apt install gcc` (Debian/Ubuntu) or `yum install gcc` (RHEL/Fedora) |
| **macOS**   | Clang (Xcode)    | `xcode-select --install`                                             |


### Build Configuration

CGO is enabled in all build configurations:

- **Development**: `CGO_ENABLED=1` (default when a C compiler is available)
- **CI/CD**: `.github/workflows/build.yml` sets `CGO_ENABLED: 1`
- **Taskfile**: `taskfile.yaml` sets `CGO_ENABLED: 1`
- **GoReleaser**: `.goreleaser.yml` sets `CGO_ENABLED: 1`

### Build Example (Windows)

```powershell
$env:CGO_ENABLED = "1"
$env:GOTOOLCHAIN = "auto"
$env:PATH = "F:/msys64/ucrt64/bin;" + $env:Path
go build -o crush.exe .
```

```
$env:CGO_ENABLED="1"; $env:GOTOOLCHAIN="auto"; $env:PATH="F:/msys64/ucrt64/bin;" + $env:Path; go build -o crush-sitter.exe .
```

### Limitations

- **Android 32-bit**: CGO is not supported on this platform
- **Static builds**: CGO dependencies may require dynamic linking
- **Binary size**: CGO dependencies increase the final binary size

## Internal Architecture

```
internal/
  agent/
    parser/
      parser.go              # Language detection, AST parsing, querying
      templates.go           # Language-specific S-expression templates
    tools/
      structural_search.go   # Tool definition, params, execution
      structural_search.md.tpl  # Tool description for the agent
    coordinator.go           # Tool registration
    common_test.go           # Test setup registration
templates/
  coder.md.tpl               # Tool funnel protocol section
```

### Key Components

`**parser.go**` — Core parsing and querying:

- `Parse(code []byte, lang string) *sitter.Node` — parses source into AST
- `Query(root *sitter.Node, querySExpr string) ([]Match, error)` — runs S-expression queries
- `DetectLanguage(filePath string) string` — maps file extensions to language names
- `GetLanguage(lang string) *sitter.Language` — returns the grammar pointer for a language
- `SupportedLanguages() []string` — returns all 19 supported languages

`**templates.go**` — Template registry:

- `Templates map[string]map[string]string` — nested map: `language → templateName → SExpression`
- `GetTemplate(lang, name) (string, bool)` — resolves a template for a language
- `TemplateNames(lang string) []string` — lists available templates for a language

`**structural_search.go**` — Tool implementation:

- `StructuralSearchParams` — tool parameters
- `findFiles()` — discovers files by glob, filtered by language extensions
- `executeStructuralSearch()` — orchestrates parsing, querying, and formatting

## Adding New Languages

To add support for a new language:

1. **Install the tree-sitter package**:
  ```bash
   go get github.com/tree-sitter/tree-sitter-<lang>
  ```
2. **Update `parser.go`**:
  - Add a `Language` constant
  - Add a case to `GetLanguage()` returning the grammar pointer
  - Add an extension → language mapping in `DetectLanguage()`
3. **Update `templates.go`**:
  - Add a new top-level key in `Templates` for the language
  - Add S-expression templates (see existing templates for reference)
4. **Update `structural_search.go`**:
  - Add the language's extensions to `findFiles()`
5. **Update this documentation** — add the language to the supported languages table

## S-Expression Reference

### Go


| Template          | S-Expression                                                                          |
| ----------------- | ------------------------------------------------------------------------------------- |
| `find_functions`  | `(function_declaration name: (identifier) @name)`                                     |
| `find_structs`    | `(type_declaration (type_spec name: (type_identifier) @name type: (struct_type)))`    |
| `find_variables`  | `(var_declaration declaration: (var_declarator name: (identifier) @name))`            |
| `find_interfaces` | `(type_declaration (type_spec name: (type_identifier) @name type: (interface_type)))` |
| `find_calls`      | `(call_expression function: (identifier) @name)`                                      |
| `find_imports`    | `(import_declaration path: (interpreted_string_literal) @path)`                       |
| `find_comments`   | `(comment) @comment`                                                                  |


### TypeScript


| Template          | S-Expression                                                                                              |
| ----------------- | --------------------------------------------------------------------------------------------------------- |
| `find_functions`  | `(function_declaration name: (identifier) @name parameters: (formal_parameters) body: (statement_block))` |
| `find_structs`    | `(class_declaration name: (type_identifier) @name body: (class_body))`                                    |
| `find_variables`  | `(lexical_declaration (variable_declarator name: (identifier) @name))`                                    |
| `find_interfaces` | `(interface_declaration name: (type_identifier) @name)`                                                   |
| `find_calls`      | `(call_expression function: (identifier) @name arguments: (arguments))`                                   |
| `find_imports`    | `(import_statement source: (string) @path)`                                                               |
| `find_comments`   | `(comment) @comment`                                                                                      |


### Python


| Template          | S-Expression                                                                            |
| ----------------- | --------------------------------------------------------------------------------------- |
| `find_functions`  | `(function_definition name: (identifier) @name parameters: (parameters) body: (block))` |
| `find_structs`    | `(class_definition name: (identifier) @name body: (block))`                             |
| `find_variables`  | `(assignment left: (identifier) @name)`                                                 |
| `find_interfaces` | *not applicable*                                                                        |
| `find_calls`      | `(call function: (identifier) @name arguments: (arguments))`                            |
| `find_imports`    | `(import_statement name: (dotted_name) @path)`                                          |
| `find_comments`   | `(comment) @comment`                                                                    |


## Dependencies


| Package                                            | Version | Purpose                               |
| -------------------------------------------------- | ------- | ------------------------------------- |
| `github.com/tree-sitter/go-tree-sitter`            | v0.25.0 | Go bindings for tree-sitter C library |
| `github.com/tree-sitter/tree-sitter-go`            | v0.25.0 | Go grammar                            |
| `github.com/tree-sitter/tree-sitter-typescript`    | v0.23.2 | TypeScript grammar                    |
| `github.com/tree-sitter/tree-sitter-javascript`    | v0.25.0 | JavaScript grammar                    |
| `github.com/tree-sitter/tree-sitter-python`        | v0.25.0 | Python grammar                        |
| `github.com/DerekStride/tree-sitter-sql`           | v0.3.11 | SQL grammar                           |
| `github.com/tree-sitter/tree-sitter-rust`          | v0.24.2 | Rust grammar                          |
| `github.com/tree-sitter/tree-sitter-java`          | v0.23.5 | Java grammar                          |
| `github.com/tree-sitter/tree-sitter-c-sharp`       | v0.23.5 | C# grammar                            |
| `github.com/tree-sitter/tree-sitter-php`           | v0.24.2 | PHP grammar                           |
| `github.com/tree-sitter/tree-sitter-cpp`           | v0.23.4 | C++ grammar                           |
| `github.com/tree-sitter/tree-sitter-c`             | v0.24.2 | C grammar                             |
| `github.com/tree-sitter/tree-sitter-bash`          | v0.25.1 | Bash grammar                          |
| `github.com/tree-sitter-grammars/tree-sitter-hcl`  | v1.2.0  | HCL (Terraform) grammar               |
| `github.com/tree-sitter/tree-sitter-ruby`          | v0.23.1 | Ruby grammar                          |
| `github.com/tree-sitter/tree-sitter-json`          | v0.24.8 | JSON grammar                          |
| `github.com/tree-sitter/tree-sitter-html`          | v0.23.2 | HTML grammar                          |
| `github.com/tree-sitter/tree-sitter-css`           | v0.23.2 | CSS grammar                           |
| `github.com/tree-sitter-grammars/tree-sitter-toml` | v0.7.0  | TOML grammar                          |
| `github.com/tree-sitter/tree-sitter-scala`         | v0.23.2 | Scala grammar                         |


All 20 packages embed C source code and require CGO compilation.