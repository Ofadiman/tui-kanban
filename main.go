package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Foreground(lipgloss.Color("241"))
)

type Task struct {
	status      status
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type Model struct {
	focused  status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/3, height-10)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[todo].Title = "Todo lists"
	m.lists[todo].SetItems([]list.Item{
		Task{
			status:      todo,
			title:       "one title",
			description: "one description",
		},
		Task{
			status:      todo,
			title:       "two title",
			description: "two description",
		},
		Task{
			status:      todo,
			title:       "three title",
			description: "three description",
		},
	})

	m.lists[inProgress].Title = "In progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{
			status:      inProgress,
			title:       "four title",
			description: "four description",
		},
	})

	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{
			status:      done,
			title:       "five title",
			description: "five description",
		},
	})
}

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		{
			if m.loaded == false {
				columnStyle.Width(msg.Width / 3)
				focusedStyle.Width(msg.Width / 3)
				columnStyle.Height(msg.Height - 20)
				focusedStyle.Height(msg.Height - 20)
				m.loaded = true
				m.initLists(msg.Width, msg.Height)
			}
		}
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "ctrl+c", "q":
				{
					m.quitting = true
					return m, tea.Quit
				}
			case "right", "l":
				{
					m.Next()
				}
			case "left", "h":
				{
					m.Prev()
				}
			}
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.quitting == true {
		return ""
	}

	if m.loaded == false {
		return ""
	}

	todoView := m.lists[todo].View()
	inProgressView := m.lists[inProgress].View()
	doneView := m.lists[done].View()

	switch m.focused {
	case todo:
		{
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		}
	case inProgress:
		{
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		}
	case done:
		{
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				focusedStyle.Render(doneView),
			)
		}
	default:
		{
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(doneView),
				columnStyle.Render(inProgressView),
			)
		}
	}
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
