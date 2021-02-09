package workingon

import (
	"errors"
	"fmt"
	"github.com/fefeme/workingon/toggl"
	"github.com/fefeme/workingon/util"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var (
	ErrorPidRequired        = errors.New("no project id found for toggl but TogglPidRequired is set to true")
	ErrorPidNotSetInMapping = errors.New("project pid not set in mapping")
)

func NewTimeEntry(cfg *Config, project string, wid int, summaryOrKey string, templateArgs map[string]string) (*toggl.TimeEntry, error) {
	var timeEntry *toggl.TimeEntry

	// Is this as key for a task in a source?
	task, err := Registry.GetTask(summaryOrKey)

	if task != nil {
		timeEntry = &toggl.TimeEntry{
			Description: fmt.Sprintf("%s: %s", task.Key, task.Summary),
		}
	} else {
		// Maybe it's an alias for a template
		tpl, _ := Configuration.GetTemplate(summaryOrKey)
		if tpl != nil {
			// We need to overwrite the startime and stoptime from the commandline
			// for example: wo add ds 21.01.2021 should work
			timeEntry, err = tpl.CreateTimeEntryFromTemplate(templateArgs)
			if err != nil {
				return nil, err
			}
		} else {
			// It is just a Summary / Description
			timeEntry = &toggl.TimeEntry{
				Description: summaryOrKey,
			}
		}
	}

	if project != "" {
		// Overwrite Project pid with command line project parameter
		pid, err := strconv.Atoi(project)
		if err != nil {
			pm, err := cfg.GetMapping(project)
			if err == nil {
				pid = pm.TogglePid
			}
		}
		timeEntry.Pid = pid
	} else {
		if task != nil {
			pm, _ := Configuration.GetMapping(task.Project.Key)
			if pm != nil {
				if pm.TogglePid == 0 {
					return nil, ErrorPidNotSetInMapping
				}
				timeEntry.Pid = pm.TogglePid
			}
		}
	}
	// If PID is still unknown, let's try to find a pid using a git repository mapping
	if timeEntry.Tid == 0 && timeEntry.Pid == 0 && cfg.Settings.TogglePidRequired {
		//
		pid := FindProjectByGitRepositoryUrl(cfg)
		if pid == 0 {
			return nil, ErrorPidRequired
		}
		timeEntry.Pid = pid
	}

	timeEntry.Wid = wid
	timeEntry.CreatedWith = toggl.CreatedWith

	return timeEntry, nil
}

func setDuration(cfg *Config,
	timeEntry *toggl.TimeEntry, startTime time.Time, stopTime time.Time, duration time.Duration, running bool) error {

	now := time.Now()

	if !running {
		if timeEntry.Start == nil || timeEntry.Start.IsZero() {
			if startTime.IsZero() {
				return errors.New("no start time given")
			}
			timeEntry.Start = &startTime
		}
		if timeEntry.Stop == nil || timeEntry.Stop.IsZero() {
			if duration == 0 {
				if stopTime.IsZero() {
					return errors.New("no stop time or duration given")
				}
				timeEntry.Stop = &stopTime
			}
			timeEntry.Duration = int64(duration.Seconds())
		}

	} else {
		if startTime.IsZero() {
			timeEntry.Start = &now
			timeEntry.Duration = now.Unix() * -1
		} else {
			timeEntry.Start = &startTime
			timeEntry.Duration = startTime.Unix() * -1
		}

	}

	err := timeEntry.Validate()

	return err

}

func AddOrStart(cmd *cobra.Command, cfg *Config,
	wid int, project string, summaryOrKey string,
	startTime time.Time, duration time.Duration,
	templateArgs map[string]string, running bool) (*toggl.TimeEntry, error) {

	timeEntry, err := NewTimeEntry(cfg, project, wid, summaryOrKey, templateArgs)
	if err != nil {
		return nil, fmt.Errorf("timeEntry: %s", err)
	}

	var stopTime time.Time
	s, err := cmd.Flags().GetString("stop")
	if err == nil && s != "" {
		stopTime, err = util.ParseTimeUTCE(s, cfg.Settings.DateLayout, cfg.Settings.DateTimeLayout, &cfg.Settings.Location)
		if err != nil {
			return nil, fmt.Errorf("unable to parse stop time: %s", err)
		}
	}

	err = setDuration(cfg, timeEntry, startTime, stopTime, duration, running)
	if err != nil {
		return nil, err
	}

	dryRun, _ := cmd.Flags().GetBool("dry")

	if !dryRun {
		cl := toggl.NewToggl(cfg.Settings.ToggleApiToken)
		timeEntry, err = cl.TimeEntries.Add(timeEntry)
		if err != nil {
			return nil, err
		}
	}

	return timeEntry, nil
}
