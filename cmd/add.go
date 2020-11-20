package cmd

import (
	"fmt"
	"github.com/fefeme/workingon/util"
	"github.com/fefeme/workingon/workingon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//"strconv"
	"time"
)

var (
	UnableToParseArgs = fmt.Errorf("unable to make sense of the arguments")
	DurationRequired  = fmt.Errorf("a duration is required")
	//NoProject         = fmt.Errorf("unable to figure out project")
)

type CommandArgs struct {
	Duration     time.Duration
	StartTime    time.Time
	SummaryOrKey string
}

type TaskDef struct {
	Name      string `mapstructure:"name"`
	TogglTask int64  `mapstructure:"toggl_task"`
	Start     string `mapstructure:"start"`
	Stop      string `mapstructure:"stop"`
	Alias     string `mapstructure:"alias"`
}

// parseArgs tries to parse and validate the optional arguments
// into meaningful values.
func parseArgs(cmd *cobra.Command, running bool, args []string, cfg *workingon.Config) (parsedArgs *CommandArgs, err error) {
	parsedArgs = &CommandArgs{
		Duration:     0,
		StartTime:    time.Time{},
		SummaryOrKey: "",
	}

	appendTime, _ := cmd.Flags().GetBool("append")

	for i, arg := range args {
		var val time.Time
		if appendTime {
			val = util.ParseDate(arg, cfg.Settings.DateLayout, &cfg.Settings.Location)
		} else {
			val = util.ParseTimeUTC(arg, cfg.Settings.DateLayout, cfg.Settings.DateTimeLayout, &cfg.Settings.Location)
		}

		if val.IsZero() {
			//fmt.Println("time not yet set.")
			val = util.ParseDateTimeUTC(arg, cfg.Settings.DateTimeLayout, &cfg.Settings.Location)
			if val.IsZero() {
				//fmt.Println(fmt.Sprintf("parameter %s is not a datetime", arg))
				val, err = util.ParseTimeUTCE(arg, cfg.Settings.DateTimeLayout, cfg.Settings.DateLayout, &cfg.Settings.Location)
				if val.IsZero() {
					//fmt.Println(fmt.Sprintf("parameter %s is not a time", arg))
					valDuration, err := time.ParseDuration(arg)
					if err != nil {
						//fmt.Println(fmt.Sprintf("parameter %s is not a duration", arg))
						if parsedArgs.SummaryOrKey != "" && i+1 == len(args) {
							return nil, UnableToParseArgs
						}
						//fmt.Println(fmt.Sprintf("parameter %s is the summary", arg))
						parsedArgs.SummaryOrKey = arg
					} else {
						//fmt.Println(fmt.Sprintf("parameter %s is a duration", arg))
						parsedArgs.Duration = valDuration
					}
				} else {
					parsedArgs.StartTime = val
				}

			} else {
				//fmt.Println(fmt.Sprintf("parameter %s is a start date and time", arg))
				parsedArgs.StartTime = val
			}
		} else {
			parsedArgs.StartTime = val
		}
	}

	if appendTime && parsedArgs.Duration == 0 && !running {
		return nil, DurationRequired
	}
	return parsedArgs, nil

}

func NewAddCommand(cfg *workingon.Config) *cobra.Command {
	var (
		addCommandArgs *CommandArgs
		addCommand     = &cobra.Command{
			Use:   "add <key|summary|template alias> <start time> <duration>",
			Short: "Add a time entry",
			Long: `Add a time entry

Either from a template set in your config file 
or by description/key, start time and duration`,

			Args: func(cmd *cobra.Command, args []string) error {
				var err error
				addCommandArgs, err = parseArgs(cmd, false, args, cfg)
				if err != nil {
					return err
				}
				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {

				templateArgs, err := cmd.Flags().GetStringToString("templateArgs")
				if err != nil {
					return err
				}

				wid, err := cmd.Flags().GetInt("wid")
				if err != nil {
					return err
				}

				project, err := cmd.Flags().GetString("project")
				if err != nil {
					return err
				}

				timeEntry, err := workingon.AddOrStart(cmd, cfg, wid, project, addCommandArgs.SummaryOrKey, addCommandArgs.StartTime,
					addCommandArgs.Duration, templateArgs, false)
				if err != nil {
					return err
				}
				fmt.Println(timeEntry.Format(cfg.Settings.DateTimeLayout, &cfg.Settings.Location))

				return nil
			},
		}
	)

	// Flags
	addCommand.Flags().StringP("stop", "n", "", "Stop Time")
	addCommand.Flags().StringP("project", "p", viper.GetString("TOGGL_PROJECT"), "Set project")
	addCommand.Flags().BoolP("dry", "d", false, "Do not create anything in toggl")
	addCommand.Flags().BoolP("append", "a", false, "Append to last time entry")
	addCommand.Flags().BoolP("fuzzy", "f", false, "Add some fuzziness to the start and stop time")
	addCommand.Flags().IntP("wid", "w", cfg.Settings.ToggleWid, "Toggle track workspace id")

	addCommand.Flags().StringToStringP("templateArgs", "t", nil, "List of named template args")

	return addCommand
}
