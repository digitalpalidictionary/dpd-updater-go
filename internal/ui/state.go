package ui

import (
	"sync"

	"github.com/digitalpalidictionary/dpd-updater-go/internal/config"
	"github.com/digitalpalidictionary/dpd-updater-go/internal/github"
)

type AppState struct {
	sync.RWMutex
	Config         *config.Config
	LatestRelease  *github.ReleaseInfo
	IsUpdateAvail  bool
	IsProcessing   bool
	StatusMessage  string
	Progress       float64
	Logs           []string
}

func NewAppState(cfg *config.Config) *AppState {
	return &AppState{
		Config: cfg,
	}
}

func (s *AppState) SetStatus(msg string, progress float64) {
	s.Lock()
	defer s.Unlock()
	s.StatusMessage = msg
	s.Progress = progress
	if msg != "" {
		s.Logs = append(s.Logs, msg)
	}
}

func (s *AppState) AddLog(msg string) {
	s.Lock()
	defer s.Unlock()
	s.Logs = append(s.Logs, msg)
}

func (s *AppState) SetProcessing(val bool) {
	s.Lock()
	defer s.Unlock()
	s.IsProcessing = val
}
