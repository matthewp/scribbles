package messagelist

import (
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matthewp/scribbles/util"
)

type Model struct {
	height   int
	items    []list.Item
	viewport viewport.Model
	list     list.Model
}

var (
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}).
			Margin(1, 0, 0, 0).
			Padding(1, 1).
			Width(100).
			Height(itemHeight).MaxHeight(itemHeight)
	selectedItemStyle = itemStyle.Copy().
				Background(lipgloss.Color("#E4BDFB")).
				Foreground(lipgloss.Color("#2E2E2E"))
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	quitTextStyle   = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func New() Model {
	items := []list.Item{}

	const defaultWidth = 20
	const listHeight = 40

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(true)
	l.Styles.PaginationStyle = paginationStyle

	return Model{
		height: listHeight,
		items:  items,
		list:   l,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		lCmd tea.Cmd
	)

	m.list, lCmd = m.list.Update(msg)

	switch msg := msg.(type) {
	case util.Event:
		if !unprintededEvents(msg.Event) {
			i := item{
				Event: msg.Event,
				Nick:  msg.Nick,
			}
			items := append(m.list.Items(), i)
			sort.Sort(ByCreated(items))
			m.list.SetItems(items)
		}
	}
	return m, tea.Batch(lCmd)
}

func (m Model) View() string {
	return m.list.View()
}

func (m Model) SetHeight(h int) {
	m.list.SetHeight(h)
	m.height = h
}
