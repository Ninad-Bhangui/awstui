package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Layout represents the main application layout
type Layout struct {
	*tview.Grid
	header      *tview.Flex
	context     *tview.TextView
	keybindings *tview.TextView
	content     tview.Primitive
}

// NewLayout creates a new application layout
func NewLayout() *Layout {
	layout := &Layout{
		Grid:        tview.NewGrid(),
		header:      tview.NewFlex(),
		context:     tview.NewTextView(),
		keybindings: tview.NewTextView(),
	}

	// Set up header
	layout.context.
		SetTextColor(tcell.ColorWhite)

	layout.keybindings.
		SetTextAlign(tview.AlignRight).
		SetTextColor(tcell.ColorWhite)

	layout.header.AddItem(layout.context, 0, 1, false)
	layout.header.AddItem(layout.keybindings, 0, 1, false)

	// Set up grid
	layout.Grid.SetRows(1, 0) // Header row and content row
	layout.Grid.SetColumns(0) // Full width
	layout.Grid.SetBorder(false)

	// Add header to grid
	layout.Grid.AddItem(layout.header, 0, 0, 1, 1, 0, 0, false)

	return layout
}

// SetContent sets the main content area
func (l *Layout) SetContent(content tview.Primitive) {
	l.content = content
	l.Grid.AddItem(content, 1, 0, 1, 1, 0, 0, true)
}

// SetContext sets the context text in the header
func (l *Layout) SetContext(text string) {
	l.context.SetText(text)
}

// SetKeybindings sets the keybindings text in the header
func (l *Layout) SetKeybindings(text string) {
	l.keybindings.SetText(text)
}

// GetContent returns the current content primitive
func (l *Layout) GetContent() tview.Primitive {
	return l.content
}
