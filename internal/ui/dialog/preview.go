package dialog

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/crush/internal/message"
	"github.com/charmbracelet/crush/internal/ui/common"
	"github.com/charmbracelet/crush/internal/ui/styles"
	"github.com/charmbracelet/x/ansi"
	uv "github.com/charmbracelet/ultraviolet"
)

const (
	// PreviewID is the identifier for the attachment preview dialog.
	PreviewID = "preview"
	// previewDialogMaxWidth is the maximum width for the preview dialog.
	previewDialogMaxWidth = 80
	// previewDialogMaxHeight is the maximum height for the preview dialog.
	previewDialogMaxHeight = 30
)

// Preview represents a dialog for previewing attachment content.
type Preview struct {
	com    *common.Common
	attach message.Attachment
	t      *styles.Styles
}

var _ Dialog = (*Preview)(nil)

// NewPreview creates a new attachment preview dialog.
func NewPreview(com *common.Common, attach message.Attachment) *Preview {
	return &Preview{
		com:    com,
		attach: attach,
		t:      com.Styles,
	}
}

// ID implements Dialog.
func (p *Preview) ID() string {
	return PreviewID
}

// HandleMsg implements Dialog.
func (p *Preview) HandleMsg(msg tea.Msg) Action {
	_ = msg
	return nil
}

// Draw implements Dialog.
func (p *Preview) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := p.t
	width := max(0, min(previewDialogMaxWidth, area.Dx()))
	height := max(0, min(previewDialogMaxHeight, area.Dy()))
	innerWidth := width - t.Dialog.View.GetHorizontalFrameSize()
	heightOffset := t.Dialog.Title.GetVerticalFrameSize() +
		t.Dialog.View.GetVerticalFrameSize()

	content := string(p.attach.Content)
	// Truncate if content exceeds the dialog height.
	maxChars := innerWidth * (height - heightOffset)
	if len(content) > maxChars {
		content = ansi.Truncate(content, maxChars, "…")
	}

	rc := NewRenderContext(t, width)
	rc.Title = "Preview: " + p.attach.FileName
	rc.AddPart(content)
	rc.Help = t.Dialog.HelpView.Render("esc: close")

	view := rc.Render()

	DrawCenterCursor(scr, area, view, nil)
	return nil
}
