/*
Copyright © 2025 prince kumar 
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil" // Consider using a logger that integrates better with Cobra
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// DataItem represents the structure of each entry in your .pk.json file.
//
//	Note the lowercase field names and json tags to match your file.
type DataItem struct {
	Cat    string `json:"Cat"`
	Isdone bool   `json:"isdone"`
	Task   string `json:"task"`
	Time   string `json:"time"`
}

// model defines the UI model for the Bubble Tea application.
type model struct {
	table  table.Model
	data   []DataItem
	colors []lipgloss.Color
}

// Styling Definitions
var (
	// Define your colors
	yellow     = lipgloss.Color("#FFFF00")
	red        = lipgloss.Color("#FF0000")
	blue       = lipgloss.Color("#0000FF")
	pink       = lipgloss.Color("#FF69B4") // Hot Pink
	thinBorder = lipgloss.RoundedBorder()
	baseStyle  = lipgloss.NewStyle().Border(thinBorder)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

// initialModel creates and returns the initial UI model.
func initialModel(data []DataItem) model {
	columns := []table.Column{
		{Title: "Task", Width: 30},
		{Title: "Category", Width: 15},
		{Title: "Time", Width: 10},
		{Title: "Done", Width: 5},
	}

	var rows []table.Row
	for _, item := range data {
		doneStr := "No"
		if item.Isdone {
			doneStr = "Yes"
		}
		rows = append(rows, table.Row{item.Task, item.Cat, item.Time, doneStr})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()). // Keep Normal Border, or replace with thinBorder
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Bold(false)

	t.SetStyles(s)

	// Now, assign random colors to some parts.
	colors := []lipgloss.Color{yellow, red, blue, pink} // Add more colors as needed
	rand.Seed(time.Now().UnixNano())                    // Initialize random seed

	// Assign a random color to the header border
	s.Header = s.Header.BorderForeground(colors[rand.Intn(len(colors))])

	//Assign a random colour to the selection
	s.Selected = s.Selected.Background(colors[rand.Intn(len(colors))])

	return model{table: t, data: data, colors: colors}
}

// Init initializes the Bubble Tea application.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the UI model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the UI.
func (m model) View() string {
	return baseStyle.BorderForeground(m.colors[rand.Intn(len(m.colors))]).Render(fmt.Sprintf( // Change the border colour
		"%s\n%s",
		m.table.View(),
		helpStyle("Press q to quit."),
	))
}

// loadDataFromFile reads the JSON data from the .pk.json file.
func loadDataFromFile(filePath string) ([]DataItem, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var data []DataItem
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return data, nil
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "This will list all the task ",
	Long:  `This command is used to list all the task.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath := "data/.pk.json"

		// Load the data from the .pk.json file
		data, err := loadDataFromFile(filePath)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), "Error loading data:", err) // Use Cobra's output streams
			os.Exit(1)                                                  // Exit with an error code
		}

		// Initialize the Bubble Tea model with the loaded data
		initialModel := initialModel(data)

		// Start the Bubble Tea program
		p := tea.NewProgram(initialModel, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), "Error running program:", err) // Use Cobra's error stream
			os.Exit(1)                                                     // Exit with an error code
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
