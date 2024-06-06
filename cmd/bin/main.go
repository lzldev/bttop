package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	_ "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lzldev/bttop/internal/sysmanager"
)

func main() {
	m := NewAppModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		panic(fmt.Errorf("program error: %w", err))
	}
}

type AppModel struct {
	width, height int
	bar           progress.Model
	comms         chan<- sysmanager.SysActorMessage
	logger        <-chan string
	sysmsg        <-chan sysmanager.SysMessage
	ram_pct       float64
	cpu_pct       float64
	counter       int
	start         time.Time
	messages      []string
}

func NewAppModel() AppModel {
	bar := progress.New(progress.WithDefaultGradient())

	bar.PercentFormat = " %.2f %%"

	logger := make(chan string)
	comms, sysmsg := sysmanager.StartSysManager(logger)

	return AppModel{
		bar:      bar,
		comms:    comms,
		logger:   logger,
		sysmsg:   sysmsg,
		start:    time.Now(),
		messages: []string{"1", "2"},
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(m.listenLogs, m.listenSys)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sysmanager.SysMessage:
		m.ram_pct = msg.RamPct
		m.cpu_pct = msg.CpuPct
		return m, m.listenSys
	case LoggerMessage:
		m.messages = append(m.messages, string(msg))
		return m, m.listenLogs
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "w":
			m.comms <- sysmanager.SpeedDown
			return m, nil
		case "e":
			m.comms <- sysmanager.SpeedUp
			return m, nil
		case "r":
			m.comms <- sysmanager.Nothing
			return m, nil
		case "D":
			m.messages = make([]string, 0)
			return m, nil
		case "x":
			m.messages = m.messages[:max(0, (len(m.messages)-1))]
			return m, nil
		case "d":
			m.messages = m.messages[min(len(m.messages)-1, 1):]
			return m, nil
		case "a":
			return m, tea.EnterAltScreen
		case "s":
			return m, tea.ExitAltScreen
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.counter++
			return m, nil
		}
	}
	return m, nil
}

func (m AppModel) View() (r string) {
	// r += "Ram: "
	// r += m.bar.ViewAs(1 - m.ram_pct)
	// r += " Free: "
	// r += m.bar.ViewAs(m.ram_pct)

	r += fmt.Sprintf("\nCount : %v\nSince: %v\n", m.counter, time.Since(m.start))

	header := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.width-2).
		Padding(1, 2)

	totalLabel := "Total "
	freeLabel := "Free "

	cw := m.width - header.GetHorizontalFrameSize()
	meterPadding := 2
	meterWidth := (cw / 2)

	title := lipgloss.NewStyle().
		Padding(0, 0).
		MarginBottom(1).
		Border(lipgloss.NormalBorder(), false).
		BorderBottom(true).
		Width(cw).
		Bold(true)

	m.bar.Width = (meterWidth - len(totalLabel)) - meterPadding
	total := lipgloss.PlaceHorizontal(meterWidth, lipgloss.Center,
		totalLabel+m.bar.ViewAs(1-m.ram_pct))

	m.bar.Width = (meterWidth - len(freeLabel)) - meterPadding
	free := lipgloss.PlaceHorizontal(meterWidth, lipgloss.Center,
		freeLabel+m.bar.ViewAs(m.ram_pct))

	m.bar.Width = cw
	r += header.Render(
		title.Render("CPU"),
		m.bar.ViewAs(m.cpu_pct),
	)
	r += "\n"

	r += header.Render(
		title.Render("Ram"),
		lipgloss.JoinHorizontal(lipgloss.Center, total, free),
	)

	st := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("#874BFD")).
		Border(lipgloss.RoundedBorder()).
		Width(m.width / 2).
		// Padding(4, 4).
		Align(lipgloss.Center)

	for i, msg := range m.messages {

		r += lipgloss.Place(m.width, 8, lipgloss.Center, lipgloss.Center, st.Render(fmt.Sprintf("%v:%v", i, msg)),
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
		)
		// r += st.Render(fmt.Sprintf("%v:%v", i, msg))
		r += "\n"
		// r += fmt.Sprintf("%v:%v\n", i, msg)
	}

	return r
}

type LoggerMessage string

func (m *AppModel) listenLogs() tea.Msg {
	m.messages = append(m.messages, "Listening")
	return LoggerMessage(<-m.logger)
}

func (m *AppModel) listenSys() tea.Msg {
	return <-m.sysmsg
}
