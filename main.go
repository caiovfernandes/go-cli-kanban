package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FF00FF"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	models []tea.Model
)

const (
	todo status = iota
	inProgress
	done
)
const (
	main status = iota
	form
)

// Task

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

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// Model

type Model struct {
	focused status
	lists   []list.Model
	err     error
}

func New() *Model {
	return &Model{}
}

func (m *Model) MoveToNext() {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, selectedTask)
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		width,
		height-5,
	)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// To Do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "Write documentation", description: "Write documentation for the project"},
		Task{status: todo, title: "Write tests", description: "Write tests for the project"},
		Task{status: todo, title: "Write code", description: "Write code for the project"},
	})

	// In Progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "SoW", description: "Write statement of work."},
		Task{status: inProgress, title: "Leverage Requirements", description: "Leverage project functional and non-functional requirements."},
		Task{status: inProgress, title: "Architectural Documentation", description: "Write documentation about the solution architecture."},
	})

	// In Progress
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "SoW", description: "Write statement of work."},
		Task{status: done, title: "Leverage Requirements", description: "Leverage project functional and non-functional requirements."},
		Task{status: done, title: "Architectural Documentation", description: "Write documentation about the solution architecture."},
	})
}

func (m *Model) Init() tea.Cmd {
	m.focused = 0
	m.lists = make([]list.Model, 3)
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		columnStyle.Width(msg.Width)
		focusedStyle.Width(msg.Width)
		m.initLists(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			if m.focused > todo {
				m.focused--
			}
		case "l", "right":
			if m.focused < done {
				m.focused++
			}
		case "enter":
			m.MoveToNext()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		getBoardStyle(m, todo, m.lists[todo]),
		getBoardStyle(m, inProgress, m.lists[inProgress]),
		getBoardStyle(m, done, m.lists[done]),
	)
}

// Functions
func getBoardStyle(m *Model, s status, listItem list.Model) string {
	if s == m.focused {
		return focusedStyle.Render(listItem.View())
	} else {
		return columnStyle.Render(listItem.View())
	}
}

func main() {
	m := New()
	m.Init()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		panic(err)
	}
}
