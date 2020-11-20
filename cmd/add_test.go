package cmd

import (
	"fmt"
	"github.com/fefeme/workingon/workingon"
	"testing"
	"time"
)

const (
	dateLayout     = "2.1.2006"
	dateTimeLayout = "2.1.2006 15:04"
)

var (
	loc   *time.Location
	tasks []workingon.Task
)

func init() {
	var err error
	if loc, err = time.LoadLocation("Europe/Berlin"); err != nil {
		panic(err)
	}
	tasks = []workingon.Task{
		{
			Key:     "TEST-1",
			Summary: "this is a test task",
			Project: workingon.Project{
				Key:  "TEST-PROJECT",
				Name: "A test project",
			},
			TogglTask: 100,
		},
	}

}

func TestParseArgsInvalid(t *testing.T) {
	argsToParse := []string{"This", "Makes", "No", "Sense"}

	cfg := workingon.Config{
		CreatedWith: "test",
		Settings:    workingon.Settings{},
		Projects:    nil,
		Templates:   nil,
		Sources:     nil,
	}

	cmd := NewAddCommand(&cfg)

	_, err := parseArgs(cmd, argsToParse, &cfg)
	if err != UnableToParseArgs {
		t.Errorf("expected error UnableToParseArgs")
	}
}

func TestParseArgsValidSummaryAndStartAndDuration(t *testing.T) {

	const (
		summary   = "A Task that started at 10 and took 1h"
		startTime = "10:00"
	)
	var (
		today = time.Now().Format(dateLayout)
	)
	today += " " + startTime

	cfg := workingon.Config{
		CreatedWith: "test",
		Settings:    workingon.Settings{},
		Projects:    nil,
		Templates:   nil,
		Sources:     nil,
	}

	cmd := NewAddCommand(&cfg)

	expTime, err := time.ParseInLocation(dateTimeLayout, today, loc)
	if err != nil {
		t.Error(err)
	}
	argsToParse := []string{"10:00", summary, "1h"}
	parsedArgs, err := parseArgs(cmd, argsToParse, &cfg)
	if err != nil {
		t.Error(err)
	}
	if parsedArgs.SummaryOrKey != summary {
		t.Errorf("expected summary to be %s, received %s", summary, parsedArgs.SummaryOrKey)
	}
	if parsedArgs.StartTime.Unix() != expTime.Unix() {
		t.Errorf("expected start time to be %s, received %s", expTime, parsedArgs.StartTime)
	}
	fmt.Println(parsedArgs)
}

func TestParseArgsTemplateAlias(t *testing.T) {
	/*
	var (
		today = time.Now().Format(dateLayout)
	)
	 */

	cfg := workingon.Config{
		CreatedWith: "test",
		Settings:    workingon.Settings{},
		Projects:    nil,
		Templates: []workingon.TemplateConfig{{
			Alias:       "ds",
			Description: "Daily Standup",
			Start:       "16:30",
			Stop:        "16:45",
			TogglTask:   40819208,
		}},
		Sources: nil,
	}

	cmd := NewAddCommand(&cfg)

	/*
	expTime, err := time.ParseInLocation(dateTimeLayout, today, loc)
	if err != nil {
		t.Error(err)
	}*/

	argsToParse := []string{"ds", "29.10.2020"}
	parsedArgs, err := parseArgs(cmd, argsToParse, &cfg)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(parsedArgs)
}
