package cmd

import (
	"fmt"
	"github.com/fefeme/workingon/util"
	"github.com/fefeme/workingon/workingon"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	//isTime = regexp.MustCompile("^\\d{1,2}:\\d{2}$")
	isDay         = regexp.MustCompile("^\\d{1,2}")
	isDayAndMonth = regexp.MustCompile("^(\\d{1,2}).(\\d{2})$")
	isTime        = regexp.MustCompile("^\\d{1,2}:\\d{2}$")
	isTimeRange   = regexp.MustCompile("^(\\d{1,2}:\\d{2})-(\\d{1,2}:\\d{2})$")
)

type ParseArgsConfig struct {
	defaultDateFormat     string
	defaultDateTimeFormat string
	defaultLocation       *time.Location
}

func ParseDateFromArg(date string, cfg *workingon.Config) time.Time {
	m := isDay.MatchString(date)

	if m {
		year, month, day := time.Now().Date()
		day, _ = strconv.Atoi(date)
		return time.Date(year, month, day, 0, 0, 0, 0, &cfg.Settings.Location)
	}

	matches := isDayAndMonth.FindStringSubmatch(date)
	if len(matches) > 0 {
		year, _, _ := time.Now().Date()
		day, _ := strconv.Atoi(matches[0])
		month, _ := strconv.Atoi(matches[1])
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, &cfg.Settings.Location)
	}

	return time.Time{}
}

func tryTime(value string, dateLayout string, dateTimeLayout string, loc *time.Location) (dt *time.Time, err error) {
	var t time.Time
	if isTime.MatchString(value) {
		t = time.Now()
		value = fmt.Sprintf("%s %s", t.Format(dateLayout), value)
		t, err = time.ParseInLocation(dateTimeLayout, value, loc)
		if err != nil {
			return nil, err
		}
		t = t.UTC()
		return &t, nil
	}
	return nil, fmt.Errorf("%s is not a time", value)

}

func tryDate(value string, dateLayout string, loc *time.Location) (*time.Time, error) {
	dt, err := time.ParseInLocation(dateLayout, value, loc)
	if err != nil {
		if value == "yesterday" {
			dt = time.Now().AddDate(0, 0, -1)
			return &dt, nil
		}
		return nil, err
	}
	return &dt, nil
}

func tryDuration(value string, config *ParseArgsConfig) (*time.Time, time.Duration, error) {
	var m = isTimeRange.FindStringSubmatch(value)

	//fmt.Printf("Match duration '%v', '%v' (%v)", value, m, len(m))
	if len(m) > 0 {
		start, err := util.ParseTimeUTCE(m[1], config.defaultDateFormat, config.defaultDateTimeFormat, config.defaultLocation)
		if err != nil {
			return nil, 0, err
		}

		end, err := util.ParseTimeUTCE(m[2], config.defaultDateFormat, config.defaultDateTimeFormat, config.defaultLocation)
		if err != nil {
			return nil, 0, err
		}
		return &start, end.Sub(start), nil
	} else {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return nil, duration, nil
		}
	}
	return nil, 0, nil
}

func ParseArgs(config *ParseArgsConfig, args []string) (start time.Time, duration time.Duration, tail []string) {
	// initialize start time
	start = time.Now()

	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		//fmt.Printf("Checking %v (%d) \n", arg, i)
		var valid = false
		dt, d, err := tryDuration(arg, config)
		if d != 0 {
			duration = d
			if dt != nil {
				start = time.Date(start.Year(), start.Month(), start.Day(), dt.Hour(), dt.Minute(), 0, 0, config.defaultLocation)
			}
			valid = true
			continue
		}
		dt, err = tryDate(arg, config.defaultDateFormat, config.defaultLocation)
		if err == nil {
			start = time.Date(dt.Year(), dt.Month(),dt.Day(), start.Hour(), start.Minute(), 0, 0, config.defaultLocation)
			valid = true
			continue
		}
		dt, err = tryTime(arg, config.defaultDateFormat, config.defaultDateTimeFormat, config.defaultLocation)
		if err == nil {
			start = time.Date(start.Year(), start.Month(), start.Day(), dt.Hour(), dt.Minute(), 0, 0, config.defaultLocation)
			valid = true
			continue
		}
		if !valid {
			tail = append(tail, arg)
			//fmt.Printf("Neither a duration or a time %v -> %v \n", arg, tail)
		}
	}

	return start, duration, tail
}
