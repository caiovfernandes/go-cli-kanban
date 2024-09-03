package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FF00FF"))

	models []tea.Model
)

const (
	todo status = iota
	inProgress
	done
)
const (
	model status = iota
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

func NewTask(status status, title, description string) Task {
	return Task{status: status, title: title, description: description}
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
		Task{status: done, title: "Project Idea", description: "zzzzz."},
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
		case "n":
			models[model] = m
			models[form].(*Form).focused = m.focused
			return models[form], nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case Task:
		task := msg
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
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

type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
}

func NewForm(focused status) *Form {
	form := &Form{focused: focused}
	form.title = textinput.New()
	form.title.Focus()
	form.description = textarea.New()
	return form
}

func (f Form) CreateTaskFromForm(m *Model) tea.Msg {
	m.lists[f.focused].InsertItem(len(m.lists[f.focused].Items()), NewTask(f.focused, f.title.Value(), f.description.Value()))
	return NewTask(f.focused, f.title.Value(), f.description.Value())
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if f.title.Focused() {
				f.title.Blur()
				f.description.Focus()
				return f, textarea.Blink
			} else {
				models[form] = f
				f.CreateTaskFromForm(models[model].(*Model))
				models[form] = NewForm(f.focused)
				return models[model], nil
			}
		}
	}
	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	} else {
		f.description, cmd = f.description.Update(msg)
		return f, cmd
	}
}

func (f Form) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, f.title.View(), f.description.View())
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
	models = []tea.Model{New(), NewForm(todo)}
	m := models[model]
	m.Init()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		panic(err)
	}
}
