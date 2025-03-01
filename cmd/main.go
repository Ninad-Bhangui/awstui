package main

import (
	"github.com/Ninad-Bhangui/awstui/aws"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	profileList, err := aws.ParseProfilesFromConfig("")
	if err != nil {
		panic(err)
	}
	textView := tview.NewTextView().
		SetText("Hello, AWS TUI!").
		SetText(profileList[0]).
		SetTextAlign(tview.AlignCenter)

	if err := app.SetRoot(textView, true).Run(); err != nil {
		panic(err)
	}
}
