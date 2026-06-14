package parser

import (
	"path/filepath"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	sitter_csharp "github.com/tree-sitter/tree-sitter-c-sharp/bindings/go"
	sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
	sitter_hcl "github.com/tree-sitter-grammars/tree-sitter-hcl/bindings/go"
	sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"
	sitter_js "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"
	sitter_py "github.com/tree-sitter/tree-sitter-python/bindings/go"
	sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
	sitter_sql "github.com/DerekStride/tree-sitter-sql/bindings/go"
	sitter_ts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
	tsgo "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

// Language represents a supported programming language.
type Language string

const (
	LanguageCSharp     Language = "csharp"
	LanguageCpp        Language = "cpp"
	LanguageHcl        Language = "hcl"
	LanguageGo         Language = "go"
	LanguageJava       Language = "java"
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
		"cpp",
		"hcl",
		"go",
		"java",
		"javascript",
		"python",
		"php",
		"rust",
		"sql",
		"typescript",
	}
}

// GetLanguage returns the tree-sitter language for the given name.
func GetLanguage(name string) *sitter.Language {
	switch name {
	case "go":
		return sitter.NewLanguage(tsgo.Language())
	case "cpp":
		return sitter.NewLanguage(sitter_cpp.Language())
	case "hcl":
		return sitter.NewLanguage(sitter_hcl.Language())
	case "java":
		return sitter.NewLanguage(sitter_java.Language())
	case "typescript":
		return sitter.NewLanguage(sitter_ts.Language_Typescript())
	case "javascript":
		return sitter.NewLanguage(sitter_js.Language())
	case "python":
		return sitter.NewLanguage(sitter_py.Language())
	case "php":
		return sitter.NewLanguage(sitter_php.Language())
	case "sql":
		return sitter.NewLanguage(sitter_sql.Language())
	case "rust":
		return sitter.NewLanguage(sitter_rust.Language())
	default:
		return sitter.NewLanguage(tsgo.Language())
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
	case ".hcl":
		return "hcl"
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
	case ".java":
		return "java"
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
				StartByte: cap.Node.StartByte(),
				EndByte:   cap.Node.EndByte(),
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
