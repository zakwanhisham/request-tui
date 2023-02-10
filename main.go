package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh/"

type model struct {
	status int
	err    error
}

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)

	if err != nil {
		return errMsg{err}
	}
	defer res.Body.Close()

	return statusMsg(res.StatusCode)
}

// func checkSomeUrl(url string) tea.Cmd {
// 	return func() tea.Msg {
// 		c := &http.Client{Timeout: 10 * time.Second}
// 		res, err := c.Get(url)
// 		if err != nil {
// 			return errMsg{err}
// 		}
// 		return statusMsg(res.StatusCode)
// 	}
// }

type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// if you want to check the server, change here
func (m model) Init() tea.Cmd { return checkServer }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}
	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m model) View() string {
	s := fmt.Sprintf("Checking %s", url)
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	} else if m.status != 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return s + "\n"
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
