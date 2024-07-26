package main

import (
    "fmt"
    "os"
    log "github.com/sirupsen/logrus"
    tea "github.com/charmbracelet/bubbletea"
    bluetoothmanager "0xKowalski1/Bluetooth-Cli/Bluetooth-Manager"
)

type model struct {
    choices           []string
    cursor            int
    bluetoothOn       bool
    bluetoothManager  *bluetoothmanager.BluetoothManager
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

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "enter", " ":
            if m.cursor == 0 {
                err := m.bluetoothManager.SetBluetoothPowerState(!m.bluetoothOn)
                if err != nil {
                    fmt.Printf("ERROR: %v", err)
                }
                m.bluetoothOn = !m.bluetoothOn
                m.updateChoiceText()
            }
        }
    case bluetoothStateMsg:
        m.bluetoothOn = bool(msg) // Ensure boolean conversion
        m.updateChoiceText()
    }
    return m, nil
}

func (m *model) updateChoiceText() {
    if m.bluetoothOn {
        m.choices[0] = "Toggle Bluetooth: ON"
    } else {
        m.choices[0] = "Toggle Bluetooth: OFF"
    }
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
    log.SetLevel(log.ErrorLevel) // Hide erroneous logs

    bluetoothManager, err := bluetoothmanager.NewBluetoothManager()
    if err != nil {
        fmt.Printf("ERROR: %v", err)
    }

    initialModel := &model{
        choices:          []string{"Loading..."},
        bluetoothManager: bluetoothManager,
    }

    p := tea.NewProgram(initialModel)
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}

