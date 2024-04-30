package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "needforheat-server-api",
	Short:        "needforheat-server-api is the needforheat server api",
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: printUsage,
}

func Execute() error {
	return rootCmd.Execute()
}

func printUsage(cmd *cobra.Command, args []string) {
	cmd.Print(cmd.UsageString())
}
