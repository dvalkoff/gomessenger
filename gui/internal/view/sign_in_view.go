package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dvalkoff/gomessenger/gui/internal/events"
)

type UserSignInAttemptedEvent struct {
	Username string
	Password string
}

type SignInView struct {
	Window      fyne.Window
	EventStream chan<- events.Event
}

func (view *SignInView) Run() {
	signInContainer := view.signInContainer()

	signUpButton := widget.NewButton("Sign up", func() {
		view.EventStream <- events.SwitchViewEvent{NewView: events.SignUpView}
	})

	fyne.DoAndWait(func() {
		view.Window.SetContent(container.NewVBox(
			widget.NewLabel("Sign in"),
			signInContainer,
			container.NewHBox(widget.NewLabel("Don't have an account?"), signUpButton),
		))
	})
}

func (view *SignInView) signInContainer() *fyne.Container {
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	signInForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		},
		OnSubmit: func() {
			username := usernameEntry.Text
			password := passwordEntry.Text
			view.EventStream <- UserSignInAttemptedEvent{
				Username: username,
				Password: password,
			}
		},
	}
	return container.NewVBox(signInForm)
}
