package cmd

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/fefeme/workingon/toggl"
	"github.com/fefeme/workingon/workingon"
	"github.com/spf13/cobra"
	"time"
)

func NewWhatCommand(cfg *workingon.Config) *cobra.Command {
	var flagToday bool

	whatCommand := &cobra.Command{
		Use:   "what",
		Short: "What are you working on?",
		Long:  `What are you working on?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := toggl.NewToggl(cfg.Settings.ToggleApiToken)

			if flagToday {
				year, month, day := time.Now().Date()
				start := time.Date(year, month, day, 0, 0, 0, 0, &cfg.Settings.Location)
				now := time.Now()
				timeEntries, err := cl.TimeEntries.List(&start, &now)
				if err != nil {
					return err
				}

				table := simpletable.New()

				table.Header = &simpletable.Header{
					Cells: []*simpletable.Cell{
						{Align: simpletable.AlignLeft, Text: "Description"},
						{Align: simpletable.AlignLeft, Text: "Start"},
						{Align: simpletable.AlignLeft, Text: "Duration"},
					},
				}


				for _, te := range timeEntries.TimeEntries {
					duration := time.Duration(te.Duration) * time.Second
					if duration < 0 {
						duration = time.Now().Sub(*te.Start)
					}
					r := []*simpletable.Cell{
						{Align: simpletable.AlignLeft, Text: te.Description},
						{Align: simpletable.AlignLeft, Text: te.Start.In(&cfg.Settings.Location).Format(cfg.Settings.DateTimeLayout)},
						{Align: simpletable.AlignLeft, Text: duration.String()},
					}

					table.Body.Cells = append(table.Body.Cells, r)

				}

				fmt.Print(table.String())

				return nil
			}

			current, err := cl.TimeEntries.Current()
			if err != nil {
				return err
			}

			var msg string
			prompt, _ := cmd.Flags().GetBool("prompt")
			if prompt {
				if current != nil {
					msg = "\u23f2 "
				} else {
					msg = "\u23f2 "
				}
			} else {
				if current != nil {
					msg = fmt.Sprintf("You are currently working on: '%s'", current.Format(cfg.Settings.DateTimeLayout, &cfg.Settings.Location))
				} else {
					msg = "You are slacking off. Go back to work!"
				}
			}

			fmt.Println(msg)

			return nil

		},
	}
	whatCommand.Flags().BoolP("prompt", "p", false, "Output an indicator for usage in a shell prompt")
	whatCommand.Flags().BoolVarP(&flagToday, "today", "t", false, "List time entries for today")

	return whatCommand
}
