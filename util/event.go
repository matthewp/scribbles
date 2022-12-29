package util

import "github.com/nbd-wtf/go-nostr"

type Event struct {
	Event nostr.Event
	Nick  string
}
