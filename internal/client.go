package internal

import (
	"time"
)

type Panel int

const (
	feeds Panel = iota
	items
	item
	dialog
	search
	help
)

type State struct {
	panel  Panel
	active bool
}

type Client struct {
	Feeds    []*Feed
	items    []*Item
	navStack []State
}

func newClient() Client {
	return Client{
		Feeds:    make([]*Feed, 0),
		items:    nil,
		navStack: make([]State, 0),
	}
}

func (c *Client) reload() {
	for {
		for _, f := range c.Feeds {
			f.update()
		}
		<-time.Tick(tickTime)
	}
}
