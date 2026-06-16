package parser

// Templates provides human-readable query names mapped to tree-sitter S-expression patterns.
// The outer key is the language name (e.g., "go", "typescript"), the inner key is the template name.
var Templates = map[string]map[string]string{
	"go": {
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
    (field_declaration_list)))
`,

		"find_variables": `
(var_declaration
  (var_spec
    name: (identifier) @name
    value: (_) @value))

(var_declaration
  (var_spec
    name: (identifier)))

(short_var_declaration
  (expression_list
    (identifier) @name)
  (expression_list))

(short_var_declaration
  (expression_list
    (identifier) @name)
  (expression_list
    (_) @value))
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
	},

	"typescript": {
		"find_functions": `
(function_declaration
  name: (identifier) @name
  parameters: (formal_parameters) @parameters
  body: (statement_block) @body)

(arrow_function
  parameters: (formal_parameters) @parameters
  body: (statement_block) @body)

(function_expression
  name: (identifier) @name
  parameters: (formal_parameters) @parameters
  body: (statement_block) @body)
`,

		"find_structs": `
(class_declaration
  name: (type_identifier) @name
  body: (class_body
    (method_definition
      name: (property_identifier) @method_name
      parameters: (formal_parameters) @parameters
      body: (statement_block) @body)))

(class_declaration
  name: (type_identifier) @name
  body: (class_body
    (field_definition
      name: (property_identifier) @field_name
      type: (type_annotation
        type: (_) @field_type))))

(type_alias_declaration
  name: (type_identifier) @name
  value: (_) @type_body)
`,

		"find_variables": `
(variable_declarator
  name: (identifier) @name
  value: (_) @value)

(variable_declarator
  name: (identifier) @name)

(variable_declarator
  name: (array_pattern) @name
  value: (_) @value)

(variable_declarator
  name: (object_pattern) @name
  value: (_) @value)
`,

		"find_interfaces": `
(interface_declaration
  name: (type_identifier) @name
  body: (interface_body
    (property_signature
      name: (property_identifier) @method_name
      type: (type_annotation
        type: (_) @method_type))) @interface_body)
`,

		"find_calls": `
(call_expression
  function: (identifier) @function_name
  arguments: (arguments) @arguments)

(call_expression
  function: (member_expression
    property: (property_identifier) @method_name)
  arguments: (arguments) @arguments)
`,

		"find_imports": `
(import_statement
  source: (string) @import_path
  (import_clause
    (named_imports
      (import_specifier
        name: (_) @import_name)))

(import_statement
  source: (string) @import_path
  (import_clause
    (namespace_import
      name: (_) @import_name)))

(import_statement
  source: (string) @import_path
  (import_require_clause
    (identifier) @import_name))
`,

		"find_comments": `
(comment) @comment
`,
	},

	"javascript": {
		"find_functions": `
(function_declaration
  name: (identifier) @name
  parameters: (formal_parameters) @parameters
  body: (statement_block) @body)

(arrow_function
  parameters: (formal_parameters) @parameters
  body: (statement_block) @body)

(function_expression
  name: (identifier) @name
  parameters: (formal_parameters) @parameters
  body: (statement_block) @body)
`,

		"find_structs": `
(class_declaration
  name: (identifier) @name
  body: (class_body
    (method_definition
      name: (property_identifier) @method_name
      parameters: (formal_parameters) @parameters
      body: (statement_block) @body)))

(class_declaration
  name: (identifier) @name
  body: (class_body
    (field_definition
      name: (property_identifier) @field_name
      value: (_) @field_value)))

(class
  name: (identifier) @name
  body: (class_body
    (method_definition
      name: (property_identifier) @method_name
      parameters: (formal_parameters) @parameters
      body: (statement_block) @body)))
`,

		"find_variables": `
(variable_declarator
  name: (identifier) @name
  value: (_) @value)

(variable_declarator
  name: (identifier) @name)

(variable_declarator
  name: (array_pattern) @name
  value: (_) @value)

(variable_declarator
  name: (object_pattern) @name
  value: (_) @value)
`,

		"find_interfaces": `
`,

		"find_calls": `
(call_expression
  function: (identifier) @function_name
  arguments: (arguments) @arguments)

(call_expression
  function: (member_expression
    property: (property_identifier) @method_name)
  arguments: (arguments) @arguments)
`,

		"find_imports": `
(import_statement
  source: (string) @import_path
  (import_clause
    (named_imports
      (import_specifier
        name: (_) @import_name)))

(import_statement
  source: (string) @import_path
  (import_clause
    (namespace_import
      name: (_) @import_name)))
`,

		"find_comments": `
(comment) @comment
`,
	},

	"python": {
		"find_functions": `
(function_definition
  name: (identifier) @name
  parameters: (parameters) @parameters
  body: (block) @body)
`,

		"find_structs": `
(class_definition
  name: (identifier) @name
  body: (block) @body)
`,

		"find_variables": `
(assignment
  left: (identifier) @name
  right: (_) @value)

(assignment
  left: (identifier) @name)

(assignment
  left: (pattern_list) @name
  right: (_) @value)
`,

		"find_interfaces": `
`,

		"find_calls": `
(call
  function: (_) @function_name
  arguments: (_))

(call
  function: (_) @method_name
  arguments: (_))
`,

		"find_imports": `
(import_statement
  (aliased_import
    name: (dotted_name) @package_name
    alias: (identifier) @import_alias))

(import_statement
  (dotted_name) @package_name)

(import_from_statement
  module_name: (dotted_name) @import_path
  (aliased_import
    name: (dotted_name) @import_name
    alias: (identifier) @import_alias))

(import_from_statement
  module_name: (dotted_name) @import_path
  (dotted_name) @import_name)
`,

		"find_comments": `
(comment) @comment
`,
	},

	"sql": {
		"find_functions": `
(create_function
  name: (object_reference
    name: (identifier) @name)
  body: (function_body) @body)
`,

		"find_structs": `
(create_table
  name: (object_reference
    name: (identifier) @name)
  body: (column_definitions) @body)
`,

		"find_select_tables": `
(select_statement
  (from_clause
    (table_name) @table_name))
`,

		"find_joins": `
(join_clause
  (table_name) @joined_table)
`,

		"find_inserts": `
(insert_statement
  (table_name) @table_name)
`,

		"find_deletes": `
(delete_statement
  (from_clause
    (table_name) @table_name))
`,

		"find_select_all": `
(select_statement
  (select_list
    (wildcard)))
`,

		"find_variables": `
`,

		"find_interfaces": `
`,

		"find_calls": `
`,

		"find_imports": `
`,

		"find_comments": `
(comment) @comment
`,
	},

	"rust": {
		"find_functions": `
(function_item
  name: (identifier) @name
  parameters: (parameters) @parameters
  body: (block) @body)

(function_item
  name: (identifier) @name
  body: (block) @body)
`,

		"find_structs": `
(struct_item
  name: (type_identifier) @name
  body: (_) @struct_body)
`,

		"find_variables": `
(let_declaration
  pattern: (identifier) @name
  value: (_) @value)

(let_declaration
  pattern: (identifier) @name)
`,

		"find_interfaces": `
(trait_item
  name: (identifier) @name)
`,

		"find_calls": `
(call_expression
  function: (_) @function_name)
`,

		"find_imports": `
(use_declaration
  name: (scoped_identifier
    path: (_) @package_name
    name: (_) @import_name))

(use_declaration
  name: (identifier) @import_name)
`,

		"find_comments": `
(line_comment) @comment

(block_comment) @comment
`,
	},

	// "java" — Java not supported (requires external scanner not present in vendored grammar)
	"php": {
		"find_functions": `
(function_definition
  name: (name) @name
  body: (compound_statement) @body)

(method_declaration
  name: (name) @name
  body: (compound_statement) @body)

(anonymous_function
  parameters: (formal_parameters) @parameters
  body: (compound_statement) @body)
`,

		"find_structs": `
(class_declaration
  name: (_) @name
  body: (declaration_list) @class_body)
`,

		"find_variables": `
(variable_name) @name

(property_declaration
  (property_element
    name: (variable_name) @name
    default_value: (_) @value))

(property_declaration
  (property_element
    name: (variable_name) @name))

(assignment_expression
  left: (variable_name) @name
  right: (_) @value)
`,

		"find_interfaces": `
(interface_declaration
  name: (name) @name
  body: (declaration_list) @interface_body)
`,

		"find_calls": `
(function_call_expression
  function: (_) @function_name
  arguments: (arguments) @arguments)

(member_call_expression
  object: (_)
  (name) @method_name
  arguments: (arguments) @arguments)

(scoped_call_expression
  scope: (_)
  (name) @method_name
  arguments: (arguments) @arguments)
`,

		"find_imports": `
(namespace_use_declaration
  name: (qualified_name) @import_path
  alias: (name) @import_name)

(namespace_use_declaration
  name: (qualified_name) @import_path)

(use_declaration
  name: (qualified_name) @import_path
  alias: (name) @import_name)
`,

		"find_comments": `
(comment) @comment
`,
	},

	"cpp": {
		"find_functions": `
(function_definition
  (function_declarator
    (identifier) @name
    (parameter_list) @parameters)
  (compound_statement) @body)

(function_definition
  (function_declarator
    (field_identifier) @name
    (parameter_list) @parameters)
  (compound_statement) @body)
`,

		"find_structs": `
(class_specifier
  name: (type_identifier) @name
  body: (field_declaration_list) @body)

(struct_specifier
  name: (type_identifier) @name
  body: (field_declaration_list) @body)
`,

		"find_variables": `
(declaration
  (init_declarator
    (identifier) @name
    value: (_) @value))

(assignment_expression
  left: (identifier) @name
  right: (_) @value)
`,

		"find_interfaces": ``,

		"find_calls": `
(call_expression
  function: (identifier) @function_name
  arguments: (argument_list) @arguments)

(call_expression
  function: (qualified_identifier
    (identifier) @function_name)
  arguments: (argument_list) @arguments)

(call_expression
  function: (field_expression
    field: (field_identifier) @method_name)
  arguments: (argument_list) @arguments)
`,

		"find_imports": `
(preproc_include
  path: (system_lib_string) @import_path)

(preproc_include
  path: (string_literal) @import_path)
`,

		"find_comments": `
(comment) @comment
`,
	},

	"hcl": {
		"find_functions": `
(attribute
  name: (identifier) @name
  value: (expression) @body)
`,

		"find_structs": `
(block
  type: (identifier) @name
  labels: (block_label) @labels
  body: (body) @body)
`,

		"find_variables": `
(attribute
  name: (identifier) @name
  value: (expression) @value)
`,

		"find_interfaces": ``,

		"find_calls": `
(function_call
  function_name: (identifier) @function_name
  arguments: (arguments) @arguments)
`,

		"find_imports": ``,

		"find_comments": `
(comment) @comment
`,
	},

	"csharp": {
		"find_functions": `
(method_declaration
  name: (identifier) @name
  parameters: (parameter_list) @parameters
  body: (block) @body)

(method_declaration
  name: (identifier) @name
  parameters: (parameter_list) @parameters
  body: (arrow_expression_clause) @body)

(constructor_declaration
  name: (identifier) @name
  parameters: (parameter_list) @parameters
  body: (block) @body)
`,

		"find_structs": `
(class_declaration
  name: (identifier) @name
  body: (declaration_list) @body)

(record_declaration
  name: (identifier) @name
  body: (declaration_list) @body)
`,

		"find_variables": `
(local_declaration_statement
  (variable_declaration
    (variable_declarator
      (identifier) @name)))
`,

		"find_interfaces": ``,

		"find_calls": `
(invocation_expression
  function: (identifier) @function_name
  arguments: (argument_list) @arguments)

(invocation_expression
  function: (member_access_expression
    (identifier) @method_name)
  arguments: (argument_list) @arguments)
`,

		"find_imports": `
(using_directive
  (identifier) @import_path)

(using_directive
  (qualified_name
    (identifier) @import_path))
`,

		"find_comments": `
(comment) @comment
`,
	},

	"ruby": {
		"find_functions": `
(method
  name: (_) @name)

(singleton_method
  name: (_) @name)
`,

		"find_structs": `
(class
  name: (_) @name)

(class
  (constant) @name)
`,

		"find_variables": `
(assignment
  left: (_) @name)
`,

		"find_interfaces": ``,

		"find_calls": `
(call
  method: (_) @method_name)
`,

		"find_imports": ``,

		"find_comments": `
(comment) @comment
`,
	},

	"json": {
		"find_functions": ``,

		"find_structs": `
(object
  (pair
    key: (string) @name
    value: (object) @body))
`,

		"find_variables": `
(pair
  key: (string) @name
  value: (_) @value)
`,

		"find_interfaces": ``,

		"find_calls": ``,

		"find_imports": ``,

		"find_comments": `
(comment) @comment
`,
	},

	"html": {
		"find_functions": ``,

		"find_structs": `
(element
  name: (tag_name) @name
  children: (element_children) @body)
`,

		"find_variables": `
(attribute
  name: (attribute_name) @name
  value: (attribute_value) @value)
`,

		"find_interfaces": ``,

		"find_calls": ``,

		"find_imports": `
(element
  name: (tag_name) @import_path
  (start_tag
    (attribute
      name: (attribute_name) @import_name)))
`,

		"find_comments": `
(comment) @comment
`,
	},

	"css": {
		"find_functions": ``,

		"find_structs": `
(rule_set
  (selectors) @name
  (block) @body)

(media_statement
  (feature_query) @name
  (block) @body)
`,

		"find_variables": `
(declaration
  (property_name) @name
  (call_expression) @value)

(declaration
  (property_name) @name
  (color_value) @value)

(declaration
  (property_name) @name
  (integer_value) @value)

(declaration
  (property_name) @name
  (float_value) @value)

(declaration
  (property_name) @name
  (string_value) @value)

(declaration
  (property_name) @name
  (plain_value) @value)
`,

		"find_interfaces": ``,

		"find_calls": `
(call_expression
  (function_name) @name
  (arguments) @arguments)
`,

		"find_imports": `
(import_statement
  (string) @import_path)
`,

		"find_comments": `
(comment) @comment
`,
	},

	"toml": {
		"find_functions": ``,

		"find_structs": `
(table
  name: (key) @name
  value: (array) @body)
`,

		"find_variables": `
(pair
  key: (key) @name
  value: (_) @value)
`,

		"find_interfaces": ``,

		"find_calls": ``,

		"find_imports": ``,

		"find_comments": `
(comment) @comment
`,
	},

	"scala": {
		"find_functions": `
(function_definition
  name: (_) @name)

(function_definition
  name: (_) @name
  parameters: (_) @parameters)

(function_definition
  name: (_) @name
  parameters: (_) @parameters
  body: (_) @body)

(function_declaration
  name: (_) @name)
`,

		"find_structs": `
(class_definition
  name: (_) @name
  body: (_) @body)

(object_definition
  name: (_) @name
  body: (_) @body)
`,

		"find_variables": `
(val_definition
  pattern: (_) @name
  value: (_) @value)

(val_definition
  pattern: (_) @name)

(var_definition
  pattern: (_) @name
  value: (_) @value)

(var_definition
  pattern: (_) @name)
`,

		"find_interfaces": `
(trait_definition
  name: (_) @name
  body: (_) @body)
`,

		"find_calls": `
(call_expression
  function: (_) @function_name
  arguments: (_))
`,

		"find_imports": `
(import_declaration
  path: (_) @import_path)
`,

		"find_comments": `
(comment) @comment
`,
	},

	"c": {
		"find_functions": `
(function_definition
  declarator: (function_declarator
    parameters: (parameter_list) @parameters)
  body: (compound_statement) @body)
`,

		"find_structs": `
(struct_specifier
  name: (type_identifier) @name
  body: (field_declaration_list) @body)

(type_definition
  (struct_specifier
    body: (field_declaration_list) @body)
  (type_identifier) @name)
`,

		"find_variables": `
(init_declarator
  declarator: (identifier) @name
  value: (_) @value)
`,

		"find_interfaces": ``,

		"find_calls": `
(call_expression
  function: (identifier) @function_name
  arguments: (argument_list) @arguments)

(call_expression
  function: (field_expression
    field: (field_identifier) @method_name)
  arguments: (argument_list) @arguments)
`,

		"find_imports": `
(preproc_include
  path: (string_literal) @import_path)

(preproc_include
  path: (system_lib_string) @import_path)
`,

		"find_comments": `
(comment) @comment
`,
	},

	"bash": {
		"find_functions": `
(function_definition
  name: (word) @name
  body: (compound_statement) @body)

(function_definition
  name: (word) @name)
`,

		"find_structs": ``,

		"find_variables": `
(variable_assignment
  name: (variable_name) @name
  value: (_) @value)

(variable_assignment
  name: (variable_name) @name)
`,

		"find_interfaces": ``,

		"find_calls": `
(command
  name: (command_name (word)) @function_name)

(command
  name: (command_name (word)) @function_name
  argument: (_) @arguments)
`,

		"find_imports": ``,

		"find_comments": `
(comment) @comment
`,
	},
}

// GetTemplate returns the tree-sitter query for the given language and template name.
// Returns an empty string and false if the template doesn't exist.
func GetTemplate(lang, name string) (string, bool) {
	langTemplates, ok := Templates[lang]
	if !ok {
		return "", false
	}
	query, ok := langTemplates[name]
	return query, ok
}

// TemplateNames returns a sorted list of available template names for the given language.
func TemplateNames(lang string) []string {
	langTemplates, ok := Templates[lang]
	if !ok {
		return nil
	}
	names := make([]string, 0, len(langTemplates))
	for name := range langTemplates {
		names = append(names, name)
	}
	return names
}
