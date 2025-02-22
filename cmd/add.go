/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

// make a slice to add these tasks to
type Task struct {
	Id     int16     `json :"id"`
	Task   string    `json:"task"`
	Cat    string    `json :"cat"`
	Time   time.Time `json:"time"`
	Isdone bool      `json:"isdone"`
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "This command is used to add the task",
	Long:  `This command is used to add the tasks in the to do list`,
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		cat, _ := cmd.Flags().GetString("cat")

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringP("Task", "t", "", "Provide the task ")
	addCmd.Flags().StringP("cat", "c", "", "Provide the cattogry ")
}
