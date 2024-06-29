package cmd

import (
	"envim/initialize"
	"github.com/spf13/cobra"
	"log"
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
			log.Fatal("Error parsing file flag")
		}

		if err != nil {
			log.Fatal("Error parsing gitignore flag")
		}

		dotnvim, err := cmd.Flags().GetBool("dotnvim")
		if err != nil {
			log.Fatal("Error parsing dotnvim flag")
		}

		if filePath, err := initialize.CreateEnvironment(); err != nil {
			log.Fatal(err)
		} else {
      log.Println("Environment created in", filePath)
    }

		if dotnvim {
			if filePath, err := initialize.CreateDotNvim(); err != nil {
				log.Println(err)
			} else {
        log.Println(".nvim folder created in", filePath)
      }
		}

		if filePath, err := initialize.CreateConfigFile(configFile); err != nil {
			log.Println(err)
		} else {
      log.Println("Config file created in", filePath)
    }

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("dotnvim", "d", false, "Create .nvim directory")
}
