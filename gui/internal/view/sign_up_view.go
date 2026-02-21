package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dvalkoff/gomessenger/gui/internal/events"
)

type SignUpView struct {
	Window      fyne.Window
	EventStream chan<- events.Event
}

func (view *SignUpView) Run() {
	signUpContainer := view.signUpContainer()

	signInButton := widget.NewButton("Sign in", func() {
		view.EventStream <- events.SwitchViewEvent{NewView: events.SignInView}
	})
	fyne.DoAndWait(func() {
		view.Window.SetContent(container.NewVBox(
			widget.NewLabel("Sign up"),
			signUpContainer,
			container.NewHBox(widget.NewLabel("Already have an account?"), signInButton),
		))
	})
}

func (view *SignUpView) signUpContainer() *fyne.Container {
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	signUpForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		},
		OnSubmit: func() {
			username := usernameEntry.Text
			password := passwordEntry.Text
			view.EventStream <- events.UserSignUpAttemptedEvent{
				Nickname: username,
				Password: password,
			}
		},
	}
	return container.NewVBox(signUpForm)
}
