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
	audioView   sessionState = iota
	audioInput  sessionState = iota
	audioOutput sessionState = iota
	entryView
)

type MainModel struct {
	state sessionState
	entry tea.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

type Node struct {
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
	var nodes []Node

	err = json.Unmarshal(out, &nodes)
	if err != nil {
		log.Fatal(err)
	}

	items := []list.Item{
		item{title: "Volume", desc: "System and playback volume"},
		item{title: "Output", desc: "Output Devices"},
		item{title: "Input", desc: "Input Devices"},
	}

	// for _, node := range nodes {
	// 	if node.Type == "PipeWire:Interface:Node" &&
	// 		(node.Info.Props.Media_class == "Audio/Sink" || node.Info.Props.Media_class == "Audio/Source") {
	// 		items = append(items, item{
	// 			title: node.Info.Props.Node_nick,
	// 			desc:  node.Info.Props.Media_class,
	// 		})
	// 	}
	// }

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Audio"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
