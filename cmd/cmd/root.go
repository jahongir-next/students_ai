package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "ChDPU AI Monitor CLI",
}

// Execute barcha buyruqlarni ishga tushiradi
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(ingestCmd)
	rootCmd.AddCommand(fetchGradesCmd)
}
