/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wjrsimpson/prem/fixtures"
)

// nextCmd represents the next command

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Prints out the next 5 fixtures for each team",
	Long: `Prints out the next 5 fixtures for each team, along with the difficulty of each fixture.
	
The fixtures will be retrieved from the FPL API and cached in the user's cache directory. You can force a refresh of the cache by using the -r flag.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fixtures.PrintNextFixtures()
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
