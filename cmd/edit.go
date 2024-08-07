/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Lachignol/martin-solving/note"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var editCmd = &cobra.Command{
	Use:   "edit [index de la tache a modifier]",
	Short: "Editer une tache existante",
	Long: `Methode en ligne de commande afin de modifier le titre d'une tache.
	le mode interactif est aussi disponible en tapant show`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
		index, err := strconv.Atoi(args[0])

		nouveautitre := m.inputs[0].Value()
		if err != nil {
			log.Println(err)
		}
		note.EditTodo(index, nouveautitre)

	},
}

var m model

func initialModel() model {
	m = model{
		inputs: make([]textinput.Model, 1),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = CursorStyle
		t.CharLimit = 109

		switch i {
		case 0:
			t.Placeholder = "Nom du nouveau titre"
			t.Focus()
			t.PromptStyle = FocusedStyle
			t.TextStyle = FocusedStyle
			// case 1:
			// 	t.Placeholder = "Description"
			// 	t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				//initialisation de la variable dans le scope global a true quand on appui sur submit
				//pour ensuite l'ajouté au notes

				readyToAdd = true
				return m, tea.Quit
			}
			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = FocusedStyle
					m.inputs[i].TextStyle = FocusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = NoStyle
				m.inputs[i].TextStyle = NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder
	for i := range m.inputs {

		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &BlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString(HelpStyle.Render("tapez esc ou ctrl+c pour quitter"))
	return FormStyle.Render(b.String())
}

func init() {
	todoCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
