package messagelist

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/matthewp/scribbles/util"
	"github.com/nbd-wtf/go-nostr"
	"gopkg.in/yaml.v2"
)

var kindNames = map[int]string{
	nostr.KindSetMetadata:            "Profile Metadata",
	nostr.KindTextNote:               "📝",
	nostr.KindRecommendServer:        "Relay Recommendation",
	nostr.KindContactList:            "Contact List",
	nostr.KindEncryptedDirectMessage: "Encrypted Message",
	nostr.KindDeletion:               "Deletion Notice",
	nostr.KindReaction:               "🤯",
}

/*
var (
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
)
*/

func unprintededEvents(evt nostr.Event) bool {
	switch evt.Kind {
	case nostr.KindEncryptedDirectMessage:
		if !evt.Tags.ContainsAny("p", nostr.Tag{util.GetConfigPubKey()}) {
			return true
		}
	case nostr.KindReaction:
		return true
	}
	return false
}

func printEvent(evt nostr.Event, nick *string, verbose bool) string {
	kind, ok := kindNames[evt.Kind]
	if !ok {
		kind = fmt.Sprint(evt.Kind)
	}

	// Don't print encrypted messages that aren't for me
	// Note we should never get here
	if unprintededEvents(evt) {
		return ""
	}

	var fromField string = shorten(evt.PubKey)

	if nick != nil && *nick != "" {
		fromField = fmt.Sprintf("%s (%s)", *nick, shorten(evt.PubKey))
	}

	if verbose {
		if nick == nil || *nick == "" {
			fromField = evt.PubKey
		} else {
			fromField = fmt.Sprintf("%s (%s)", *nick, evt.PubKey)
		}
	}

	var ukind = lipgloss.NewStyle().
		Padding(0, 1).
		SetString(kind)
	var ufrom = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#04B575")).
		Padding(0, 1).
		SetString(fromField)

	// The actual post content
	var content string = ""

	switch evt.Kind {
	case nostr.KindSetMetadata:
		var metadata util.Metadata
		err := json.Unmarshal([]byte(evt.Content), &metadata)
		if err != nil {
			content = fmt.Sprintf("Invalid JSON: '%s',\n  %s",
				err.Error(), evt.Content)
		} else {
			y, _ := yaml.Marshal(metadata)
			spl := strings.Split(string(y), "\n")
			for i, v := range spl {
				spl[i] = "  " + v
			}
			content = strings.Join(spl, "\n")
		}
	case nostr.KindTextNote:
		content = "  " + strings.ReplaceAll(evt.Content, "\n", "\n  ")
	case nostr.KindReaction:
		content = fmt.Sprintf("%+v", evt)
	case nostr.KindRecommendServer:
	case nostr.KindContactList:
	case nostr.KindEncryptedDirectMessage:
	default:
		content = evt.Content
	}

	ht := humanize.Time(evt.CreatedAt)
	inner := fmt.Sprintf("%s %s  %s\n\n%s", ukind, ufrom, ht, content)
	return inner
}

func shorten(id string) string {
	if len(id) < 12 {
		return id
	}
	return id[0:4] + "..." + id[len(id)-4:]
}
