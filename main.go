package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var MainCmd = &cobra.Command{
	Use:   "failsave",
	Short: "Utility for saving shell environment before command execution",
}

func init() {
	MainCmd.AddCommand(RunCmd)
}

func Main() error {
	return MainCmd.Execute()
}

func main() {
	err := Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
