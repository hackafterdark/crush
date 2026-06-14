Search source code using tree-sitter AST queries. This is the "sniper rifle" for finding code by syntax structure — use it before grep for finding functions, structs, variables, interfaces, calls, imports, or comments.

**TOOL FUNNEL PROTOCOL:**
1. **Try** `structural_search` **first** for finding functions, structs, variables, interfaces, calls, imports, or comments by syntax structure.
2. **Fallback to** `grep` **if** the pattern is too complex or the file is too large to parse.
3. **LSP tools** (references, diagnostics) **only** for cross-file symbol resolution or type information.

**Supported Languages:**
- **Go** (`*.go`) — `parameter_list`, `block`, `field_identifier`
- **C++** (`*.cpp`, `*.cc`, `*.cxx`, `*.hpp`, `*.hxx`) — `function_declarator`, `compound_statement`, `class_specifier`, `field_expression`
- **HCL** (`*.hcl`) — `block`, `attribute`, `identifier`, `body`
- **C++** (`*.cpp`, `*.cc`, `*.cxx`, `*.hpp`, `*.hxx`) — `function_definition`, `class_specifier`, `preproc_include`, `call_expression`
- **TypeScript** (`*.ts`, `*.tsx`) — `formal_parameters`, `statement_block`, `property_identifier`
- **JavaScript** (`*.js`, `*.jsx`) — `formal_parameters`, `statement_block`, `property_identifier`
- **Python** (`*.py`) — `parameters`, `block`, `call`
- **SQL** (`*.sql`) — `create_function`, `create_table`, `create_view`, `select`, `insert`, `update`, `delete`
- **Rust** (`*.rs`) — `function_item`, `struct_item`, `enum_item`, `let_declaration`, `use_declaration`
- **Java** (`*.java`) — `method_declaration`, `class_declaration`, `interface_declaration`, `variable_declaration`, `method_invocation`
- **PHP** (`*.php`) — `function_definition`, `class_declaration`, `interface_declaration`, `variable_name`, `function_call_expression`, `namespace_use_declaration`

**Notes:**
- Templates are language-specific. An S-expression for Go won't work for Python.
- If `language` is not specified, it is auto-detected from the `include` pattern or file extensions.
- Python comments are `extra: true` nodes (like Go) — they appear in the AST when queried.

Available templates: {{ range $index, $element := .AvailableTemplates }}{{ if $index }}, {{ end }}{{ $element }}{{ end }}
