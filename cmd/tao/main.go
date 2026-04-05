package main

import (
	"fmt"
	"os"

	"tao-agent/internal/workflow"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tao",
	Short: "Tao Agent is a lightweight configurable agent framework",
}

var runCmd = &cobra.Command{
	Use:   "run [workflow]",
	Short: "Run a workflow",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		w, err := workflow.Load(path)
		if err != nil {
			fmt.Printf("Error loading workflow: %v\n", err)
			os.Exit(1)
		}

		// Initialize engine (agents will be resolved during execution)
		engine := workflow.NewEngine()

		fmt.Printf("Executing workflow: %s\n", w.Name)
		if err := engine.Execute(w); err != nil {
			fmt.Printf("Execution failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
