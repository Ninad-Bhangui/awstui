package ui

import (
	"context"
	"fmt"
	"time"

	awsservices "github.com/Ninad-Bhangui/awstui/aws/services"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ResourceList represents a list of AWS resources
type ResourceList struct {
	*tview.Table
	resourceType string
	cfg          config.Config
}

// Column represents a table column configuration
type Column struct {
	Title string
	Key   string
	Width int
}

// Resource column configurations
var resourceColumns = map[string][]Column{
	"ec2": {
		{"ID", "id", 20},
		{"Name", "name", 30},
		{"Type", "type", 15},
		{"State", "state", 10},
		{"Private IP", "private_ip", 15},
		{"Public IP", "public_ip", 15},
	},
	"ecr": {
		{"Name", "name", 40},
		{"URI", "uri", 60},
		{"Images", "images", 10},
		{"Created", "created", 20},
	},
	"lambda": {
		{"Name", "name", 40},
		{"Runtime", "runtime", 15},
		{"Memory", "memory", 10},
		{"Last Modified", "modified", 20},
	},
	"secrets": {
		{"Name", "name", 40},
		{"Last Modified", "modified", 20},
		{"Days Until Rotation", "rotation", 15},
	},
}

// NewResourceList creates a new resource list
func NewResourceList(resourceType string, cfg config.Config) *ResourceList {
	list := &ResourceList{
		Table:        tview.NewTable().SetSelectable(true, false),
		resourceType: resourceType,
		cfg:          cfg,
	}

	// Set up table
	list.SetBorder(true)
	list.SetTitle(getResourceTitle(resourceType))
	list.SetTitleAlign(tview.AlignLeft)

	// Set up headers
	columns := resourceColumns[resourceType]
	for i, col := range columns {
		cell := tview.NewTableCell(col.Title).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(1)
		if col.Width > 0 {
			cell.SetMaxWidth(col.Width)
		}
		list.SetCell(0, i, cell)
	}

	// Load data
	list.LoadData()

	return list
}

// LoadData loads resource data from AWS
func (l *ResourceList) LoadData() {
	ctx := context.Background()

	switch l.resourceType {
	case "ec2":
		instances, err := awsservices.ListEC2Instances(ctx, l.cfg)
		if err != nil {
			l.showError(err)
			return
		}

		for i, inst := range instances {
			row := i + 1
			l.SetCell(row, 0, tview.NewTableCell(inst.ID))
			l.SetCell(row, 1, tview.NewTableCell(inst.Name))
			l.SetCell(row, 2, tview.NewTableCell(inst.Type))
			stateCell := tview.NewTableCell(inst.State)
			if inst.State == "running" {
				stateCell.SetTextColor(tcell.ColorGreen)
			} else if inst.State == "stopped" {
				stateCell.SetTextColor(tcell.ColorRed)
			}
			l.SetCell(row, 3, stateCell)
			l.SetCell(row, 4, tview.NewTableCell(inst.PrivateIP))
			l.SetCell(row, 5, tview.NewTableCell(inst.PublicIP))
		}

	case "ecr":
		repos, err := awsservices.ListECRRepositories(ctx, l.cfg)
		if err != nil {
			l.showError(err)
			return
		}

		for i, repo := range repos {
			row := i + 1
			l.SetCell(row, 0, tview.NewTableCell(repo.Name))
			l.SetCell(row, 1, tview.NewTableCell(repo.URI))
			l.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%d", repo.ImageCount)))
			l.SetCell(row, 3, tview.NewTableCell(formatTime(repo.CreatedAt)))
		}

	case "lambda":
		functions, err := awsservices.ListLambdaFunctions(ctx, l.cfg)
		if err != nil {
			l.showError(err)
			return
		}

		for i, fn := range functions {
			row := i + 1
			l.SetCell(row, 0, tview.NewTableCell(fn.Name))
			l.SetCell(row, 1, tview.NewTableCell(fn.Runtime))
			l.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%d", fn.MemorySize)))
			l.SetCell(row, 3, tview.NewTableCell(formatTime(fn.LastModified)))
		}

	case "secrets":
		secrets, err := awsservices.ListSecrets(ctx, l.cfg)
		if err != nil {
			l.showError(err)
			return
		}

		for i, secret := range secrets {
			row := i + 1
			l.SetCell(row, 0, tview.NewTableCell(secret.Name))
			l.SetCell(row, 1, tview.NewTableCell(formatTime(secret.LastModified)))
			rotationCell := tview.NewTableCell("-")
			if secret.DaysUntilRotation >= 0 {
				rotationCell.SetText(fmt.Sprintf("%d", secret.DaysUntilRotation))
			}
			l.SetCell(row, 2, rotationCell)
		}
	}
}

func (l *ResourceList) showError(err error) {
	l.Clear()
	l.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("Error: %v", err)).SetTextColor(tcell.ColorRed))
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func getResourceTitle(resourceType string) string {
	switch resourceType {
	case "ec2":
		return "EC2 Instances"
	case "ecr":
		return "ECR Repositories"
	case "lambda":
		return "Lambda Functions"
	case "secrets":
		return "Secrets Manager"
	default:
		return "Resources"
	}
}
