/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"envim/run"
	"envim/validate"
	"log"
	"slices"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Open neovim inside the environment",
	Long: `Run command will trigger the start of neovim inside the environment.
  If the environment does not exist in the working directory, it will just be translated as "nvim" command.
  This behaviour can be overwriten by using the -c, --conservative flag.`,
	Run: func(cmd *cobra.Command, args []string) {

    conservative := slices.Contains(args, "--conservative")
		// run nvim with arguments if environment is not initialized
		if !validate.EnvimExists && !conservative {
			run.RunDefault(args)
			return
		}

		if err := validate.ValidateEnvironment(); err != nil {
			log.Fatal(err)
		}
		if err := run.Run("v0.10.0", nil, args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
  runCmd.Flags().Bool("conservative", false, "Don't alias nvim if envim is not initialized")
	runCmd.DisableFlagParsing = true
}
