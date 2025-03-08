package ui

import (
	"context"
	"fmt"

	"github.com/Ninad-Bhangui/awstui/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProfileSelector represents the profile selection screen
type ProfileSelector struct {
	*tview.List
	profileMgr *aws.ProfileManager
	onSelect   func(aws.Profile, config.Config)
}

// NewProfileSelector creates a new profile selection screen
func NewProfileSelector(onSelect func(aws.Profile, config.Config)) *ProfileSelector {
	selector := &ProfileSelector{
		List:       tview.NewList(),
		profileMgr: aws.NewProfileManager(),
		onSelect:   onSelect,
	}

	// Basic list setup
	selector.SetBorder(true)
	selector.SetTitle("AWS Profiles")
	selector.SetTitleAlign(tview.AlignLeft)
	selector.SetHighlightFullLine(true)
	selector.SetSelectedBackgroundColor(tcell.ColorBlue)

	// Set up selection handler
	selector.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if index >= 0 && index < len(selector.profileMgr.GetAllProfiles()) {
			profile := selector.profileMgr.GetAllProfiles()[index]

			// Try to load AWS config
			ctx := context.Background()
			cfg, err := selector.profileMgr.LoadConfig(ctx, profile.Name)
			if err != nil {
				// Just add error as new item for now
				selector.AddItem(fmt.Sprintf("Error: %v", err), "", 0, nil)
				return
			}
			onSelect(profile, cfg)
		}
	})

	return selector
}

// LoadProfiles loads and displays the AWS profiles
func (s *ProfileSelector) LoadProfiles() error {
	// Load profiles
	if err := s.profileMgr.LoadProfiles(); err != nil {
		s.AddItem(fmt.Sprintf("Error loading profiles: %v", err), "", 0, nil)
		return err
	}

	// Clear existing items
	s.Clear()

	// Add profiles to list
	var activeProfileIndex int
	for i, p := range s.profileMgr.GetAllProfiles() {
		name := p.Name
		if p.Region != "" {
			name += fmt.Sprintf(" (region: %s)", p.Region)
		}
		if p.IsDefault {
			name += " [default]"
		}
		if p.IsFromEnv {
			name += " [active]"
			activeProfileIndex = i
		}
		s.AddItem(name, "", 0, nil)
	}

	// Set initial selection
	s.SetCurrentItem(activeProfileIndex)

	return nil
}
