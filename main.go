// TODO - add some colour
// TODO - filter bigger lists, so like if there are twenty items add another character to filter the list

// Log - I am making a map of the keys and their index and corresponding values so I don't have to write 20 if statements and switch cases
// to do anything with the keys

package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var globalDirectory string

type KeyValues struct {
	value string
	index int
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	directory string
	err       error
	keyMaps   map[tea.KeyType]KeyValues
}

// see what this is like
func initialPath() string {
	c := exec.Command("pwd")

	output, err := c.Output()

	if err != nil {
		// not quite sure what to add here yet
	}

	return string(output)
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80
	ti.SetValue(initialPath())

	keyMapper := make(map[tea.KeyType]KeyValues)

	keyMapper[tea.KeyCtrlA] = KeyValues{value: "a", index: 0}
	keyMapper[tea.KeyCtrlS] = KeyValues{value: "s", index: 1}
	keyMapper[tea.KeyCtrlD] = KeyValues{value: "d", index: 2}
	keyMapper[tea.KeyCtrlQ] = KeyValues{value: "q", index: 3}
	keyMapper[tea.KeyCtrlW] = KeyValues{value: "w", index: 4}
	keyMapper[tea.KeyCtrlE] = KeyValues{value: "e", index: 5}
	keyMapper[tea.KeyCtrlZ] = KeyValues{value: "z", index: 6}
	keyMapper[tea.KeyCtrlX] = KeyValues{value: "x", index: 7}

	return model{
		textInput: ti,
		directory: lsd("-p"),
		keyMaps:   keyMapper,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func updateInput(m model, inputValue string) {
	m.textInput.SetValue(inputValue)
	m.textInput.CursorEnd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.directory = lsd(m.textInput.Value())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlA,
			tea.KeyCtrlS,
			tea.KeyCtrlD,
			tea.KeyCtrlQ,
			tea.KeyCtrlW,
			tea.KeyCtrlE,
			tea.KeyCtrlZ,
			tea.KeyCtrlX:

			currentKey := m.keyMaps[msg.Type]
			list := lsd(m.textInput.Value())
			input := strings.Split(list, "\n")

			var updateValue = m.textInput.Value()
			if string(updateValue[len(updateValue)-1]) != "/" {
				updateValue = updateValue + "/"
			}

			m.textInput.SetValue(updateValue + input[currentKey.index])
			m.textInput.CursorEnd()

			return m, cmd

		case tea.KeyBackspace:
			line := strings.Split(m.textInput.Value(), "/")
			m.textInput.SetValue(strings.Join(line[:len(line)-1], "/"))
			m.textInput.CursorEnd()
			return m, cmd

		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	// putting this here is slow
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		m.textInput.View(),
		lsd3(m, m.textInput.Value()),
	) + "\n"
}

func lsd(path string) string {
	// TODO: something like split the path here and if value is not / or . then do nothing
	// just to make sure it doesn't error when it is going off
	if path != "" {

		c := exec.Command("ls", path)

		output, err := c.Output()

		if err != nil {
			// I got rid of this because it blowing everything up
			// fmt.Println("Failed to run cmd:", err)
			// os.Exit(1)
			return ""
		}

		return string(output)
	} else {
		return ""
	}

}

func lsd3(m model, path string) string {
	// TODO: something like split the path here and if value is not / or . then do nothing
	// just to make sure it doesn't error when it is going off
	// TODO: name this better later

	if path != "" {

		c := exec.Command("ls", path)

		output, err := c.Output()

		if err != nil {
			// I got rid of this because it blowing everything up
			// fmt.Println("Failed to run cmd:", err)
			// os.Exit(1)
			return ""
		}

		list := strings.Split(string(output), "\n")
		shortcutAppender := func(key tea.KeyType, i int, value string) {
			if i == m.keyMaps[key].index {
				key := m.keyMaps[key]
				list[i] = key.value + " " + value
			}
		}

		for i, value := range list {
			shortcutAppender(tea.KeyCtrlA, i, value)
			shortcutAppender(tea.KeyCtrlS, i, value)
			shortcutAppender(tea.KeyCtrlD, i, value)
			shortcutAppender(tea.KeyCtrlQ, i, value)
			shortcutAppender(tea.KeyCtrlW, i, value)
			shortcutAppender(tea.KeyCtrlE, i, value)
			shortcutAppender(tea.KeyCtrlZ, i, value)
			shortcutAppender(tea.KeyCtrlX, i, value)
		}

		return strings.Join(list, "\n")
	} else {
		return ""
	}

}
