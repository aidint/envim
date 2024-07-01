package cmd

import (
	"encoding/json"
	"envim/config"
	"envim/install"
	"envim/validate"
	"log"
	"os"
	"path"
	"github.com/spf13/cobra"
	lua "github.com/yuin/gopher-lua"
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
		if err := validate.ValidateEnvironment(); err != nil {
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

		L := lua.NewState()
    defer L.Close()

		configMap, err := config.ReadConfig(L, file)
		if err != nil {
			log.Fatal(err)
		}

    bytes, err := json.MarshalIndent(configMap, "", "  ")
    if err != nil {
      log.Fatal("Error marshalling config map")
    }

    err = os.WriteFile(path.Join(".envim", "config.json"), bytes, 0644)

		installMap, err := install.Install(L, configMap, force)
		if err != nil {
			log.Fatal(err)
		}

		ps, err := json.MarshalIndent(installMap, "", "  ")
		if err := os.WriteFile(path.Join(".envim", "envim.json"), ps, 0644); err != nil {
			log.Fatal(err)
		}
		log.Printf("Installed dependencies: \n%s", string(ps))
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().Bool("force", false, "Force install")
}
