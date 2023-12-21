/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wjrsimpson/prem/fixtures"
)

// blanksCmd represents the blanks command
var blanksCmd = &cobra.Command{
	Use:   "blanks",
	Short: "Print out the gameweeks in which some teams have no fixture",
	Long:  `Print out the gameweeks in which some teams have no fixture`,
	Run: func(cmd *cobra.Command, args []string) {
		fixtures.PrintBlanks()
	},
}

func init() {
	rootCmd.AddCommand(blanksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// blanksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// blanksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
