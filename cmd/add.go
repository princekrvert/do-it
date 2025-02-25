/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// make a slice to add these tasks to
type Task struct {
	Task   string    `json:"task"`
	Cat    string    `json :"cat"`
	Time   time.Time `json:"time"`
	Isdone bool      `json:"isdone"`
}

// Function to add data into json file

// Function to add JSON data to the .pk.json file
func addDataToPKJSON(filePath string, newData interface{}) error {
	// 1. Read the existing data from the file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644) // Open for read/write, create if not exists
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	var existingData []interface{} //  Assume it's a JSON array

	if fileInfo.Size() > 0 { //Only try to read if there is content
		byteValue, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		// 2. Unmarshal the existing data (assuming it's an array)
		err = json.Unmarshal(byteValue, &existingData)
		if err != nil {

			//If its not an array, try unmarshaling into a single object
			var singleObject interface{}
			errSingle := json.Unmarshal(byteValue, &singleObject)
			if errSingle == nil {
				existingData = append(existingData, singleObject)
			} else {
				existingData = []interface{}{} // Initialize to an empty array if unmarshaling fails.
			}

		}
	} else {
		existingData = []interface{}{} // Initialize if file is empty
	}

	// 3. Append the new data to the existing data
	existingData = append(existingData, newData)

	// 4. Marshal the updated JSON data back into bytes
	updatedJSON, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	// 5. Write the updated JSON data back to the file
	err = file.Truncate(0) // Clear existing content
	if err != nil {
		return fmt.Errorf("error truncating file: %w", err)
	}
	_, err = file.Seek(0, 0) // Reset pointer to the beginning of the file

	_, err = file.Write(updatedJSON)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "This command is used to add the task",
	Long:  `This command is used to add the tasks in the to do list`,
	Run: func(cmd *cobra.Command, args []string) {
		task, _ := cmd.Flags().GetString("task")
		cat, _ := cmd.Flags().GetString("cat")
		first := Task{task, cat, time.Now(), false}
		//json_data, err := json.Marshal(first)
		//if err != nil {
		//	fmt.Println("\033[31;1m Some error occured")
		//}
		//fmt.Println(string(json_data))
		err := addDataToPKJSON("data/.pk.json", first)
		if err != nil {
			fmt.Println(err)
		}

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
	addCmd.Flags().StringP("task", "t", "", "Provide the task ")
	addCmd.Flags().StringP("cat", "c", "", "Provide the cattogry ")
}
