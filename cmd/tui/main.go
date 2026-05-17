package main

import (
	"fmt"
	"os"
	"strconv"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"home-server-hub/internal/models"
)

var docStyle = lipgloss.NewStyle().Margin(1, 1)

type Item struct {
	models.Application
}

// implement the list.Item interface
func (i Item) FilterValue() string {
	return i.Name
}

func (i Item) Title() string {
	return i.Name
}

func (i Item) Description() string {
	if i.URL != "" {
		return i.URL
	}

	return ("http://" + i.IP + ":" + strconv.Itoa(int(i.Port)))
}

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	v := tea.NewView(docStyle.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func main() {
	items := []list.Item{
		Item{
			Application: models.Application{
				Name: "Jellyfin",
				URL:  "http://jellyfin.home",
			},
		},

		Item{
			Application: models.Application{
				Name: "AdGuard",
				IP:   "192.168.0.100",
				Port: 8082,
			},
		},
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Applications"

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
