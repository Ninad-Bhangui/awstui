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

// HomeScreen represents the main services overview screen
type HomeScreen struct {
	*tview.Table
	onServiceSelect func(service string)
}

// Service represents an AWS service with its details
type Service struct {
	Name        string
	Description string
	Command     string
}

// Available services
var services = []Service{
	{"EC2 Instances", "Manage virtual servers in the cloud", "ec2"},
	{"ECR Repositories", "Manage Docker container images", "ecr"},
	{"Lambda Functions", "Run code without provisioning servers", "lambda"},
	{"Secrets Manager", "Store and manage sensitive information", "secrets"},
}

// NewHomeScreen creates a new home screen
func NewHomeScreen(onServiceSelect func(service string)) *HomeScreen {
	home := &HomeScreen{
		Table:           tview.NewTable().SetSelectable(true, false),
		onServiceSelect: onServiceSelect,
	}

	// Set up table
	home.SetBorder(true)
	home.SetTitle("AWS Services")
	home.SetTitleAlign(tview.AlignLeft)

	// Set up headers
	home.SetCell(0, 0, tview.NewTableCell("Service").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	home.SetCell(0, 1, tview.NewTableCell("Description").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	home.SetCell(0, 2, tview.NewTableCell("Quick Access").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	// Add services
	for i, service := range services {
		home.SetCell(i+1, 0, tview.NewTableCell(service.Name).SetTextColor(tcell.ColorWhite))
		home.SetCell(i+1, 1, tview.NewTableCell(service.Description).SetTextColor(tcell.ColorWhite))
		home.SetCell(i+1, 2, tview.NewTableCell(":"+service.Command).SetTextColor(tcell.ColorGreen))
	}

	// Set up selection handler
	home.SetSelectedFunc(func(row, col int) {
		if row > 0 && row <= len(services) {
			home.onServiceSelect(services[row-1].Command)
		}
	})

	return home
}
