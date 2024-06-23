package cmd

import (
	env "envim/environment"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new envim environment",
	Long: `Create a new envim environment in the current directory alognside a config file.
If there is already a config file, the creation of the config will be skipped.
If used with -g, --gitignore flag, it will append the .envim directory to the .gitignore file.
If used with -d, --dotnvim flag, it will create the .nvim folder in the current directory as well.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("Error parsing file flag")
			return
		}

		gitignore, err := cmd.Flags().GetBool("gitignore")
		if err != nil {
			fmt.Println("Error parsing gitignore flag")
			return
		}

		dotnvim, err := cmd.Flags().GetBool("dotnvim")
		if err != nil {
			fmt.Println("Error parsing dotnvim flag")
			return
		}

		validationErrs := env.ValidateCreation(map[string]env.FlagData{
			"dotnvim":   {Active: dotnvim, Value: ""},
			"gitignore": {Active: gitignore, Value: ""},
			"file":      {Active: true, Value: configFile},
		})

		if len(validationErrs) > 0 {
			for _, err := range validationErrs {
				log.Println(err)
			}
			os.Exit(1)
		}

		if err := env.CreateEnvironment(); err != nil {
			log.Fatal(err)
		}

		if dotnvim {
			if err := env.CreateDotNvim(); err != nil {
				log.Fatal(err)
			}
		}

		if gitignore {
			if err := env.AppendToGitignore(); err != nil {
				log.Fatal(err)
			}
		}

		if err := env.CreateConfigFile(configFile); err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("gitignore", "g", false, "Append .envim to .gitignore")
	initCmd.Flags().BoolP("dotnvim", "d", false, "Create .nvim directory")
}
