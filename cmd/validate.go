package cmd

import (
	"envim/validate"
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate config file",
	Long: `Validate the structure of envim config file. The file should return a table
  that specific keys according to documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
    file, err := cmd.Flags().GetString("file")
    if err != nil {
      fmt.Println("Error parsing file name")
      return
    }
    if ok, err := validate.Validate(file); ok {
      fmt.Println("All checks passed!")
    } else {
      fmt.Println("File is invalid:\n", err)
    }
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
  validateCmd.Flags().StringP("file", "f", "envim.lua", "name of the config file")
}
