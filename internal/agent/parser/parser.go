package parser

import (
	"path/filepath"
	"strconv"
	"strings"

	lang_bash "github.com/charmbracelet/crush/internal/agent/parser/bash"
	lang_cpp "github.com/charmbracelet/crush/internal/agent/parser/cpp"
	lang_c "github.com/charmbracelet/crush/internal/agent/parser/c"
	lang_css "github.com/charmbracelet/crush/internal/agent/parser/css"
	lang_csharp "github.com/charmbracelet/crush/internal/agent/parser/csharp"
	lang_go "github.com/charmbracelet/crush/internal/agent/parser/go"
	lang_hcl "github.com/charmbracelet/crush/internal/agent/parser/hcl"
	lang_html "github.com/charmbracelet/crush/internal/agent/parser/html"
	// Java support requires an external scanner (scanner.c) that is not
	// present in the vendored grammar. See:
	//   https://github.com/tree-sitter/tree-sitter-java
	// lang_java "github.com/charmbracelet/crush/internal/agent/parser/java"
	lang_javascript "github.com/charmbracelet/crush/internal/agent/parser/javascript"
	lang_json "github.com/charmbracelet/crush/internal/agent/parser/json"
	lang_php "github.com/charmbracelet/crush/internal/agent/parser/php"
	lang_python "github.com/charmbracelet/crush/internal/agent/parser/python"
	lang_ruby "github.com/charmbracelet/crush/internal/agent/parser/ruby"
	lang_rust "github.com/charmbracelet/crush/internal/agent/parser/rust"
	lang_scala "github.com/charmbracelet/crush/internal/agent/parser/scala"
	lang_toml "github.com/charmbracelet/crush/internal/agent/parser/toml"
	lang_typescript "github.com/charmbracelet/crush/internal/agent/parser/typescript"

	"golang.org/x/exp/slices"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// Match represents a single query match result.
type Match struct {
	Index    int
	Captures []QueryResult
}

// QueryResult represents a captured node in a query match.
type QueryResult struct {
	Capture   string
	Text      string
	StartByte uint32
	EndByte   uint32
	StartPos  Pos
	EndPos    Pos
}

// Pos represents a position in source code.
type Pos struct {
	Row    uint
	Column uint
}

// Language represents a supported programming language.
type Language string

const (
	LanguageCSharp     Language = "csharp"
	LanguageBash       Language = "bash"
	LanguageCpp        Language = "cpp"
	LanguageC          Language = "c"
	LanguageHcl        Language = "hcl"
	LanguageGo         Language = "go"
	LanguageJSON       Language = "json"
	LanguageHtml       Language = "html"
	LanguageCss        Language = "css"
	LanguageToml       Language = "toml"
	LanguageScala      Language = "scala"
	LanguageRuby       Language = "ruby"
	// LanguageJava     Language = "java" — Java not supported (requires external scanner)
	LanguageJavaScript Language = "javascript"
	LanguagePython     Language = "python"
	LanguagePHP        Language = "php"
	LanguageRust       Language = "rust"
	LanguageSQL        Language = "sql"
	LanguageTypeScript Language = "typescript"
)

// SupportedLanguages returns the list of supported language names.
func SupportedLanguages() []string {
	return []string{
		"csharp",
		"c",
		"bash",
		"cpp",
		"hcl",
		"go",
		// "java" — requires external scanner not present in vendored grammar
		"javascript",
		"json",
		"html",
		"css",
		"toml",
		"scala",
		"python",
		"php",
		"rust",
		"sql",
		"ruby",
		"typescript",
	}
}

// GetLanguage returns the tree-sitter language for the given name.
func GetLanguage(name string) *sitter.Language {
	switch name {
	case "go":
		return lang_go.GetLanguage()
	case "cpp":
		return lang_cpp.GetLanguage()
	case "c":
		return lang_c.GetLanguage()
	case "bash":
		return lang_bash.GetLanguage()
	case "hcl":
		return lang_hcl.GetLanguage()
	case "csharp":
		return lang_csharp.GetLanguage()
	// Java not supported — requires external scanner not present in vendored grammar
	case "typescript":
		return lang_typescript.GetLanguage()
	case "javascript":
		return lang_javascript.GetLanguage()
	case "python":
		return lang_python.GetLanguage()
	case "php":
		return lang_php.GetLanguage()
	case "sql":
		return lang_go.GetLanguage() // SQL not yet vendored — falls back to Go
	case "rust":
		return lang_rust.GetLanguage()
	case "ruby":
		return lang_ruby.GetLanguage()
	case "json":
		return lang_json.GetLanguage()
	case "html":
		return lang_html.GetLanguage()
	case "css":
		return lang_css.GetLanguage()
	case "toml":
		return lang_toml.GetLanguage()
	case "scala":
		return lang_scala.GetLanguage()
	default:
		return lang_go.GetLanguage()
	}
}

// DetectLanguage returns the language name based on file extension.
func DetectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "go"
	case ".cpp", ".cc", ".cxx", ".hpp", ".hxx":
		return "cpp"
	case ".c", ".h":
		return "c"
	case ".sh":
		return "bash"
	case ".hcl":
		return "hcl"
	case ".rb":
		return "ruby"
	case ".json":
		return "json"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".toml":
		return "toml"
	case ".scala", ".sbt":
		return "scala"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	case ".py":
		return "python"
	case ".php":
		return "php"
	case ".sql":
		return "sql"
	case ".rs":
		return "rust"
	// ".java" — Java not supported (requires external scanner)
	case ".cs":
		return "csharp"
	default:
		return "go"
	}
}

