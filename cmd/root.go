package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "twomes-backoffice-api",
	Short:        "twomes-backoffice-api is the twomes backoffice API server",
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
