Search Go source code using tree-sitter AST queries. This is the "sniper rifle" for finding code by syntax structure — use it before grep for finding functions, structs, variables, interfaces, calls, imports, or comments.

**TOOL FUNNEL PROTOCOL:**
1. **Try** `structural_search` **first** for finding functions, structs, variables, interfaces, calls, imports, or comments by syntax structure.
2. **Fallback to** `grep` **if** the pattern is too complex or the file is too large to parse.
3. **LSP tools** (references, diagnostics) **only** for cross-file symbol resolution or type information.

Available templates: {{ range $index, $element := .AvailableTemplates }}{{ if $index }}, {{ end }}{{ $element }}{{ end }}