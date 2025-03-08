package main

import (
	"fmt"
	"os"

	"github.com/Ninad-Bhangui/awstui/aws"
	"github.com/Ninad-Bhangui/awstui/ui"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	layout := ui.NewLayout(app)
	pages := tview.NewPages()

	var currentConfig config.Config

	// Create profile selector
	profileSelector := ui.NewProfileSelector(func(profile aws.Profile, cfg config.Config) {
		// Store config for later use
		currentConfig = cfg

		// Create home screen
		homeScreen := ui.NewHomeScreen(func(service string) {
			showResourceList(app, layout, pages, service, currentConfig)
		})

		// Update layout
		layout.SetContent(homeScreen)
		layout.SetContext(fmt.Sprintf("Profile: %s", profile.Name))
		layout.SetKeybindings("<?> Help • <:> Quick Nav • <q> Back")
		layout.SetStatus("Ready")
	})

	// Load profiles
	if err := profileSelector.LoadProfiles(); err != nil {
		fmt.Printf("Error loading profiles: %v\n", err)
		os.Exit(1)
	}

	// Set up initial screen
	layout.SetContent(profileSelector)
	layout.SetContext("Select AWS Profile")
	layout.SetKeybindings("<?> Help • <q> Quit")

	// Set up pages
	pages.AddPage("main", layout, true, true)

	// Set up input capture for quick navigation
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Let layout handle its own keys first
		if layout.GetContent() == profileSelector {
			return event
		}

		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case ':':
				// Show resource prompt
				var resourceType string
				var modal *tview.InputField
				modal = tview.NewInputField().
					SetLabel(":").
					SetFieldWidth(20).
					SetDoneFunc(func(key tcell.Key) {
						if key == tcell.KeyEnter {
							resourceType = modal.GetText()
							pages.RemovePage("modal")
							showResourceList(app, layout, pages, resourceType, currentConfig)
						} else if key == tcell.KeyEscape {
							pages.RemovePage("modal")
							app.SetFocus(layout.GetContent())
						}
					})

				modal.SetBorder(true)
				modal.SetTitle("Quick Navigation")
				modal.SetTitleAlign(tview.AlignLeft)

				pages.AddPage("modal", modal, true, true)
				app.SetFocus(modal)
				return nil
			case 'q':
				if layout.GetContent() != profileSelector {
					// Return to profile selector
					layout.SetContent(profileSelector)
					layout.SetContext("Select AWS Profile")
					layout.SetKeybindings("<?> Help • <q> Quit")
					return nil
				}
				app.Stop()
				return nil
			}
		case tcell.KeyEscape:
			if layout.GetContent() != profileSelector {
				// Return to profile selector
				layout.SetContent(profileSelector)
				layout.SetContext("Select AWS Profile")
				layout.SetKeybindings("<?> Help • <q> Quit")
				return nil
			}
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}

func showResourceList(app *tview.Application, layout *ui.Layout, pages *tview.Pages, resourceType string, cfg config.Config) {
	switch resourceType {
	case "ec2", "ecr", "lambda", "secrets":
		list := ui.NewResourceList(resourceType, cfg)
		layout.SetContent(list)
		layout.SetContext(fmt.Sprintf("Viewing %s", list.GetTitle()))
		layout.SetKeybindings("<?> Help • <:> Quick Nav • <q> Back")
		app.SetFocus(list)
	default:
		// Show error modal
		modal := tview.NewModal().
			SetText(fmt.Sprintf("Unknown resource type: %s", resourceType)).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("error")
				app.SetFocus(layout.GetContent())
			})

		pages.AddPage("error", modal, true, true)
		app.SetFocus(modal)
	}
}
