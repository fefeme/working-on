package cmd

import (
	"fmt"
	"github.com/fefeme/workingon/workingon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewStartCommand(cfg *workingon.Config) *cobra.Command {
	var (
		appendTo bool
		dry      bool
		cont     bool
		project  string
		commandArgs *CommandArgs
	)

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Start working on a task",
		Long:  `Start working on a task`,
		Args: func(cmd *cobra.Command, args []string) error {
			var err error
			commandArgs, err = parseArgs(cmd, false, args, cfg)
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

			timeEntry, err := workingon.AddOrStart(cmd, cfg, wid, project, commandArgs.SummaryOrKey, commandArgs.StartTime,
				commandArgs.Duration, templateArgs, true)
			if err != nil {
				return err
			}
			fmt.Println(timeEntry.Format(cfg.Settings.DateTimeLayout, &cfg.Settings.Location))

			return nil

		},
	}
	startCommand.Flags().BoolVarP(&appendTo, "append", "a", false, "Use stop time of last time entry as start time for this task")
	startCommand.Flags().BoolVarP(&cont, "continue", "c", false, "Continue last task")
	startCommand.Flags().BoolVarP(&dry, "dry", "d", false, "Do not create anything in toggl")
	startCommand.Flags().StringVarP(&project, "project", "p", viper.GetString("TOGGL_PROJECT"), "Set project")
	startCommand.Flags().StringToStringP("templateArgs", "t", nil, "List of named template args")
	startCommand.Flags().IntP("wid", "w", cfg.Settings.ToggleWid, "Toggle track workspace id")


	return startCommand
}
