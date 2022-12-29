package sidebar

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	list list.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

func New() Model {

	items := []list.Item{
		item{title: "home"},
		item{title: "notifications"},
		item{title: "messages"},
		item{title: "profile"},
		item{title: "settings"},
	}

	del := list.NewDefaultDelegate()
	del.ShowDescription = false
	l := list.New(items, del, 0, 0)
	l.Title = "My Fave Things"
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	return Model{
		list: l,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}
