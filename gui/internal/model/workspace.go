package model

import (
	"sync"

	"github.com/dvalkoff/gomessenger/gui/internal/events"
	"github.com/dvalkoff/gomessenger/gui/internal/integration/repository"
)

type Workspace struct {
	Id  int
	URL string
}

var currentWorkspace *Workspace
var workspaceLock *sync.RWMutex = &sync.RWMutex{}

func CurrentWorkspace() *Workspace {
	workspaceLock.RLock()
	defer workspaceLock.RUnlock()

	return currentWorkspace
}

func setCurrentWorkspace(newWorkspace *Workspace) {
	workspaceLock.Lock()
	defer workspaceLock.Unlock()

	currentWorkspace = newWorkspace
}

type WorkspaceService interface {
	InitWorkspace(events.InitWorkspaceEvent)
	SaveWorkspace(events.WorkspaceSetEvent)
}

type workspaceService struct {
	workspaceRepository repository.WorkspaceRepository
	eventStream         chan<- events.Event
}

func NewWorkspaceService(
	workspaceRepository repository.WorkspaceRepository,
	eventStream chan<- events.Event,
) WorkspaceService {
	return &workspaceService{
		workspaceRepository: workspaceRepository,
		eventStream:         eventStream,
	}
}

func (s *workspaceService) InitWorkspace(event events.InitWorkspaceEvent) {
	spaceRow, err := s.workspaceRepository.GetCurrentSpace()
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	if spaceRow == nil {
		s.eventStream <- events.SwitchViewEvent{NewView: events.WorkspaceView}
		return
	}
	newWorkspace := &Workspace{
		Id:  spaceRow.Id,
		URL: spaceRow.URL,
	}
	setCurrentWorkspace(newWorkspace)
	s.eventStream <- events.InitUserEvent{}
}

func (s *workspaceService) SaveWorkspace(event events.WorkspaceSetEvent) {
	row := repository.WorkspaceRow{
		URL: event.URL,
	}
	workspace, err := s.workspaceRepository.AddWorkspace(row)
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	err = s.workspaceRepository.SetWorkspaceCurrent(workspace.Id)
	if err != nil {
		s.eventStream <- events.ErrorNotificationEvent{Err: err}
		return
	}
	newWorkspace := &Workspace{
		Id:  workspace.Id,
		URL: workspace.URL,
	}
	setCurrentWorkspace(newWorkspace)
	s.eventStream <- events.SwitchViewEvent{NewView: events.SignUpView}
}
