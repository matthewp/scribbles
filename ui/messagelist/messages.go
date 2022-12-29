package messagelist

import (
	"github.com/charmbracelet/bubbles/list"
)

type Message struct {
	Created int64
	Text    string
}

type ByCreated []list.Item

func (a ByCreated) Len() int      { return len(a) }
func (a ByCreated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCreated) Less(i, j int) bool {
	ui := a[i].(item)
	uj := a[j].(item)
	return ui.Event.CreatedAt.Unix() > uj.Event.CreatedAt.Unix()
}
