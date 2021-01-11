package cmd

import (
	. "gopkg.in/check.v1"
	"strings"
	"testing"
	"time"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ParseArgsSuite struct {
	today          time.Time
	yesterday      time.Time
	parseArgConfig ParseArgsConfig
	location       *time.Location
}

var _ = Suite(&ParseArgsSuite{})

func (s *ParseArgsSuite) SetUpTest(c *C) {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	s.location = location
	s.today = time.Now()
	s.yesterday = s.today.AddDate(0, 0, -1)
	s.parseArgConfig = ParseArgsConfig{
		defaultDateFormat:     "2.1.2006",
		defaultDateTimeFormat: "2.1.2006 15:04",
		defaultLocation:       location,
	}

}

func (s *ParseArgsSuite) TestParseArgsYesterday(c *C) {
	args := strings.Split("yesterday 15:30-17:30", " ")
	dt, d, tail := ParseArgs(&s.parseArgConfig, args)

	c.Assert(dt.Month(), Equals, s.yesterday.Month())
	c.Assert(dt.Day(), Equals, s.yesterday.Day())
	c.Assert(dt.Year(), Equals, s.yesterday.Year())
	// Note: Times are in UTC
	c.Assert(dt.Hour(), Equals, 14)
	c.Assert(dt.Minute(), Equals, 30)

	c.Assert(d, Equals, time.Hour * 2)
	c.Assert(tail, HasLen, 0)
}

func (s *ParseArgsSuite) TestParseArgsWithKeyAndDateAndTimeRange(c *C) {
	args := strings.Split("4.1.2021 11:30-12:30 A-KEY", " ")

	dt, d, tail := ParseArgs(&s.parseArgConfig, args)

	c.Assert(dt.Month(), Equals,time.January)
	c.Assert(dt.Day(), Equals, 4)
	c.Assert(dt.Year(), Equals, 2021)
	c.Assert(dt.Hour(), Equals, 10)
	c.Assert(dt.Minute(), Equals, 30)

	c.Assert(d, Equals, time.Hour * 1)
	c.Assert(tail, HasLen, 1)
	c.Assert(tail[0], Equals, "A-KEY")

}

func (s *ParseArgsSuite) TestParseArgsWithKeyAndStartTime(c *C) {
	args := strings.Split("4.1.2021 A-KEY 8:00", " ")
	dt, d, tail := ParseArgs(&s.parseArgConfig, args)

	c.Assert(dt.Month(), Equals,time.January)
	c.Assert(dt.Day(), Equals, 4)
	c.Assert(dt.Year(), Equals, 2021)
	c.Assert(dt.Hour(), Equals, 7)
	c.Assert(dt.Minute(), Equals, 0)
	c.Assert(d, Equals, 0 * time.Second)
	c.Assert(tail, HasLen, 1)
	c.Assert(tail[0], Equals, "A-KEY")

}

func (s *ParseArgsSuite) TestParseArgsWithKeyStartTimeAndDuration(c *C) {
	args := strings.Split("4.1.2021 15:30 2h", " ")
	dt, d, tail := ParseArgs(&s.parseArgConfig, args)

	c.Assert(dt.Month(), Equals,time.January)
	c.Assert(dt.Day(), Equals, 4)
	c.Assert(dt.Year(), Equals, 2021)
	c.Assert(dt.Hour(), Equals, 14)
	c.Assert(dt.Minute(), Equals, 30)
	c.Assert(d, Equals, 2 * time.Hour)
	c.Assert(tail, HasLen, 0)
}

func (s *ParseArgsSuite) TestParseArgsWithKeyAndTimeRange(c *C) {
	args := strings.Split("A-KEY 13:30-17:00", " ")

	dt, d, tail := ParseArgs(&s.parseArgConfig, args)

	c.Assert(dt.Month(), Equals,time.January)
	c.Assert(dt.Day(), Equals, s.today.Day())
	c.Assert(dt.Year(), Equals, s.today.Year())
	c.Assert(dt.Hour(), Equals, 12)
	c.Assert(dt.Minute(), Equals, 30)
	c.Assert(d, Equals, 210 * time.Minute)
	c.Assert(tail, HasLen, 1)
	c.Assert(tail[0], Equals, "A-KEY")
}


func (s *ParseArgsSuite) TestParseArgsWithSummary(c *C) {
	args := strings.Split("4.1.2021 15:30 2h This is a summary", " ")
	dt, d, tail := ParseArgs(&s.parseArgConfig, args)

	c.Assert(dt.Month(), Equals,time.January)
	c.Assert(dt.Day(), Equals, 4)
	c.Assert(dt.Year(), Equals, 2021)
	c.Assert(dt.Hour(), Equals, 14)
	c.Assert(dt.Minute(), Equals, 30)
	c.Assert(d, Equals, 2 * time.Hour)
	c.Assert(tail, HasLen, 4)
	c.Assert(strings.Join(tail,  " "), Equals, "This is a summary")
}

/*7
func TestParseArgs(t *testing.T) {


	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}

	var config = ParseArgsConfig{
		defaultDateFormat:     "2.1.2006",
		defaultDateTimeFormat: "2.1.2006 15:04",
		defaultLocation:       location,
	}

	var (
		args []string
		dt   time.Time
		d    time.Duration
		tail []string
	)

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)


	args = strings.Split("yesterday 15:30-17:30", " ")
	dt, d, tail = ParseArgs(&config, args)

	if !(dt.Month() == today.Month() && dt.Day() == yesterday.Day() && dt.Year() == yesterday.Year() && dt.Hour() == 15 && dt.Minute() == 30) {
		t.Errorf("Expected %v as start time, but got %v", yesterday, dt)
	}

	if d != time.Hour*2 {
		t.Fail()
	}
	if dt.Truncate(24*time.Hour) != time.Now().AddDate(0, 0, -1).Truncate(24*time.Hour) {
		t.Errorf("Expected yesterdays date, but got %v", dt)
	}
	if len(tail) > 0 {
		t.Errorf("length of tail should be 0, but is %d", len(tail))
	}

	args = strings.Split("4.1.2021 15:30-17:30 A-KEY", " ")
	dt, d, tail = ParseArgs(&config, args)
	if !(dt.Month() == 1 && dt.Day() == 4 && dt.Year() == 2012 && dt.Hour() == 14 && dt.Minute() == 30) {
		t.Errorf("Expected 4.1.2021, 15:30 as start time, but got %v", dt)
	}
	if len(tail) != 1 {
		t.Errorf("length of tail should be 1, but is %d", len(tail))
	}
	if tail[0] != "A-KEY" {
		t.Errorf("Should be A-KEY, but was %v", tail[0])
	}

	args = strings.Split("4.1.2021 MOET-297 8:00", " ")
	dt, d, tail = ParseArgs(&config, args)
	if !(dt.Month() == 1 && dt.Day() == 4 && dt.Year() == 2012 && dt.Hour() == 14 && dt.Minute() == 30) {
		t.Errorf("Expected 4.1.2021, 15:30 as start time, but got %v", dt)
	}

	args = strings.Split("4.1.2021 15:30 2h", " ")
	dt, d, tail = ParseArgs(&config, args)
	if !(dt.Month() == 1 && dt.Day() == 4 && dt.Year() == 2012 && dt.Hour() == 14 && dt.Minute() == 30) {
		t.Errorf("Expected 4.1.2021, 15:30 as start time, but got %v", dt)
	}

	args = strings.Split("MOET-297 13:30-17:00", " ")
	dt, d, tail = ParseArgs(&config, args)
	if !(dt.Month() == today.Month() && dt.Year() == today.Year() && dt.Day() == today.Day()) {
		t.Errorf("Start date not correct")
	}



}
*/
