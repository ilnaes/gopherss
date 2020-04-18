package internal

import (
	"errors"
	"runtime"
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
	navStack     []State
	feedSelected int
	itemSelected []int
	input        string
	cursor       int
}

func newClient() Client {
	return Client{
		Feeds: make([]*Feed, 0),
		navStack: []State{
			{
				route:  main,
				active: feeds,
			},
		},
		feedSelected: 0,
		itemSelected: []int{},
	}
}

func (c *Client) openBrowser() {
	if runtime.GOOS == "darwin" {
	}
}

func (c *Client) addFeed(url string) {
	feed, err := feedFromURL(url)
	if err != nil {
		return
	}

	c.Feeds = append(c.Feeds, feed)
	c.itemSelected = append(c.itemSelected, 0)
}

func (c *Client) scrollUp() {
	if c.peekState().active == items {
		if c.itemSelected[c.feedSelected] > 0 {
			c.itemSelected[c.feedSelected] -= 1
		}
	}
	if c.peekState().active == feeds {
		if c.feedSelected > 0 {
			c.feedSelected -= 1
		}
	}
}

func (c *Client) scrollDown() {
	if c.peekState().active == items {
		c.itemSelected[c.feedSelected] += 1
	}
	if c.peekState().active == feeds && c.feedSelected < len(c.Feeds)-1 {
		c.feedSelected += 1
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

	feed := c.Feeds[c.feedSelected]

	feed.mu.Lock()
	defer feed.mu.Unlock()

	l := len(feed.Items)
	for i := range feed.Items {
		if !feed.Items[l-i-1].Deleted {
			items = append(items, feed.Items[l-i-1].Title+"  "+
				strings.Replace(feed.Items[l-i-1].getDescription(), "\n", " ", -1))
		}
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
