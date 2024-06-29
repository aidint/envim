package cmd

import (
	"encoding/json"
	"envim/initialize"
	"envim/install"
	"log"
	"os"
	"path"
	"github.com/spf13/cobra"
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
    //Validate evironment
    if err := initialize.ValidateEnvironment(); err != nil {
      log.Fatal(err)
    }

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal("Error parsing file name")
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatal("Error parsing force flag")
		}

    m, err := install.Install(file, force)
    if err != nil {
      log.Fatal(err)
    }
    
    ps, err := json.MarshalIndent(m, "", "  ")
    if err := os.WriteFile(path.Join(".envim", "installed.json"), ps, 0644); err != nil {
      log.Fatal(err)
    }
    log.Printf("Installed dependencies: \n%s", string(ps))
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().Bool("force", false, "Force install")
}
