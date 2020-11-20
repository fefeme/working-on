package toggl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fefeme/workingon/util"
	"math"
	"math/rand"
	"net/url"
	"time"
)

const Endpoint = "time_entries"
const CreatedWith = "working_on"

type TimeEntryList struct {
	Count       int
	TimeEntries []TimeEntry
}

type TimeEntryRequest struct {
	TimeEntry *TimeEntry `json:"time_entry"`
}

type TimeEntryResponse struct {
	Data *TimeEntry `json:"data,omitempty"`
}

type TimeEntryClient interface {
	List(*TimeEntryList, error)
	String() string
}

type TimeEntries struct {
	client *Client
}

func (t *TimeEntry) String() string {
	return fmt.Sprintf("%s (%d) from %s to %s for %s", t.Description, t.Pid, t.Start,
		t.Stop, time.Duration(t.Duration)*time.Second)
}

func (t *TimeEntry) Format(dfLayout string, loc *time.Location) string {
	start := ""
	if t.Start != nil {
		start = t.Start.In(loc).Format(dfLayout)
	}

	d := t.Duration
	if d < 0 {
		return fmt.Sprintf("\"%s\" %d at %s", t.Description, t.Pid, start)
	}

	return fmt.Sprintf("\"%s\" %s for %s (%d)", t.Description, start,
		time.Duration(d)*time.Second, t.Pid)
}

func (t *TimeEntry) Fuzz() {
	sig := [2]int{-1, 1}[rand.Intn(2)]
	fuzzyTime := t.Start.Add(time.Duration(rand.Intn(180)*sig) * time.Second)

	t.Start = &fuzzyTime
}

func (t *TimeEntry) Validate() error {
	if t.Duration == 0 {
		if t.Start == nil || t.Start.IsZero() {
			return fmt.Errorf("no start time given, unable to calculate duration")
		}
		if t.Stop == nil {
			return fmt.Errorf("no stop time given, unable to calculate duration")
		}
		duration := t.Stop.Sub(*t.Start)
		if math.Abs(duration.Hours()) > 999 {
			return fmt.Errorf("something went wrong - duration is more than 999 hours")
		}
		t.Duration = duration.Milliseconds() / 1000
	}
	if t.Start == nil || t.Start.IsZero() {
		return fmt.Errorf("something went wrong - no start time set")
	}

	t.Stop = util.TimeInUTC(t.Stop)
	t.Start = util.TimeInUTC(t.Start)

	return nil
}

func (t *TimeEntries) Start(timeEntry *TimeEntry) (*TimeEntry, error) {
	return t.Add(timeEntry)
}

func (t *TimeEntries) Add(timeEntry *TimeEntry) (*TimeEntry, error) {
	message, err := t.client.NewMessage("POST", Endpoint, &TimeEntryRequest{TimeEntry: timeEntry})

	if err != nil {
		return nil, err
	}

	data, err := t.client.SendRequest(message)
	if err != nil {
		return nil, err
	}

	var res TimeEntryResponse

	err = json.Unmarshal(*data, &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (t *TimeEntries) List(start *time.Time, end *time.Time) (*TimeEntryList, error) {

	base, _ := url.Parse(Endpoint)

	params := url.Values{}

	if start != nil {
		params.Add("start_date", start.Format(time.RFC3339))
		params.Add("end_date", end.Format(time.RFC3339))
		base.RawQuery = params.Encode()
	}

	message, err := t.client.NewMessage("GET", base.String(), nil)
	if err != nil {
		return nil, err
	}

	data, err := t.client.SendRequest(message)

	if err != nil {
		return nil, err
	}

	var timeEntries []TimeEntry
	err = json.Unmarshal(*data, &timeEntries)

	if err != nil {
		return nil, err
	}

	return &TimeEntryList{
		TimeEntries: timeEntries,
		Count:       len(timeEntries),
	}, nil
}

func (t *TimeEntries) Current() (*TimeEntry, error) {
	message, err := t.client.NewMessage("GET", fmt.Sprintf("%s/current", Endpoint), nil)
	if err != nil {
		return nil, err
	}

	data, err := t.client.SendRequest(message)
	if err != nil {
		return nil, err
	}

	var res TimeEntryResponse
	err = json.Unmarshal(*data, &res)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}
	return res.Data, nil
}

func (t *TimeEntries) MostRecent() (*TimeEntry, error) {
	timeEntries, err := t.List(nil, nil)
	if err != nil {
		return nil, err
	}
	if len(timeEntries.TimeEntries) > 0 {
		return &timeEntries.TimeEntries[len(timeEntries.TimeEntries)-1], nil
	}
	return nil, nil
}

func (t *TimeEntries) StopCurrent() (*TimeEntry, error) {
	timeEntry, err := t.Current()
	if err != nil {
		return nil, err
	}

	if timeEntry == nil {
		return nil, errors.New("no time entry is currently running")
	}

	message, err := t.client.NewMessage("PUT", fmt.Sprintf("%s/%d/stop", Endpoint, timeEntry.Id), nil)
	if err != nil {
		return nil, err
	}

	data, err := t.client.SendRequest(message)

	if err != nil {
		return nil, err
	}

	var res TimeEntryResponse
	err = json.Unmarshal(*data, &res)

	if err != nil {
		return nil, err
	}

	return res.Data, nil

}
