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

	// Create a function to handle profile selection
	onProfileSelect := func(profile aws.Profile, cfg config.Config) {
		app.Stop()
	}

	// Create profile selector
	profileSelector := ui.NewProfileSelector(app, onProfileSelect)
	if err := profileSelector.LoadProfiles(); err != nil {
		panic(err)
	}

	// Handle Esc/q to quit
	profileSelector.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

	if err := app.SetRoot(profileSelector, true).Run(); err != nil {
		panic(err)
	}
}
