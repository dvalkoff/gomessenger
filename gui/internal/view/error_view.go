package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type ErrorView struct {
	Window fyne.Window
	Err    error
}

func (view *ErrorView) Run() {
	fyne.DoAndWait(func() {
		errorDialog := dialog.NewError(view.Err, view.Window)
		errorDialog.Show()
	})
}
