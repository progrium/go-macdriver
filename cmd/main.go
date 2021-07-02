package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version string

	rootCmd = &cobra.Command{
		Version: Version,
		Use:     "macdriver",
		Short:   "FIXME",
	}
)

func init() {
	rootCmd.AddCommand(genCmd)
}

func Execute() {
	fatal(rootCmd.Execute())
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
