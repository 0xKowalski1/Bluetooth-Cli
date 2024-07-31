package main

import (
	"fmt"
	"os"

	bluetoothmanager "0xKowalski1/Bluetooth-Cli/Bluetooth-Manager"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

type model struct {
	choices          []string
	cursor           int
	bluetoothOn      bool
	scanning         bool
	devices          []string
	bluetoothManager *bluetoothmanager.BluetoothManager
}

func (m model) Init() tea.Cmd {
	return m.checkInitialBluetoothState
}

func (m model) checkInitialBluetoothState() tea.Msg {
	bluetoothStatus, err := m.bluetoothManager.GetBluetoothStatus()
	if err != nil {
		panic(err)
	}
	return bluetoothStateMsg(bluetoothStatus)
}

type bluetoothStateMsg bool

type deviceDiscoveredMsg struct {
	name    string
	address string
	rssi    int16
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			switch m.cursor {
			case 0:
				err := m.bluetoothManager.SetBluetoothPowerState(!m.bluetoothOn)
				if err != nil {
					fmt.Printf("ERROR: %v", err)
				}
				m.bluetoothOn = !m.bluetoothOn
				m.updateChoiceText()
			case 1:
				if m.scanning {
					m.bluetoothManager.StopScan()
					m.scanning = false
				} else {
					m.startScan()
					m.scanning = true
				}
				m.updateChoiceText()
			}
		}
	case bluetoothStateMsg:
		m.bluetoothOn = bool(msg)
		m.updateChoiceText()
	case deviceDiscoveredMsg:
		m.devices = append(m.devices, fmt.Sprintf("Name: %s, Address: %s, RSSI: %d", msg.name, msg.address, msg.rssi))
		m.updateChoiceText()
	}
	return m, nil
}

func (m *model) updateChoiceText() {
	newChoices := make([]string, 2+len(m.devices))

	if m.bluetoothOn {
		newChoices[0] = "Toggle Bluetooth: ON"
	} else {
		newChoices[0] = "Toggle Bluetooth: OFF"
	}
	if m.scanning {
		newChoices[1] = "Stop Scan"
	} else {
		newChoices[1] = "Start Scan"
	}

	for _, device := range m.devices {
		newChoices = append(newChoices, device)
	}

	m.choices = newChoices
}

func (m *model) startScan() {
	err := m.bluetoothManager.StartScan()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}

	// Listen for discovered devices
	go func() {
		for ev := range m.bluetoothManager.Discovery {
			if ev.Type == adapter.DeviceRemoved {
				continue
			}

			dev, err := device.NewDevice1(ev.Path)
			if err != nil {
				log.Errorf("%s: %s", ev.Path, err)
				continue
			}

			if dev == nil {
				log.Errorf("%s: not found", ev.Path)
				continue
			}

			msg := deviceDiscoveredMsg{
				name:    dev.Properties.Name,
				address: dev.Properties.Address,
				rssi:    dev.Properties.RSSI,
			}
			m.Update(msg) // Update model with discovered device
		}
	}()
}

func (m model) View() string {
	s := "Bluetooth Manager\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s
}

func main() {
	log.SetLevel(log.ErrorLevel) // Avoid erroneous go-bluetooth errors

	bluetoothManager, err := bluetoothmanager.NewBluetoothManager()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}

	initialModel := &model{
		choices:          []string{"Toggle Bluetooth: OFF", "Start Scan"},
		bluetoothManager: bluetoothManager,
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
