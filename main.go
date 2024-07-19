package main

import (
	"buurro/tuition/audio"
	"buurro/tuition/bluetooth"
	"buurro/tuition/style"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type sessionState int

const (
	mainView      sessionState = iota
	audioView     sessionState = iota
	bluetoothView sessionState = iota
)

type MainModel struct {
	state     sessionState
	prevState sessionState
	size      tea.WindowSizeMsg
	list      list.Model

	audio     audio.AudioModel
	bluetooth bluetooth.BluetoothModel
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) View() string {
	switch m.state {
	case mainView:
		return style.DocStyle.Render(m.list.View())
	case audioView:
		return m.audio.View()
	case bluetoothView:
		return m.bluetooth.View()
	}
	return ""
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.state {
	case mainView:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := style.DocStyle.GetFrameSize()
			m.size = tea.WindowSizeMsg{
				Width:  msg.Width - h,
				Height: msg.Height - v,
			}
			m.list.SetSize(m.size.Width, m.size.Height)
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.prevState = mainView
				current := m.list.SelectedItem().FilterValue()
				if current == "Audio" {
					m.audio = audio.AudioModel{}
					m.state = audioView
					return m, m.audio.Init()
				}
				if current == "Bluetooth" {
					m.bluetooth = bluetooth.GetModel(m.size)
					m.state = bluetoothView
					return m, m.bluetooth.Init()
				}
			}
		}
		m.list, cmd = m.list.Update(msg)
	case audioView:
		m.audio, cmd = m.audio.Update(msg)
	case bluetoothView:
		m.bluetooth, cmd = m.bluetooth.Update(msg)
	}
	return m, cmd
}

func getModel(size tea.WindowSizeMsg) MainModel {
	items := []list.Item{
		item{title: "Audio", desc: "System audio"},
		item{title: "Bluetooth", desc: "Bluetooth devices and settings"},
	}
	list := list.New(items, list.NewDefaultDelegate(), size.Width, size.Height)
	list.Title = "ayo"

	return MainModel{state: mainView, list: list}
}

func main() {
	m := getModel(tea.WindowSizeMsg{Width: 0, Height: 0})
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
