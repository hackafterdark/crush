package tools

import (
	"context"
	"fmt"
	"slices"

	"charm.land/fantasy"
	"github.com/charmbracelet/crush/internal/agent/tools/mcp"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/charmbracelet/crush/internal/otel"
	"github.com/charmbracelet/crush/internal/permission"
	"go.opentelemetry.io/otel/attribute"
)

// mcpTransportAttr returns OTel attributes describing the MCP server's
// transport (stdio/pipe, HTTP, etc.) based on the client session state.
func mcpTransportAttr(mcpName string) []attribute.KeyValue {
	info, ok := mcp.GetState(mcpName)
	if !ok {
		return nil
	}
	if info.Client == nil {
		return nil
	}
	attrs := []attribute.KeyValue{
		attribute.String("mcp.session.id", mcpName),
		attribute.String("mcp.protocol.version", "2025-06-18"),
	}
	if info.Client.Transport() != "" {
		attrs = append(attrs, attribute.String("network.transport", info.Client.Transport()))
	}
	if info.Client.TransportURL() != "" {
		attrs = append(attrs, attribute.String("network.protocol.name", info.Client.TransportURL()))
	}
	return attrs
}

// whitelistDockerTools contains Docker MCP tools that don't require permission.
var whitelistDockerTools = []string{
	"mcp_docker_mcp-find",
	"mcp_docker_mcp-add",
	"mcp_docker_mcp-remove",
	"mcp_docker_mcp-config-set",
	"mcp_docker_code-mode",
}

// GetMCPTools gets all the currently available MCP tools.
func GetMCPTools(permissions permission.Service, cfg *config.ConfigStore, wd string) []*Tool {
	var result []*Tool
	for mcpName, tools := range mcp.Tools() {
		for _, tool := range tools {
			result = append(result, &Tool{
				mcpName:     mcpName,
				tool:        tool,
				permissions: permissions,
				workingDir:  wd,
				cfg:         cfg,
			})
		}
	}
	return result
}

// Tool is a tool from a MCP.
type Tool struct {
	mcpName         string
	tool            *mcp.Tool
	cfg             *config.ConfigStore
	permissions     permission.Service
	workingDir      string
	providerOptions fantasy.ProviderOptions
}

func (m *Tool) SetProviderOptions(opts fantasy.ProviderOptions) {
	m.providerOptions = opts
}

func (m *Tool) ProviderOptions() fantasy.ProviderOptions {
	return m.providerOptions
}

func (m *Tool) Name() string {
	return fmt.Sprintf("mcp_%s_%s", m.mcpName, m.tool.Name)
}

func (m *Tool) MCP() string {
	return m.mcpName
}

func (m *Tool) MCPToolName() string {
	return m.tool.Name
}

func (m *Tool) Info() fantasy.ToolInfo {
	parameters := make(map[string]any)
	required := make([]string, 0)

	if input, ok := m.tool.InputSchema.(map[string]any); ok {
		if props, ok := input["properties"].(map[string]any); ok {
			parameters = props
		}
		if req, ok := input["required"].([]any); ok {
			// Convert []any -> []string when elements are strings
			for _, v := range req {
				if s, ok := v.(string); ok {
					required = append(required, s)
				}
			}
		} else if reqStr, ok := input["required"].([]string); ok {
			// Handle case where it's already []string
			required = reqStr
		}
	}

	return fantasy.ToolInfo{
		Name:        m.Name(),
		Description: m.tool.Description,
		Parameters:  parameters,
		Required:    required,
	}
}

func (m *Tool) Run(ctx context.Context, params fantasy.ToolCall) (fantasy.ToolResponse, error) {
	ctx, span := otel.StartSpan(ctx, "execute_tool mcp")
	defer span.End()
	span.SetAttributes(
		attribute.String("gen_ai.tool.name", m.Name()),
		attribute.String("gen_ai.tool.call.id", params.ID),
		attribute.String("gen_ai.tool.call.arguments", params.Input),
		// MCP-specific attributes to distinguish from native tools.
		attribute.String("mcp.method.name", "tools/call"),
		attribute.String("gen_ai.operation.name", "execute_tool"),
	)
	// Add transport attributes from the MCP session state.
	for _, attr := range mcpTransportAttr(m.mcpName) {
		span.SetAttributes(attr)
	}
	sessionID := GetSessionFromContext(ctx)
	if sessionID == "" {
		return fantasy.ToolResponse{}, fmt.Errorf("session ID is required for creating a new file")
	}

	// Skip permission for whitelisted Docker MCP tools.
	if !slices.Contains(whitelistDockerTools, params.Name) {
		permissionDescription := fmt.Sprintf("execute %s with the following parameters:", m.Info().Name)
		p, err := m.permissions.Request(
			ctx,
			permission.CreatePermissionRequest{
				SessionID:   sessionID,
				ToolCallID:  params.ID,
				Path:        m.workingDir,
				ToolName:    m.Info().Name,
				Action:      "execute",
				Description: permissionDescription,
				Params:      params.Input,
			},
		)
		if err != nil {
			return fantasy.ToolResponse{}, err
		}
		if !p {
			return NewPermissionDeniedResponse(), nil
		}
	}

	result, err := mcp.RunTool(ctx, m.cfg, m.mcpName, m.tool.Name, params.Input)
	if err != nil {
		otel.SetErrorStatus(span, err.Error())
		return fantasy.NewTextErrorResponse(err.Error()), nil
	}

	var response fantasy.ToolResponse
	switch result.Type {
	case "image", "media":
		if !GetSupportsImagesFromContext(ctx) {
			modelName := GetModelNameFromContext(ctx)
			return fantasy.NewTextErrorResponse(fmt.Sprintf("This model (%s) does not support image data.", modelName)), nil
		}

		if result.Type == "image" {
			response = fantasy.NewImageResponse(result.Data, result.MediaType)
		} else {
			response = fantasy.NewMediaResponse(result.Data, result.MediaType)
		}
		response.Content = result.Content
	default:
		response = fantasy.NewTextResponse(result.Content)
	}

	// Record the tool result on the span (opt-in per MCP semconv).
	if response.Content != "" {
		span.SetAttributes(attribute.String("gen_ai.tool.call.result", response.Content))
	}
	return response, nil
}
