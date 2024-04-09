package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Check health of the service",
		RunE:  handleHealth,
	}

	rootCmd.AddCommand(healthCmd)
}

func handleHealth(cmd *cobra.Command, args []string) error {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/healthcheck", nil)
	if err != nil {
		fmt.Println("creating request failed")
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("healthcheck failed")
		os.Exit(1)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("healthcheck failed")
		os.Exit(1)
	}

	return nil
}
