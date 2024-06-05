package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envim",
	Short: "envim is a tool to manage different neovim configurations for different projects.",
  Long: `envim is a tool to manage different neovim configurations for different projects.`,

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

