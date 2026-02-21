package events

type Event any

type BackendEvent struct {
}

type ViewEnum int

const (
	WorkspaceView = iota
	SignInView
	SignUpView
	MainAppView
)

type SwitchViewEvent struct {
	NewView ViewEnum
}

type ErrorNotificationEvent struct {
	Err error
}

type InitWorkspaceEvent struct{}

type InitUserEvent struct{}

type WorkspaceSetEvent struct {
	URL string
}

type UserSignUpAttemptedEvent struct {
	Nickname string
	Password string
}
