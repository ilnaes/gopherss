package internal

import (
	"errors"
	"strings"
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

type Route int

const (
	main Route = iota
)

type State struct {
	route  Route
	active Panel
}

type Client struct {
	Feeds        []*Feed
	items        []*Item
	navStack     []State
	feedSelected int
	itemSelected int
}

func newClient() Client {
	return Client{
		Feeds: make([]*Feed, 0),
		items: nil,
		navStack: []State{
			{
				route:  main,
				active: feeds,
			},
		},
		feedSelected: 0,
		itemSelected: 0,
	}
}

func (c *Client) scrollUp() {
	if c.itemSelected > 0 {
		c.itemSelected -= 1
	}
}

func (c *Client) scrollDown() {
	c.itemSelected += 1
}

func (c *Client) reload() {
	for {
		for _, f := range c.Feeds {
			f.update()
		}
		<-time.Tick(tickTime)
	}
}

func (c *Client) pushState(s State) {
	c.navStack = append(c.navStack, s)
}

func (c *Client) peekState() *State {
	return &c.navStack[len(c.navStack)-1]
}

func (c *Client) popState() (State, error) {
	n := len(c.navStack)
	if n == 1 {
		return State{}, errors.New("Can't pop anymore")
	}

	res := c.navStack[n-1]
	c.navStack = c.navStack[:n-1]

	return res, nil
}

func (c *Client) getState(s State) *State {
	n := len(c.navStack)
	if n == 0 {
		return nil
	}

	return &c.navStack[n-1]
}

func (c *Client) getItems() []string {
	items := make([]string, 0)

	if len(c.Feeds) == 0 {
		return items
	}

	c.Feeds[0].mu.Lock()
	defer c.Feeds[0].mu.Unlock()

	l := len(c.Feeds[0].Items)
	for i := range c.Feeds[0].Items {
		items = append(items, c.Feeds[0].Items[l-i-1].Title+"  "+
			strings.Replace(c.Feeds[0].Items[l-i-1].getDescription(), "\n", " ", -1))
	}

	return items
}

func (c *Client) getFeeds() []string {
	feeds := make([]string, 0)

	for _, f := range c.Feeds {
		f.mu.Lock()
		feeds = append(feeds, f.Title)
		f.mu.Unlock()
	}
	return feeds
}
