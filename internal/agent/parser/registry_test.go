package parser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	lang_typescript "github.com/charmbracelet/crush/internal/agent/parser/typescript"
)

func TestRegistryDefaults(t *testing.T) {
	r := NewQueryRegistry()
	q, ok := r.GetTemplate("go", "find_functions")
	if !ok {
		t.Fatal("expected find_functions query to be present by default")
	}
	if q == "" {
		t.Fatal("expected non-empty default query")
	}
}

func TestRegistryOverrides(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "crush-query-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	queriesDir := filepath.Join(tempDir, ".crush", "queries")
	if err := os.MkdirAll(queriesDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// 1. Write single query overriding default.
	overrideYaml := `
id: find_functions
description: "Custom go function finder"
language: go
query: "(function_declaration) @func"
guidance: "Custom guidance"
`
	if err := os.WriteFile(filepath.Join(queriesDir, "override.yaml"), []byte(overrideYaml), 0o644); err != nil {
		t.Fatal(err)
	}

	// 2. Write list of queries containing a new one.
	newListYaml := `
- id: find_custom_thing
  description: "Custom search"
  language: go
  query: "(type_declaration) @type"
- id: find_comments
  description: "Custom comments"
  language: go
  query: "(comment) @comment"
`
	if err := os.WriteFile(filepath.Join(queriesDir, "custom.yaml"), []byte(newListYaml), 0o644); err != nil {
		t.Fatal(err)
	}

	r := NewQueryRegistry()
	if err := r.Reload(tempDir); err != nil {
		t.Fatal(err)
	}

	// Verify override.
	q, ok := r.GetTemplate("go", "find_functions")
	if !ok {
		t.Fatal("expected find_functions to be present")
	}
	if q != "(function_declaration) @func" {
		t.Errorf("expected overridden query, got %q", q)
	}

	cap, ok := r.GetCapability("go", "find_functions")
	if !ok {
		t.Fatal("expected capability to exist")
	}
	if cap.Description != "Custom go function finder" {
		t.Errorf("expected overridden description, got %q", cap.Description)
	}
	if cap.Guidance != "Custom guidance" {
		t.Errorf("expected overridden guidance, got %q", cap.Guidance)
	}

	// Verify new query from list.
	qCustom, ok := r.GetTemplate("go", "find_custom_thing")
	if !ok {
		t.Fatal("expected new query to be loaded")
	}
	if qCustom != "(type_declaration) @type" {
		t.Errorf("expected new query, got %q", qCustom)
	}

	// Verify list query override.
	qComment, ok := r.GetTemplate("go", "find_comments")
	if !ok {
		t.Fatal("expected find_comments to exist")
	}
	if qComment != "(comment) @comment" {
		t.Errorf("expected overridden comments, got %q", qComment)
	}
}

func formatNode(node *sitter.Node, code []byte, indent string) string {
	res := ""
	if node.IsNamed() {
		res += fmt.Sprintf("%s%s [%d-%d] -> %q\n", indent, node.Kind(), node.StartByte(), node.EndByte(), node.Utf8Text(code))
	}
	for i := uint(0); i < node.ChildCount(); i++ {
		res += formatNode(node.Child(i), code, indent+"  ")
	}
	return res
}

func TestFindTSFunctions(t *testing.T) {
	parser := sitter.NewParser()
	lang := lang_typescript.GetLanguage()
	parser.SetLanguage(lang)

	code := []byte("const x = `hello`;")
	tree := parser.ParseCtx(context.Background(), code, nil)
	if tree == nil {
		t.Fatal("expected non-nil tree")
	}
	defer tree.Close()

	ast := formatNode(tree.RootNode(), code, "")
	if strings.Contains(ast, "ERROR") {
		t.Errorf("AST contains ERROR node:\n%s", ast)
	}
}
