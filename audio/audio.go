package audio

import (
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type AudioModel struct {
	deviceList []string
	error      error
}

type errorMsg struct {
	error
}

func (e errorMsg) Error() string {
	return e.error.Error()
}

type deviceListMsg struct {
	devices []string
}

func (m AudioModel) View() string {
	if m.error == nil && m.deviceList == nil {
		return "loading..."
	}
	if m.error != nil {
		return m.error.Error()
	}
	if m.deviceList != nil {
		return strings.Join(m.deviceList, ", ")
	}
	return ""
}

func (m AudioModel) Init() tea.Cmd {
	return fetchAudioDevices
}

func (m AudioModel) Update(msg tea.Msg) (AudioModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case errorMsg:
		m.error = msg.error
	case deviceListMsg:
		m.deviceList = msg.devices
	}

	return m, cmd
}

func fetchAudioDevices() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get("https://marco.ooo/")
	if err != nil {
		return errorMsg{err}
	}
	defer res.Body.Close() // nolint:errcheck

	return deviceListMsg{devices: []string{"a", "b", "c"}}
}
