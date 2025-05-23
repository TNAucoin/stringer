/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tnaucoin/stringer/parser"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "scan a directory for Github CompositeActions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]
		actions, err := parser.ParseCompositeActions(root)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if len(actions) == 0 {
			fmt.Println("No composite actions found")
			return
		}

		for _, a := range actions {
			fmt.Printf("🔹 %s — %s\n", a.Name, a.Description)
			fmt.Printf("   Inputs: %v\n", a.Inputs)
			fmt.Printf("   Outputs: %v\n\n", a.Outputs)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
