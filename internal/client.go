package internal

import (
	"errors"
	"runtime"
	"strings"
	"sync"
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
	searchOn     bool
	item         *Item
	*sync.Mutex
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
		Mutex:        &sync.Mutex{},
	}
}

func (c *Client) updateItem() {
	f := c.feedSelected
	if len(c.Feeds) > 0 && len(c.Feeds[f].Items) > 0 {
		c.item = c.Feeds[f].Items[c.itemSelected[f]]
	}
}

func (c *Client) openBrowser() {
	if runtime.GOOS == "darwin" {
	}
}

func (c *Client) scrollUp() {
	c.Lock()
	defer c.Unlock()
	if c.peekState().active == items {
		if c.itemSelected[c.feedSelected] > 0 {
			c.itemSelected[c.feedSelected] -= 1
			c.updateItem()
		}
	}
	if c.peekState().active == feeds {
		if c.feedSelected > 0 {
			c.feedSelected -= 1
		}
	}
}

func (c *Client) scrollDown() {
	c.Lock()
	defer c.Unlock()
	if c.peekState().active == items {
		if c.itemSelected[c.feedSelected] < len(c.Feeds[c.feedSelected].Items)-1 {
			c.itemSelected[c.feedSelected] += 1
			c.updateItem()
		}
	}
	if c.peekState().active == feeds && c.feedSelected < len(c.Feeds)-1 {
		c.feedSelected += 1
	}
}

func (c *Client) reload() {
	for {
		c.Lock()
		for _, f := range c.Feeds {
			go func(f *Feed) {
				f.update()
			}(f)
		}
		c.Unlock()

		// TODO: change this
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

func (c *Client) getItems() []string {
	items := make([]string, 0)

	if len(c.Feeds) == 0 {
		return items
	}

	c.Lock()

	feed := c.Feeds[c.feedSelected]
	feed.mu.Lock()

	for _, i := range feed.Items {
		items = append(items, i.Title+"  "+
			strings.Replace(i.getDescription(), "\n", " ", -1))
	}

	feed.mu.Unlock()
	c.Unlock()

	return items
}

func (c *Client) getFeeds() []string {
	feeds := make([]string, 0)

	c.Lock()
	for _, f := range c.Feeds {
		f.mu.Lock()
		feeds = append(feeds, f.Title)
		f.mu.Unlock()
	}
	c.Unlock()

	return feeds
}

func (c *Client) removeFeed() {
	c.Lock()
	defer c.Unlock()

	c.Feeds = append(c.Feeds[:c.feedSelected], c.Feeds[c.feedSelected+1:]...)
	c.itemSelected = append(c.itemSelected[:c.feedSelected],
		c.itemSelected[c.feedSelected+1:]...)

	if c.feedSelected >= len(c.Feeds) {
		c.feedSelected--
	}

}

func (c *Client) addFeed(url string) {
	feed, err := feedFromURL(url)
	if err != nil {
		return
	}

	c.Lock()
	c.Feeds = append(c.Feeds, feed)
	c.itemSelected = append(c.itemSelected, 0)
	c.Unlock()
}
