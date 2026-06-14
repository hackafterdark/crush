package tools

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"charm.land/fantasy"
	"github.com/charmbracelet/crush/internal/agent/parser"
	"github.com/charmbracelet/crush/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

//go:embed structural_search.md.tpl
var structuralSearchDescriptionTmpl string // Change to string for easier template management

var structuralSearchDescriptionTpl = template.Must(
	template.New("structuralSearchDescription").
		Parse(structuralSearchDescriptionTmpl),
)

// Simplify this with the above and handle in structural_search.md.tpl instead
// var structuralSearchDescriptionTpl = template.Must(
// 	template.New("structuralSearchDescription").
// 		Funcs(template.FuncMap{
// 			"join": func(sep string, parts []string) string {
// 				if len(parts) == 0 {
// 					return ""
// 				}
// 				var sb strings.Builder
// 				sb.WriteString(parts[0])
// 				for _, p := range parts[1:] {
// 					sb.WriteString(sep)
// 					sb.WriteString(p)
// 				}
// 				return sb.String()
// 			},
// 		}).
// 		Parse(string(structuralSearchDescriptionTmpl)),
// )

type structuralSearchDescriptionData struct {
	AvailableTemplates []string
}

func structuralSearchDescription() string {
	return renderTemplate(structuralSearchDescriptionTpl, structuralSearchDescriptionData{
		AvailableTemplates: parser.TemplateNames(),
	})
}

// StructuralSearchParams are the parameters for the structural_search tool.
type StructuralSearchParams struct {
	// TemplateName is the name of the pre-built query template to use.
	TemplateName string `json:"template_name" description:"The name of the query template to use. Available: find_functions, find_structs, find_variables, find_interfaces, find_calls, find_imports, find_comments."`
	// Path is the directory to search in. Defaults to the current working directory.
	Path string `json:"path,omitempty" description:"The directory to search in. Defaults to the current working directory."`
	// Include is a file pattern to filter by (e.g., "*.go", "internal//*.go").
	Include string `json:"include,omitempty" description:"File pattern to include in the search (e.g., '*.go', 'internal//*.go'). Defaults to '*.go'."`
	// MaxResults is the maximum number of results to return.
	MaxResults int `json:"max_results,omitempty" description:"Maximum number of results to return (default: 100)."`
}

// StructuralSearchCapture represents a single capture within a match.
type StructuralSearchCapture struct {
	// Capture name (e.g., "name", "function_name", "field_name")
	Capture string `json:"capture"`
	// The matched text
	Text string `json:"text"`
	// Line number (1-indexed)
	Line int `json:"line"`
	// Column number (0-indexed)
	Column int `json:"column"`
}

// StructuralSearchMatch represents a complete match across files.
type StructuralSearchMatch struct {
	// File path where the match was found
	File string `json:"file"`
	// Match index within the file
	MatchIndex int `json:"match_index"`
	// All captures for this match
	Captures []StructuralSearchCapture `json:"captures"`
}

// structuralSearchResponse is the metadata returned with the response.
type structuralSearchResponse struct {
	Matches       []StructuralSearchMatch `json:"matches"`
	TotalMatches  int                     `json:"total_matches"`
	FilesSearched int                     `json:"files_searched"`
}

const (
	StructuralSearchToolName = "structural_search"
)

