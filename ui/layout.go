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
	statusBar   *tview.TextView
	helpPanel   *tview.TextView
	app         *tview.Application
	showHelp    bool
}

// NewLayout creates a new application layout
func NewLayout(app *tview.Application) *Layout {
	layout := &Layout{
		Grid:        tview.NewGrid(),
		header:      tview.NewFlex(),
		context:     tview.NewTextView(),
		keybindings: tview.NewTextView(),
		statusBar:   tview.NewTextView(),
		helpPanel:   tview.NewTextView(),
		app:         app,
		showHelp:    false,
	}

	// Set up header
	layout.context.
		SetTextColor(tcell.ColorWhite)

	layout.keybindings.
		SetTextAlign(tview.AlignRight).
		SetTextColor(tcell.ColorWhite)

	layout.header.AddItem(layout.context, 0, 1, false)
	layout.header.AddItem(layout.keybindings, 0, 1, false)

	// Set up status bar
	layout.statusBar.
		SetTextColor(tcell.ColorWhite).
		SetText("Ready")

	// Set up help panel
	layout.helpPanel.SetBorder(true)
	layout.helpPanel.SetTitle("Help & Key Bindings")
	layout.helpPanel.SetTitleAlign(tview.AlignLeft)
	layout.helpPanel.SetText(`
[::b]Navigation Keys[::-]
  ↑/k         : Move up
  ↓/j         : Move down
  Enter       : Select item
  
[::b]General Commands[::-]
  ?           : Toggle help
  q/Esc       : Quit/Back
  
[::b]AWS Resources (coming soon)[::-]
  1           : EC2 Instances
  2           : S3 Buckets
  3           : Lambda Functions

[::b]Tips[::-]
  • Use vim-style navigation (j/k) for faster movement
  • Press ? again to return to previous screen
  • More features coming soon!`)

	// Set up grid
	layout.Grid.SetRows(1, 0, 1) // Header, content, status bar
	layout.Grid.SetColumns(0)    // Full width
	layout.Grid.SetBorder(false)

	// Add header and status bar to grid
	layout.Grid.AddItem(layout.header, 0, 0, 1, 1, 0, 0, false)
	layout.Grid.AddItem(layout.statusBar, 2, 0, 1, 1, 0, 0, false)

	// Set up input capture for help toggle
	layout.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '?' {
			layout.ToggleHelp()
			return nil
		}
		return event
	})

	return layout
}

// SetContent sets the main content area
func (l *Layout) SetContent(content tview.Primitive) {
	// Remove existing content if any
	if l.content != nil {
		l.Grid.RemoveItem(l.content)
	}

	l.content = content
	l.Grid.AddItem(content, 1, 0, 1, 1, 0, 0, true)
	l.app.SetFocus(content)
}

// SetContext sets the context text in the header
func (l *Layout) SetContext(text string) {
	l.context.SetText(text)
}

// SetKeybindings sets the keybindings text in the header
func (l *Layout) SetKeybindings(text string) {
	l.keybindings.SetText(text)
}

// SetStatus sets the status bar text
func (l *Layout) SetStatus(text string) {
	l.statusBar.SetText(text)
}

// ToggleHelp toggles the help panel visibility
func (l *Layout) ToggleHelp() {
	l.showHelp = !l.showHelp
	if l.showHelp {
		// Show help panel
		l.SetContent(l.helpPanel)
		l.SetContext("Help")
		l.SetKeybindings("<?> Back")
		l.SetStatus("Viewing help")
	} else {
		// Return to previous content
		l.SetContent(l.content)
		l.app.SetFocus(l.content)
	}
}

// GetContent returns the current content primitive
func (l *Layout) GetContent() tview.Primitive {
	return l.content
}
