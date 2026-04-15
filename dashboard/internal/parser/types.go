package parser

import "time"

type DispatchEntry struct {
	StoryID          string
	Repo             string
	DispatchedAt     time.Time
	Status           string
	CompletionReport string
	ElapsedTime      time.Duration
}

type CompletionReport struct {
	StoryID     string
	Status      string
	Title       string
	Created     time.Time
	LastUpdated time.Time
	ParentEpic  string
	Phase       string
	Repos       []string
	FilePath    string
}

type EscalationReport struct {
	StoryID   string
	Reason    string
	Timestamp time.Time
	FilePath  string
}

type BacklogStory struct {
	StoryID     string
	Title       string
	Status      string
	Repo        string
	DependsOn   string
	Priority    string
	Description string
}

type BacklogOverview struct {
	Done       []BacklogStory
	InProgress []BacklogStory
	Ready      []BacklogStory
	Blocked    []BacklogStory
	Cancelled  []BacklogStory
}

type SessionLog struct {
	FileName  string
	FilePath  string
	Status    string
	Timestamp time.Time
}

type DispatchFile struct {
	StoryID      string
	DispatchedAt time.Time
	DispatchedBy string
	Attempt      int
	MaxRetries   int
	FilePath     string
}
