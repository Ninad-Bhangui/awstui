package ui

import (
	"fmt"

	"github.com/Ninad-Bhangui/awstui/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HomePage represents the main application page after profile selection
type HomePage struct {
	*tview.TextView
	profile aws.Profile
	cfg     config.Config
}

// NewHomePage creates a new home page
func NewHomePage(profile aws.Profile, cfg config.Config) *HomePage {
	home := &HomePage{
		TextView: tview.NewTextView(),
		profile:  profile,
		cfg:      cfg,
	}

	// Basic setup
	home.SetBorder(true)
	home.SetTitle("AWS TUI")
	home.SetTitleAlign(tview.AlignLeft)
	home.SetTextColor(tcell.ColorWhite)

	// Set some dummy content
	content := fmt.Sprintf("Welcome to AWS TUI!\n\n"+
		"Active Profile: %s\n"+
		"Region: %s\n\n"+
		"More features coming soon...",
		profile.Name,
		profile.Region)

	home.SetText(content)

	return home
}
