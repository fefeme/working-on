package toggl

import "time"

type TaskList struct {
	Count int
	Tasks []Task
}

// A toggle task
type Task struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	Pid              int       `json:"pid"`
	Wid              int       `json:"wid"`
	Uid              int       `json:"udi"`
	EstimatedSeconds int       `json:"estimated_seconds"`
	TrackedSeconds   int       `json:"tracked_seconds"`
	Active           bool      `json:"active"`
	At               time.Time `json:"at"`
}

// A toggl time entry
type TimeEntry struct {
	Id          int        `json:"id,omitempty"`
	Description string     `json:"description"`
	Wid         int        `json:"wid"`
	Pid         int        `json:"pid,omitempty"`
	Tid         int        `json:"tid,omitempty"`
	Billable    bool       `json:"billable"`
	Start       *time.Time `json:"start"`
	Stop        *time.Time `json:"stop"`
	Duration    int64      `json:"duration,omitempty"`
	CreatedWith string     `json:"created_with"`
	Tags        []string   `json:"tags,omitempty"`
	Duronly     *bool      `json:"duronly,omitempty"`
	At          *time.Time `json:"at,omitempty"`
}

// A toggl workspace
type Workspace struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Premium bool   `json:"premium"`
}

// A toggl project
type Project struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Wid            int       `json:"wid"`
	Cid            int       `json:"cid"`
	Active         bool      `json:"active" `
	IsPrivate      bool      `json:"is_private"`
	Template       bool      `json:"util"`
	TemplateId     int       `json:"template_id"`
	Billable       bool      `json:"billable"`
	AutoEstimates  bool      `json:"auto_estimates"`
	EstimatedHours int       `json:"estimated_hours"`
	At             time.Time `json:"at"`
	Color          string    `json:"color"`
	Rate           float32   `json:"rate"`
	CreateAt       time.Time `json:"created_at"`
}
