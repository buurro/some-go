package bluetooth

import (
	// "bufio"
	"buurro/tuition/style"

	// "os/exec"
	// "strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type errorMsg struct {
	error
}

func (e errorMsg) Error() string {
	return e.error.Error()
}

type deviceListMsg struct {
	devices []device
}

type BluetoothModel struct {
	deviceList []device
	list       list.Model
	error      error
}

func (m BluetoothModel) Init() tea.Cmd {
	return fetchListData
}

func (m BluetoothModel) View() string {
	if m.error == nil && m.deviceList == nil {
		return "loading..."
	}
	if m.error != nil {
		return m.error.Error()
	}
	if m.deviceList != nil {
		return style.DocStyle.Render(m.list.View())
	}
	return ""
}

func (m BluetoothModel) Update(msg tea.Msg) (BluetoothModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case errorMsg:
		m.error = msg.error
	case deviceListMsg:
		m.deviceList = msg.devices
		items := make([]list.Item, len(m.deviceList))
		for i, dev := range m.deviceList {
			items[i] = dev
		}
		m.list.SetItems(items)
	case tea.WindowSizeMsg:
		h, v := style.DocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.list, cmd = m.list.Update(msg)
	}
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

type device struct {
	name, address string
}

func (d device) Title() string       { return d.name }
func (d device) Description() string { return d.address }
func (d device) FilterValue() string { return d.name }

func fetchListData() tea.Msg {
	var devices []device
	var err error

	devices, err = FetchPairedDevices()
	// devices, err = []device{
	// 	{name: "AAAAAA", address: "12:af:23:00"},
	// 	{name: "BBB BB", address: "13:bf:13:01"},
	// }, nil

	if err != nil {
		return errorMsg{err}
	}

	return deviceListMsg{devices}
}

func GetModel(size tea.WindowSizeMsg) BluetoothModel {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), size.Width, size.Height)
	list.KeyMap.Quit.SetKeys("q")
	list.Title = "Bluetooth"
	return BluetoothModel{
		list: list,
	}
}
