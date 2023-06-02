// TODO - search through strings that contain substrings
// TODO - use to lower case so I dont have to type caps
// TODO - add some colour
package main

import (
	"fmt"
	"os"
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
	p := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr))
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	m, err := p.Run()

	if err != nil {

	}

	if m, ok := m.(model); ok && m.textInput.Value() != "" {
		fmt.Printf("%s\n", m.textInput.Value())
	}

}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	directory string
	path      string
	err       error
	keyMaps   map[tea.KeyType]KeyValues
}

func initialPath() string {
	c := exec.Command("pwd")

	output, err := c.Output()

	if err != nil {
		// not quite sure what to add here yet
	}

	return string(output)
}

func makeKeyMaps() map[tea.KeyType]KeyValues {
	keyMapper := make(map[tea.KeyType]KeyValues)

	keyMapper[tea.KeyCtrlA] = KeyValues{value: "a", index: 0}
	keyMapper[tea.KeyCtrlS] = KeyValues{value: "s", index: 1}
	keyMapper[tea.KeyCtrlD] = KeyValues{value: "d", index: 2}
	keyMapper[tea.KeyCtrlQ] = KeyValues{value: "q", index: 3}
	keyMapper[tea.KeyCtrlW] = KeyValues{value: "w", index: 4}
	keyMapper[tea.KeyCtrlE] = KeyValues{value: "e", index: 5}
	keyMapper[tea.KeyCtrlZ] = KeyValues{value: "z", index: 6}
	keyMapper[tea.KeyCtrlX] = KeyValues{value: "x", index: 7}

	return keyMapper
}

func lastCharacter(str string) string {
	return string(str[len(str)-1])
}

func listDirectory(path string) string {
	strippedPath := strings.TrimSuffix(path, "\n")
	keyMaps := makeKeyMaps()
	shortcutAppender := func(list []string, key tea.KeyType, i int, value string) {
		remove_newline := strings.TrimPrefix(value, "\n")
		if i == keyMaps[key].index && remove_newline != "" {
			key := keyMaps[key]
			list[i] = key.value + " " + value
		}
	}

	if strippedPath != "" {

		c := exec.Command("ls", strippedPath)

		output, err := c.Output()

		// Putting a partial path into the ls command here causes an error
		// e.g. ~/play/code/b
		// the b represents the beginning of the bubble-tea directory
		// the command cant handle that and errors out here
		// this is where we then list the last path and then filter results
		// TODO - check for this before executing the first ls
		if err != nil {
			previous_output := strings.Split(strippedPath, "/")
			previous_output_minus_after_slash := strings.Join(previous_output[:len(previous_output)-1], "/")

			co := exec.Command("ls", previous_output_minus_after_slash)

			innerOutput, innerErr := co.Output()
			if innerErr != nil {

			}

			var filtered = []string{}
			for _, value := range strings.Split(string(innerOutput), "\n") {
				if strings.HasPrefix(string(value), previous_output[len(previous_output)-1]) {
					filtered = append(filtered, string(value))
				}
			}

			for i, value := range filtered {
				shortcutAppender(filtered, tea.KeyCtrlA, i, value)
				shortcutAppender(filtered, tea.KeyCtrlS, i, value)
				shortcutAppender(filtered, tea.KeyCtrlD, i, value)
				shortcutAppender(filtered, tea.KeyCtrlQ, i, value)
				shortcutAppender(filtered, tea.KeyCtrlW, i, value)
				shortcutAppender(filtered, tea.KeyCtrlE, i, value)
				shortcutAppender(filtered, tea.KeyCtrlZ, i, value)
				shortcutAppender(filtered, tea.KeyCtrlX, i, value)
			}

			return strings.Join(filtered, "\n")
		}

		list := strings.Split(string(output), "\n")

		for i, value := range list {
			shortcutAppender(list, tea.KeyCtrlA, i, value)
			shortcutAppender(list, tea.KeyCtrlS, i, value)
			shortcutAppender(list, tea.KeyCtrlD, i, value)
			shortcutAppender(list, tea.KeyCtrlQ, i, value)
			shortcutAppender(list, tea.KeyCtrlW, i, value)
			shortcutAppender(list, tea.KeyCtrlE, i, value)
			shortcutAppender(list, tea.KeyCtrlZ, i, value)
			shortcutAppender(list, tea.KeyCtrlX, i, value)
		}

		// 		return strings.Join(list, "\n")

		return strings.Join(list, "\n")

	} else {
		return " path isnt empty "
	}

}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80
	ti.SetValue(initialPath())

	return model{
		textInput: ti,
		path:      initialPath(),
		directory: listDirectory(initialPath()),
		keyMaps:   makeKeyMaps(),
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

	m.directory = listDirectory(m.textInput.Value())
	m.path = m.textInput.Value()

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
			list := listDirectory(m.textInput.Value())
			input := strings.Split(list, "\n")

			var updateValue = m.textInput.Value()

			remove_tea_key := strings.Split(input[currentKey.index], " ")[1]

			// if last character is not /
			split_path := strings.Split(updateValue, "/")
			path_after_last_slash := split_path[len(split_path)-1]
			if strings.HasPrefix(remove_tea_key, path_after_last_slash) {
				if string(updateValue[len(updateValue)-1]) != "/" {
					// apparently i dont need this just yet
					// updateValue = updateValue + "/"
				}

				// I am here and first letter prefixes work but multi letter ones dont
				m.textInput.SetValue(strings.TrimSuffix(updateValue, path_after_last_slash) + remove_tea_key + "/")
			} else {
				if string(updateValue[len(updateValue)-1]) != "/" {
					updateValue = updateValue + "/"
				}
				m.textInput.SetValue(updateValue + remove_tea_key + "/")
			}

			m.textInput.CursorEnd()

			return m, cmd

		case tea.KeyBackspace:
			line := strings.Split(m.textInput.Value(), "/")
			if lastCharacter(m.textInput.Value()) != "/" {
				m.textInput.SetValue(strings.Join(line[:len(line)-1], "/") + "/")
			} else {
				m.textInput.SetValue(strings.Join(line[:len(line)-2], "/") + "/")
			}
			m.textInput.CursorEnd()
			return m, cmd

		case tea.KeyEnter:
			// in here I want to exit and change the directory to m.textInput.value()
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
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
		listDirectory(m.textInput.Value()),
	) + "\n"
}
