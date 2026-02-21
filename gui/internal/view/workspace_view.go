package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/dvalkoff/gomessenger/gui/internal/events"
)

type WorkspaceView struct {
	Window      fyne.Window
	EventStream chan<- events.Event
}

func (view *WorkspaceView) Run() {
	urlEntry := widget.NewEntry()

	workspaceForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Workspace URL", Widget: urlEntry},
		},
		OnSubmit: func() {
			url := urlEntry.Text
			view.EventStream <- events.WorkspaceSetEvent{URL: url}
		},
	}
	fyne.DoAndWait(func() {
		view.Window.SetContent(workspaceForm)
	})
}