// Parse parses source code using tree-sitter and returns the AST root node.
func Parse(code []byte, lang string) *sitter.Node {
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(GetLanguage(lang))
	tree := parser.Parse(code, nil)
	return tree.RootNode()
}

func Query(root *sitter.Node, code []byte, querySExpr string) ([]Match, error) {
	language := root.Language()
	query, queryErr := sitter.NewQuery(language, querySExpr)
	if queryErr != nil {
		return nil, queryErr
	}
	defer query.Close()

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	// Use the correct API: cursor.Matches
	matches := cursor.Matches(query, root, code)

	var results []Match
	matchCount := 0
	for {
		// Try to get the next match directly
		match := matches.Next()
		if match == nil {
			break
		}

		if len(match.Captures) == 0 {
			continue
		}
		var captures []QueryResult
		for _, cap := range match.Captures {
			captureName := ""
			if int(cap.Index) < len(query.CaptureNames()) {
				captureName = query.CaptureNames()[cap.Index]
			}
			captures = append(captures, QueryResult{
				Capture:   captureName,
				Text:      nodeToString(&cap.Node, code),
				StartByte: uint32(cap.Node.StartByte()),
				EndByte:   uint32(cap.Node.EndByte()),
				StartPos: Pos{
					Row:    cap.Node.StartPosition().Row,
					Column: cap.Node.StartPosition().Column,
				},
				EndPos: Pos{
					Row:    cap.Node.EndPosition().Row,
					Column: cap.Node.EndPosition().Column,
				},
			})
		}
		results = append(results, Match{
			Index:    matchCount,
			Captures: captures,
		})
		matchCount++
	}

	return results, nil
}

// nodeToString converts a tree-sitter node to its string representation.
func nodeToString(node *sitter.Node, source []byte) string {
	return node.Utf8Text(source)
}

// FindCaptures finds all captures matching a given name in the results.
func FindCaptures(matches []Match, captureName string) []QueryResult {
	var results []QueryResult
	for _, m := range matches {
		for _, c := range m.Captures {
			if c.Capture == captureName {
				results = append(results, c)
			}
		}
	}
	return results
}

// DeduplicateByPosition removes duplicate results that share the same start position and capture.
func DeduplicateByPosition(results []QueryResult) []QueryResult {
	seen := make(map[string]bool)
	var deduped []QueryResult
	for _, r := range results {
		key := r.Capture + ":" + strconv.Itoa(int(r.StartPos.Row)) + ":" + strconv.Itoa(int(r.StartPos.Column))
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, r)
		}
	}
	return deduped
}

// NodeToPos converts a sitter.Point to our Pos type.
func NodeToPos(p sitter.Point) Pos {
	return Pos{Row: p.Row, Column: p.Column}
}

// NodeChildren returns the named children of a node.
func NodeChildren(node *sitter.Node) []*sitter.Node {
	var children []*sitter.Node
	for i := uint(0); i < node.NamedChildCount(); i++ {
		child := node.NamedChild(i)
		if child != nil {
			children = append(children, child)
		}
	}
	return children
}

// NodeDescendants returns all descendants of a node (depth-first).
func NodeDescendants(node *sitter.Node) []*sitter.Node {
	var result []*sitter.Node
	var visit func(*sitter.Node)
	visit = func(n *sitter.Node) {
		result = append(result, n)
		for i := uint(0); i < n.NamedChildCount(); i++ {
			child := n.NamedChild(i)
			if child != nil {
				visit(child)
			}
		}
	}
	visit(node)
	return result
}

// Reverse reverses a slice in place.
func Reverse[T any](s []T) {
	slices.Reverse(s)
}
