/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wjrsimpson/prem/fixtures"
)

// doublesCmd represents the doubles command
var doublesCmd = &cobra.Command{
	Use:   "doubles",
	Short: "Print out gameweeks where one or more teams play twice",
	Long:  `Print out gameweeks where one or more teams play twice, and list the fixtures`,
	Run: func(cmd *cobra.Command, args []string) {
		fixtures.PrintDoubles()
	},
}

func init() {
	rootCmd.AddCommand(doublesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doublesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doublesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
