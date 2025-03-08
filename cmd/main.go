package main

import (
	"fmt"

	"github.com/Ninad-Bhangui/awstui/aws"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Get available AWS profiles
	profiles, err := aws.GetProfiles()
	if err != nil {
		panic(err)
	}

	// Create a list of profile information
	var profileInfo []string
	var selectedIndex int
	for i, p := range profiles {
		info := p.Name
		if p.Region != "" {
			info += fmt.Sprintf(" (region: %s)", p.Region)
		}
		if p.IsDefault {
			info += " [default]"
		}
		if p.IsFromEnv {
			info += " [active]"
			selectedIndex = i
		}
		profileInfo = append(profileInfo, info)
	}

	// Create profile list
	list := tview.NewList().
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorBlue)

	// Add profiles to the list
	for _, info := range profileInfo {
		list.AddItem(info, "", 0, nil)
	}

	// Set the current selection to the active profile
	list.SetCurrentItem(selectedIndex)

	// Create a frame around the list with a title
	frame := tview.NewFrame(list).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("AWS Profiles", true, tview.AlignCenter, tcell.ColorWhite)

	if err := app.SetRoot(frame, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
