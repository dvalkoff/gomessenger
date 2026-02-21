package application

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"github.com/dvalkoff/gomessenger/gui/internal/events"
	"github.com/dvalkoff/gomessenger/gui/internal/model"
	"github.com/dvalkoff/gomessenger/gui/internal/view"
)

type AppState struct {
	workspace *model.Workspace
	user      *model.User

	contacts *model.Contacts
	chats    *model.Chats

	currentChat *model.Chat
}

type EventProcessor interface {
	Run()
}

type eventProcessor struct {
	appState    AppState
	eventStream chan events.Event

	workspaceView *view.WorkspaceView
	signInView    *view.SignInView
	signUpView    *view.SignUpView

	workspaceService model.WorkspaceService
	userService      model.UserService

	window fyne.Window
}

func NewEventProcessor(
	eventStream chan events.Event,
	appWindow fyne.Window,

	workspaceView *view.WorkspaceView,
	signInView *view.SignInView,
	signUpView *view.SignUpView,

	workspaceService model.WorkspaceService,
	userService model.UserService,
) EventProcessor {
	return &eventProcessor{
		eventStream: eventStream,
		window:      appWindow,
		appState:    AppState{},

		workspaceView: workspaceView,
		signInView:    signInView,
		signUpView:    signUpView,

		workspaceService: workspaceService,
		userService:      userService,
	}
}

func (sm *eventProcessor) Run() {
	sm.eventStream <- events.InitWorkspaceEvent{}
	sm.processEvents()
}

func (sm *eventProcessor) processEvents() {
	for {
		event := <-sm.eventStream
		slog.Info("Processing event", event)
		switch eventTyped := event.(type) {
		case events.InitWorkspaceEvent:
			sm.workspaceService.InitWorkspace(eventTyped)
			sm.appState.workspace = model.CurrentWorkspace()
		case events.InitUserEvent:
			sm.userService.InitUser(eventTyped)
			sm.appState.user = model.CurrentUser()
		case events.WorkspaceSetEvent:
			sm.workspaceService.SaveWorkspace(eventTyped)
		case events.UserSignUpAttemptedEvent:
			sm.userService.CreateUser(eventTyped)
			sm.appState.user = model.CurrentUser()
		case view.UserSignInAttemptedEvent:
			sm.eventStream <- events.ErrorNotificationEvent{Err: fmt.Errorf("Signing in is not implemented")}

		case events.SwitchViewEvent:
			sm.handleSwitchView(eventTyped)
		case events.ErrorNotificationEvent:
			sm.handleErrorView(eventTyped)
		default:
			slog.Error("Event of unknown type")
		}
	}
}

func (sm *eventProcessor) handleSwitchView(event events.SwitchViewEvent) {
	switch event.NewView {
	case events.WorkspaceView:
		sm.workspaceView.Run()
	case events.SignInView:
		sm.signInView.Run()
	case events.SignUpView:
		sm.signUpView.Run()
	}
}

func (sm *eventProcessor) handleErrorView(event events.ErrorNotificationEvent) {
	errorDialog := view.ErrorView{
		Err:    event.Err,
		Window: sm.window,
	}
	errorDialog.Run()
}
