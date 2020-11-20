package cmd

import (
	"fmt"
	"github.com/fefeme/workingon/toggl"
	"github.com/fefeme/workingon/workingon"
	"github.com/spf13/cobra"
)

func NewStopCommand(cfg *workingon.Config) *cobra.Command {
	var stopCommand = &cobra.Command{
		Use:   "stop",
		Short: "Stop currently running timer",
		Long:  `Stop currently running timer`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := toggl.NewToggl(cfg.Settings.ToggleApiToken)

			timeEntry, err := cl.TimeEntries.StopCurrent()
			if err != nil {
				return err
			}

			fmt.Printf("Stopped %s. \n", timeEntry)

			return nil
		},
	}
	return stopCommand
}
