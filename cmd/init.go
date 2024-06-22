/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

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
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
  initCmd.Flags().BoolP("gitignore", "g", false, "Append .envim to .gitignore")
  initCmd.Flags().BoolP("dotnvim", "d", false, "Create .nvim directory")
}
