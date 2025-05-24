/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tnaucoin/stringer/internal/gitfetcher"
	"github.com/tnaucoin/stringer/internal/store"
	"github.com/tnaucoin/stringer/parser"
	"github.com/tnaucoin/stringer/types"
)

var (
	outputPath string
	cachePath  string
	forceScan  bool
	repo       string
	ref        string
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "scan a directory for Github CompositeActions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var actions []types.CompositeAction
		root := args[0]
		if repo != "" {
			opts := gitfetcher.Options{
				Repo: repo,
				Ref:  ref,
			}
			repoActions, err := gitfetcher.FetchCompositeActionsFromRepo(opts)
			if err != nil {
				fmt.Printf("failed to fetch github repo %s with ref: %s: %v\n", opts.Repo, opts.Ref, err)
				os.Exit(1)
			}
			actions = append(actions, repoActions...)
		} else {
			localFileActions, err := parser.ParseCompositeActions(root)

			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			actions = append(actions, localFileActions...)
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

		if outputPath != "" {
			if err := store.SaveActions(actions, outputPath); err != nil {
				fmt.Println("Failed to write output JSON:", err)
				os.Exit(1)
			}
		} else {
			valid, _ := store.IsCacheValid(root, cachePath)
			if !valid || forceScan {
				if err := store.SaveActionsWithHash(actions, root, cachePath); err != nil {
					fmt.Println("failed to write interal cache:", err)
					os.Exit(1)
				}
				fmt.Println("Updating internal actions cache")
			}
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to write parsed actions to JSON")
	scanCmd.Flags().StringVar(&cachePath, "cache", ".stringercache.json", "Path to store internal action cache")
	scanCmd.Flags().BoolVar(&forceScan, "force", false, "Force cache refresh")
	scanCmd.Flags().StringVar(&repo, "repo", "", "Github repo to scan composite actions from (my-org/my-repo)")
	scanCmd.Flags().StringVar(&ref, "ref", "main", "Git ref to use when scanning a Github repo (e.g. branch, tag")
	rootCmd.AddCommand(scanCmd)
}
