/*
Copyright © 2025 Prince Kumar
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update <task_id> --cat <new_category> --isdone <true/false> --task <new_task_description>",
	Short: "Update a task in the data/.pk.json file.",
	Long: `Update the "Cat", "isdone", and "task" fields of a task based on its ID in the data/.pk.json file.

Example:
  pk update 1 --cat "Work" --isdone true --task "Finish the report"
  pk update 2 --task "Grocery Shopping"
  pk update 3 --isdone false
`,
	Args: cobra.ExactArgs(1), // Expect exactly one argument: the task ID
	Run: func(cmd *cobra.Command, args []string) {
		taskIDStr := args[0]
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			fmt.Println("Invalid task ID.  Must be an integer:", err)
			return
		}

		newCategory, err := cmd.Flags().GetString("cat")
		if err != nil {
			fmt.Println("Error getting category flag:", err)
			return
		}

		isDoneStr, err := cmd.Flags().GetString("isdone")
		if err != nil {
			fmt.Println("Error getting isdone flag:", err)
			return
		}

		var newIsDone *bool // Use a pointer to handle cases where isdone is not set
		if isDoneStr != "" {
			isDone, err := strconv.ParseBool(isDoneStr)
			if err != nil {
				fmt.Println("Invalid value for 'isdone'. Must be 'true' or 'false':", err)
				return
			}
			newIsDone = &isDone
		}

		newTask, err := cmd.Flags().GetString("task")
		if err != nil {
			fmt.Println("Error getting task flag:", err)
			return
		}

		err = updateTask(taskID, newCategory, newIsDone, newTask)
		if err != nil {
			fmt.Println("Error updating task:", err)
			return
		}

		fmt.Println("Task", taskID, "updated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Define flags for updating the task fields
	updateCmd.Flags().String("cat", "", "The new category for the task.")
	updateCmd.Flags().String("isdone", "", "Set to 'true' or 'false' to update the completion status of the task.") // Allow empty value if user doesn't want to update.
	updateCmd.Flags().String("task", "", "The new description of the task.")
}

const dataFilePath = "data/.pk.json"

func updateTask(taskID int, newCategory string, newIsDone *bool, newTask string) error {
	// 1. Read the existing data from the JSON file.
	file, err := os.Open(dataFilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// 2. Unmarshal the JSON data into a slice of Task structs.
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// 3. Find the task with the matching ID.
	found := false
	for i := range tasks {
		if tasks[i].ID == taskID {
			found = true

			// 4. Update the fields if new values were provided.
			if newCategory != "" {
				tasks[i].Cat = newCategory
			}
			if newIsDone != nil {
				tasks[i].Isdone = *newIsDone
			}
			if newTask != "" {
				tasks[i].Task = newTask
			}
			break // Task found and updated, exit the loop
		}
	}

	if !found {
		return errors.New("task with ID " + strconv.Itoa(taskID) + " not found")
	}

	// 5. Marshal the updated slice of Task structs back into JSON.
	updatedData, err := json.MarshalIndent(tasks, "", "  ") // Pretty print with indent
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// 6. Write the updated JSON data back to the file.
	err = ioutil.WriteFile(dataFilePath, updatedData, 0644) // Use 0644 permissions
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
