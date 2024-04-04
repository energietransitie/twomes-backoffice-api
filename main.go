package main

import (
	"github.com/energietransitie/twomes-backoffice-api/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.Execute())
}
