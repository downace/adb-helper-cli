package cmd

import (
	"downace/adb-helper-cli/internal/adb"
	"os"

	"github.com/spf13/cobra"
)

const cmdGroupApp = "app"

var rootCmd = &cobra.Command{
	Use:   "adb-helper",
	Short: "CLI tool to simplify some ADB operations",
}

func init() {
	rootCmd.AddGroup(&cobra.Group{ID: cmdGroupApp})
	rootCmd.PersistentFlags().StringVarP(&adb.Binary, "adb", "a", "adb", "ADB binary path")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
