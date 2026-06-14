package parser

import (
	"slices"
	"strconv"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tsgo "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

// QueryResult represents a single match from a tree-sitter query.
type QueryResult struct {
	// Capture name (e.g., "function", "name", "parameter")
	Capture string `json:"capture"`
	// The matched text content
	Text string `json:"text"`
	// Start byte offset in the source
	StartByte uint `json:"start_byte"`
	// End byte offset in the source
	EndByte uint `json:"end_byte"`
	// Start position (row, column) 0-indexed
	StartPos Pos `json:"start_position"`
	// End position (row, column) 0-indexed
	EndPos Pos `json:"end_position"`
}

// Pos represents a position in the source code.
type Pos struct {
	Row    uint `json:"row"`
	Column uint `json:"column"`
}

// Match represents a complete match with all its captures.
type Match struct {
	// Index of the match within the query results
	Index int `json:"index"`
	// Captures within this match
	Captures []QueryResult `json:"captures"`
}

// Parse parses Go source code using tree-sitter and returns the AST root node.
func Parse(code []byte) *sitter.Node {
	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(sitter.NewLanguage(tsgo.Language()))
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
