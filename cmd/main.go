package main

import (
	"github.com/Ninad-Bhangui/awstui/aws"
	"github.com/Ninad-Bhangui/awstui/ui"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create main layout
	layout := ui.NewLayout()

	// Create a function to handle profile selection
	onProfileSelect := func(profile aws.Profile, cfg config.Config) {
		// Create and show home page
		homePage := ui.NewHomePage(profile, cfg)
		layout.SetContent(homePage)
		layout.SetContext(profile.Name)   // Show active profile in context
		layout.SetKeybindings("<q> Quit") // Update keybindings for home page
		app.SetFocus(homePage)
	}

	// Create profile selector
	profileSelector := ui.NewProfileSelector(onProfileSelect)
	if err := profileSelector.LoadProfiles(); err != nil {
		panic(err)
	}

	// Set initial content and header
	layout.SetContent(profileSelector)
	layout.SetContext("Select AWS Profile")
	layout.SetKeybindings("<↑/↓> Navigate • <Enter> Select • <Esc/q> Quit")

	// Handle global input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.Stop()
			return nil
		}
		switch event.Rune() {
		case 'q':
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
