package agent

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"

	"charm.land/fantasy"
)

const (
	loopDetectionWindowSize  = 10
	loopDetectionMaxRepeats  = 5
	reasoningLoopWindowSize  = 5
	reasoningLoopMaxRepeats  = 2
)

// hasRepeatedToolCalls checks whether the agent is stuck in a loop by looking
// at recent steps. It examines the last windowSize steps and returns true if
// any tool-call signature appears more than maxRepeats times.
func hasRepeatedToolCalls(steps []fantasy.StepResult, windowSize, maxRepeats int) bool {
	if len(steps) < windowSize {
		return false
	}

	window := steps[len(steps)-windowSize:]
	counts := make(map[string]int)

	for _, step := range window {
		sig := getToolInteractionSignature(step.Content)
		if sig == "" {
			continue
		}
		counts[sig]++
		if counts[sig] > maxRepeats {
			return true
		}
	}

	return false
}

// getToolInteractionSignature computes a hash signature for the tool
// interactions in a single step's content. It pairs tool calls with their
// results (matched by ToolCallID) and returns a hex-encoded SHA-256 hash.
// If the step contains no tool calls, it returns "".
func getToolInteractionSignature(content fantasy.ResponseContent) string {
	toolCalls := content.ToolCalls()
	if len(toolCalls) == 0 {
		return ""
	}

	// Index tool results by their ToolCallID for fast lookup.
	resultsByID := make(map[string]fantasy.ToolResultContent)
	for _, tr := range content.ToolResults() {
		resultsByID[tr.ToolCallID] = tr
	}

	h := sha256.New()
	for _, tc := range toolCalls {
		output := ""
		if tr, ok := resultsByID[tc.ToolCallID]; ok {
			output = toolResultOutputString(tr.Result)
		}
		io.WriteString(h, tc.ToolName)
		io.WriteString(h, "\x00")
		io.WriteString(h, tc.Input)
		io.WriteString(h, "\x00")
		io.WriteString(h, output)
		io.WriteString(h, "\x00")
	}
	return hex.EncodeToString(h.Sum(nil))
}

// toolResultOutputString converts a ToolResultOutputContent to a stable string
// representation for signature comparison.
func toolResultOutputString(result fantasy.ToolResultOutputContent) string {
	if result == nil {
		return ""
	}
	if text, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentText](result); ok {
		return text.Text
	}
	if errResult, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentError](result); ok {
		if errResult.Error != nil {
			return errResult.Error.Error()
		}
		return ""
	}
	if media, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentMedia](result); ok {
		return media.Data
	}
	return ""
}

// isReasoningOnlyStep returns the reasoning text from a step if the step
// contains only reasoning content (no tool calls and no final text).
// It returns the combined reasoning text and true if the step is reasoning-only,
// or an empty string and false otherwise.
func isReasoningOnlyStep(content fantasy.ResponseContent) (string, bool) {
	if len(content) == 0 {
		return "", false
	}

	var reasoningText strings.Builder
	hasReasoning := false
	hasOther := false

	for _, part := range content {
		switch part.(type) {
		case fantasy.ReasoningContent:
			hasReasoning = true
		case fantasy.ToolCallContent, fantasy.ToolResultContent:
			hasOther = true
		case fantasy.TextContent:
			// Text content at the step level is the final response text,
			// not reasoning. A step with final text is considered "progress".
			hasOther = true
		}
	}

	if !hasReasoning || hasOther {
		return "", false
	}

	// Collect all reasoning text from this step.
	for _, part := range content {
		if rc, ok := part.(fantasy.ReasoningContent); ok {
			reasoningText.WriteString(rc.Text)
		}
	}

	return reasoningText.String(), true
}

// hasRepeatedThinking detects when the model is stuck in a thinking loop —
// repeatedly producing reasoning content without making progress (no tool
// calls, no final text). It returns true when the same reasoning content
// repeats across consecutive steps.
func hasRepeatedThinking(steps []fantasy.StepResult) bool {
	if len(steps) < reasoningLoopWindowSize {
		return false
	}

	// Look at the last window of steps.
	window := steps[len(steps)-reasoningLoopWindowSize:]

	// Collect reasoning-only step texts.
	var reasoningTexts []string
	for _, step := range window {
		if text, ok := isReasoningOnlyStep(step.Content); ok {
			reasoningTexts = append(reasoningTexts, text)
		}
	}

	// We need at least 2 reasoning-only steps to detect a loop.
	if len(reasoningTexts) < 2 {
		return false
	}

	// Check if any reasoning text repeats more than maxRepeats times.
	counts := make(map[string]int)
	for _, text := range reasoningTexts {
		counts[text]++
		if counts[text] > reasoningLoopMaxRepeats {
			return true
		}
	}

	return false
}
