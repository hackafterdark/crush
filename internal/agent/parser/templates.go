package parser

// Templates provides human-readable query names mapped to tree-sitter S-expression patterns.
var Templates = map[string]string{
	"find_functions": `
(function_declaration
  name: (identifier) @name
  parameters: (parameter_list) @parameters
  body: (block) @body)

(function_declaration
  name: (identifier) @name
  parameters: (parameter_list) @parameters)

(method_declaration
  receiver: (parameter_list) @receiver
  name: (field_identifier) @name
  parameters: (parameter_list) @parameters
  body: (block) @body)
`,

	"find_structs": `
(type_spec
  name: (type_identifier) @name
  type: (struct_type
    fields: (struct_field_declaration_list
      (struct_field_declaration
        name: (field_identifier) @field_name
        type: (_) @field_type))) @struct_body)
`,

	"find_variables": `
(var_declaration
  (var_spec
    name: (identifier) @name
    value: (_) @value))

(var_declaration
  (var_spec
    name: (identifier) @name))
`,

	"find_interfaces": `
(type_spec
  name: (type_identifier) @name
  type: (interface_type
    (interface_type_elements
      (interface_type_element
        name: (field_identifier) @method_name
        type: (_) @method_type))) @interface_body)
`,

	"find_calls": `
(call_expression
  function: (identifier) @function_name
  arguments: (argument_list) @arguments)

(call_expression
  function: (selector_expression
    field: (field_identifier) @method_name)
  arguments: (argument_list) @arguments)
`,

	"find_imports": `
(import_declaration
  (import_spec
    name: (_) @package_name
    path: (interpreted_string_literal) @import_path))
`,

	"find_comments": `
(comment) @comment
`,
}

// GetTemplate returns the tree-sitter query for the given template name.
// Returns an empty string and false if the template doesn't exist.
func GetTemplate(name string) (string, bool) {
	query, ok := Templates[name]
	return query, ok
}

// TemplateNames returns a sorted list of available template names.
func TemplateNames() []string {
	names := make([]string, 0, len(Templates))
	for name := range Templates {
		names = append(names, name)
	}
	return names
}
