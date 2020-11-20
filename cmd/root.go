package cmd

import (
	"fmt"
	"github.com/fefeme/workingon/workingon"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "wo",
		Short: "Working on helps you track what you're working on.",
		Long: `                     __   .__                                
__  _  _____________|  | _|__| ____    ____     ____   ____  
\ \/ \/ /  _ \_  __ \  |/ /  |/    \  / ___\   /  _ \ /    \ 
 \     (  <_> )  | \/    <|  |   |  \/ /_/  > (  <_> )   |  \
  \/\_/ \____/|__|  |__|_ \__|___|  /\___  /   \____/|___|  /
                         \/       \//_____/               \/ 

`,
		SilenceUsage: true,
	}
)

func Execute() {
	cfg, err := workingon.InitConfig()

	for _, source := range workingon.Registry.RegisteredSources {
		err := source.Configure(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if err != nil {
		panic(err)
	}
	rootCmd.AddCommand(
		NewAddCommand(cfg),
		NewProjectsCommand(cfg),
		NewStartCommand(cfg),
		NewTasksCommand(cfg),
		NewWhatCommand(cfg),
		NewStopCommand(cfg),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
