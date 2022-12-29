package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matthewp/scribbles/ui/messagelist"
	"github.com/matthewp/scribbles/ui/sidebar"
	"github.com/matthewp/scribbles/util"
	"github.com/mitchellh/go-homedir"
	"github.com/nbd-wtf/go-nostr"
)

func main() {
	flag.StringVar(&config.DataDir, "datadir", "~/.config/nostr",
		"Base directory for configurations and data from Nostr.")
	flag.Parse()
	config.DataDir, _ = homedir.Expand(config.DataDir)
	os.Mkdir(config.DataDir, 0700)

	// logger config
	log.SetPrefix("<> ")

	// parse config
	path := filepath.Join(config.DataDir, "config.json")
	_, err := os.Open(path)
	if err != nil {
		saveConfig(path)
		_, _ = os.Open(path)
	}
	f, _ := os.Open(path)
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Fatal("can't parse config file " + path + ": " + err.Error())
		return
	}
	config.Init()

	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	go func() {
		initNostr()
		var keys []string

		nameMap := map[string]string{}

		for _, follow := range config.Following {
			keys = append(keys, follow.Key)

			if follow.Name != "" {
				nameMap[follow.Key] = follow.Name
			}
		}

		_, all, _ := pool.Sub(nostr.Filters{{Authors: keys}})
		for event := range nostr.Unique(all) {
			nick, ok := nameMap[event.PubKey]
			if !ok {
				nick = ""
			}
			// If we don't already have a nick for this user, and they are announcing their
			// new name, let's use it.
			if nick == "" {
				if event.Kind == nostr.KindSetMetadata {
					var metadata Metadata
					err := json.Unmarshal([]byte(event.Content), &metadata)
					if err != nil {
						//log.Println("Failed to parse metadata.")
						continue
					}

					nick = metadata.Name
					nameMap[nick] = event.PubKey
				}
			}
			p.Send(util.Event{
				Event: event,
				Nick:  nick,
			})
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	messagelist messagelist.Model
	textarea    textarea.Model
	sidebar     sidebar.Model
	err         error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messagelist: messagelist.New(),
		sidebar:     sidebar.New(),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		mlCmd tea.Cmd
		sbCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.messagelist, mlCmd = m.messagelist.Update(msg)
	m.sidebar, sbCmd = m.sidebar.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newHeight := msg.Height - m.textarea.Height()
		m.messagelist.SetHeight(newHeight)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
			/*case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()*/
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, mlCmd, sbCmd)
}

func (m model) View() string {
	rightPane := lipgloss.JoinVertical(lipgloss.Left, m.messagelist.View(), m.textarea.View())
	return lipgloss.JoinHorizontal(lipgloss.Top, m.sidebar.View(), rightPane)
}
