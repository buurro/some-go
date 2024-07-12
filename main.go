package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	entryView       sessionState = iota
	audioView       sessionState = iota
	audioInputView  sessionState = iota
	audioOutputView sessionState = iota
)

type MainModel struct {
	state sessionState
	entry tea.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// item
type item struct {
	title, desc string
	targetState sessionState
}

func (i item) Title() string             { return i.title }
func (i item) Description() string       { return i.desc }
func (i item) FilterValue() string       { return i.title }
func (i item) TargetState() sessionState { return i.targetState }

// model
type model struct {
	list  list.Model
	state sessionState
	items map[sessionState][]list.Item
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) TargetState(item list.Item) sessionState {

	for _, i := range m.items[m.state] {
		if i.FilterValue() == item.FilterValue() {
			return audioInputView
		}
	}
	return entryView
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			m.state = audioView
			m.list.SetItems(m.items[m.state])
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

// node
type node struct {
	Id   int32
	Type string
	Info struct {
		Props struct {
			Node_nick   string `json:"node.nick"`
			Media_class string `json:"media.class"`
		}
	}
}

func main() {
	out, err := exec.Command("pw-dump").Output()
	if err != nil {
		log.Fatal(err)
	}
	var nodes []node

	err = json.Unmarshal(out, &nodes)
	if err != nil {
		log.Fatal(err)
	}

	var items = make(map[sessionState][]list.Item)

	items[entryView] = []list.Item{
		item{title: "Audio", desc: "System audio", targetState: audioView},
		item{title: "Bluetooth", desc: "Bluetooth devices and settings"},
	}

	items[audioView] = []list.Item{
		item{title: "Volume", desc: "System and playback volume"},
		item{title: "Output", desc: "Output Devices", targetState: audioOutputView},
		item{title: "Input", desc: "Input Devices", targetState: audioInputView},
	}

	items[audioInputView] = []list.Item{}
	items[audioOutputView] = []list.Item{}

	for _, node := range nodes {
		if node.Type == "PipeWire:Interface:Node" {
			item := item{
				title: node.Info.Props.Node_nick,
				desc:  node.Info.Props.Media_class,
			}
			if node.Info.Props.Media_class == "Audio/Source" {
				items[audioInputView] = append(items[audioInputView], item)
			}
			if node.Info.Props.Media_class == "Audio/Sink" {
				items[audioOutputView] = append(items[audioOutputView], item)
			}
		}
	}

	m := model{items: items}

	m.list = list.New(m.items[m.state], list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Audio"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
