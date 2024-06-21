package main

import (
	"github.com/sanehe/workday-sign/internal/sign"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:     "sane",
	Short:   "sane: 小工具",
	Long:    `sane: 小工具`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(sign.CmdSign)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