func findGoFiles(workingDir, path, include string) ([]string, error) {
	searchPath := path
	if searchPath == "" {
		searchPath = workingDir
	}

	var files []string
	if include != "" {
		// Use filepath.Glob for glob patterns
		globPattern := filepath.Join(searchPath, include)
		matches, err := filepath.Glob(globPattern)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			info, err := os.Stat(m)
			if err != nil {
				continue
			}
			if !info.IsDir() {
				files = append(files, m)
			}
		}
	} else {
		// Default: find all .go files
		err := filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() && strings.HasSuffix(path, ".go") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func formatResults(results []StructuralSearchMatch, maxResults int) string {
	if len(results) == 0 {
		return "No matches found"
	}

	var sb strings.Builder
	if len(results) >= maxResults {
		fmt.Fprintf(&sb, "Found at least %d matches (truncated)\n\n", maxResults)
	} else {
		fmt.Fprintf(&sb, "Found %d matches\n\n", len(results))
	}

	currentFile := ""
	for _, result := range results {
		if currentFile != result.File {
			if currentFile != "" {
				sb.WriteString("\n")
			}
			currentFile = result.File
			sb.WriteString(fmt.Sprintf("=== %s ===\n", result.File))
		}
		for i, cap := range result.Captures {
			if i == 0 {
				fmt.Fprintf(&sb, "  Line %d, Col %d: %s\n", cap.Line, cap.Column, cap.Text)
			} else {
				fmt.Fprintf(&sb, "    %-15s: %s\n", cap.Capture, cap.Text)
			}
		}
	}

	return sb.String()
}

func executeStructuralSearch(ctx context.Context, workingDir string, params StructuralSearchParams) (fantasy.ToolResponse, error) {
	searchPath := params.Path
	if searchPath == "" {
		searchPath = workingDir
	}

	files, err := findGoFiles(workingDir, searchPath, params.Include)
	if err != nil {
		return fantasy.NewTextErrorResponse("error finding files: " + err.Error()), nil
	}

	if len(files) == 0 {
		return fantasy.NewTextResponse("No Go files found matching the pattern"), nil
	}

	// Get the query template
	query, ok := parser.GetTemplate(params.TemplateName)
	if !ok {
		available := strings.Join(parser.TemplateNames(), ", ")
		return fantasy.NewTextErrorResponse("unknown template: " + params.TemplateName + ". Available: " + available), nil
	}

	maxResults := params.MaxResults
	if maxResults <= 0 {
		maxResults = 100
	}

	var allResults []StructuralSearchMatch
	filesSearched := 0

	for _, file := range files {
		// Read file
		code, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Parse with tree-sitter
		root := parser.Parse(code)

		// Run query
		matches, err := parser.Query(root, code, query)
		if err != nil {
			// Log query error but continue with other files
			continue
		}

		filesSearched++

		// Convert matches to response format
		for _, match := range matches {
			var captures []StructuralSearchCapture
			for _, cap := range match.Captures {
				captures = append(captures, StructuralSearchCapture{
					Capture: cap.Capture,
					Text:    cap.Text,
					Line:    int(cap.StartPos.Row) + 1, // Convert to 1-indexed
					Column:  int(cap.StartPos.Column),
				})
			}
			allResults = append(allResults, StructuralSearchMatch{
				File:       file,
				MatchIndex: match.Index,
				Captures:   captures,
			})

			if len(allResults) >= maxResults {
				goto done
			}
		}
	}

done:
	if len(allResults) == 0 {
		return fantasy.WithResponseMetadata(
			fantasy.NewTextResponse("No matches found"),
			structuralSearchResponse{
				Matches:       nil,
				TotalMatches:  0,
				FilesSearched: filesSearched,
			},
		), nil
	}

	return fantasy.WithResponseMetadata(
		fantasy.NewTextResponse(formatResults(allResults, maxResults)),
		structuralSearchResponse{
			Matches:       allResults,
			TotalMatches:  len(allResults),
			FilesSearched: filesSearched,
		},
	), nil
}

// NewStructuralSearchTool creates a new structural search tool.
func NewStructuralSearchTool(workingDir string) fantasy.AgentTool {
	return fantasy.NewAgentTool(
		"structural_search",
		structuralSearchDescription(),
		func(ctx context.Context, params StructuralSearchParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			ctx, span := otel.StartSpan(ctx, "execute_tool structural_search")
			defer span.End()
			span.SetAttributes(
				attribute.String("gen_ai.tool.name", StructuralSearchToolName),
				attribute.String("gen_ai.tool.call.id", call.ID),
				attribute.String("gen_ai.tool.call.arguments", call.Input),
			)
			return executeStructuralSearch(ctx, workingDir, params)
		},
	)
}
