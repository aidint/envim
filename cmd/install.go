package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
  "envim/install"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs all the dependencies of the project",
	Long: `Installs all the dependencies described in envim configuration file (i.e. envim.lua).
The command will install all the packages locally for the project. For neovim installation,
it will not use the local enviornment, rather it will install it in a central location (i.e. ~/.envim/neovim).
The command will skip neovim installation if it is already installed in the central location. This behaviour
can be overwriten by using the --force flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("Error parsing file name")
			return
		}
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Println("Error parsing force flag")
			return
		}

    installed, err := install.Install(file, force)
    if err != nil {
      fmt.Println(err)
      return
    }
    if !installed {
      fmt.Println("Installation failed")
    } else {
      fmt.Println("Installation successful")
    }
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().Bool("force", false, "Force install")
}
