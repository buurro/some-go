package main

import (
	"buurro/tuition/audio"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type sessionState int

const (
	entryView sessionState = iota
	audioView sessionState = iota
)

type MainModel struct {
	state sessionState
	list  list.Model

	audio audio.AudioModel
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) View() string {
	switch m.state {
	case entryView:
		return docStyle.Render(m.list.View())
	case audioView:
		return m.audio.View()
	}
	return ""
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	switch m.state {
	case entryView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				current := m.list.SelectedItem().FilterValue()
				if current == "Audio" {
					m.audio = audio.AudioModel{}
					m.state = audioView
					return m, m.audio.Init()
				}
			}
		}
		m.list, cmd = m.list.Update(msg)
	case audioView:
		m.audio, cmd = m.audio.Update(msg)
	}

	return m, cmd
}

func main() {
	items := []list.Item{
		item{title: "Audio", desc: "System audio"},
		item{title: "Bluetooth", desc: "Bluetooth devices and settings"},
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "ayo"

	m := MainModel{state: entryView, list: list}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
