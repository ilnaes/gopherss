package internal

import (
	"errors"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type Panel int

const (
	feeds Panel = iota
	items
	box
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
	boxTop       int
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
		itemSelected: []int{0},
		Mutex:        &sync.Mutex{},
	}
}

func (c *Client) updateItem() {
	var feed []*Item
	if c.feedSelected > 0 {
		feed = c.Feeds[c.feedSelected-1].Items
	} else {
		feed = c.getAll()
	}

	if len(feed) > 0 {
		c.item = feed[c.itemSelected[c.feedSelected]]
	}
}

func (c *Client) openBrowser() {
	if runtime.GOOS == "darwin" {
	}
}

func (c *Client) scrollBoxUp() {
	if c.boxTop > 0 {
		c.boxTop--
	}
}

func (c *Client) scrollBoxDown() {
	c.boxTop++
	// TODO: error check this
}

func (c *Client) scrollListUp() {
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

// get a merged list of all items
// should be called while holding lock
func (c *Client) getAll() []*Item {
	items := make([]*Item, 0)

	for _, f := range c.Feeds {
		f.mu.Lock()
		items = mergeItems(items, f.Items, false)
		f.mu.Unlock()
	}

	return items
}

func (c *Client) scrollListDown() {
	c.Lock()
	defer c.Unlock()
	if c.peekState().active == items {
		var feed []*Item
		if c.feedSelected == 0 {
			feed = c.getAll()
		} else {
			feed = c.Feeds[c.feedSelected-1].Items
		}

		if c.itemSelected[c.feedSelected] < len(feed)-1 {
			c.itemSelected[c.feedSelected] += 1
			c.updateItem()
		}
	}
	if c.peekState().active == feeds && c.feedSelected < len(c.Feeds) {
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
	res := make([]string, 0)
	var items []*Item

	c.Lock()
	defer c.Unlock()

	if c.feedSelected == 0 {
		items = c.getAll()
	} else {
		feed := c.Feeds[c.feedSelected-1]
		items = feed.Items

		feed.mu.Lock()
		defer feed.mu.Unlock()
	}

	for _, i := range items {
		res = append(res, i.Title+"  "+
			strings.Replace(i.getDescription(), "\n", " ", -1))
	}

	return res
}

func (c *Client) getFeeds() ([]string, []bool) {
	feeds := []string{"All"}
	updating := []bool{false}

	c.Lock()
	for _, f := range c.Feeds {
		f.mu.Lock()
		feeds = append(feeds, f.Title)
		updating = append(updating, f.updating)
		updating[0] = updating[0] || f.updating
		f.mu.Unlock()
	}
	c.Unlock()

	return feeds, updating
}

func (c *Client) removeFeed() {
	if c.feedSelected == 0 {
		// don't delete All feed
		return
	}

	c.Lock()
	defer c.Unlock()

	c.Feeds = append(c.Feeds[:c.feedSelected-1], c.Feeds[c.feedSelected:]...)
	c.itemSelected = append(c.itemSelected[:c.feedSelected-1],
		c.itemSelected[c.feedSelected:]...)

	if c.feedSelected > len(c.Feeds) {
		c.feedSelected--
	}

}

func (c *Client) addFeed(url string) {
	feed, err := feedFromURL(url)
	if err != nil {
		return
	}

	c.Lock()
	i := sort.Search(len(c.Feeds), func(i int) bool {
		return strings.Compare(c.Feeds[i].Title, feed.Title) >= 0
	})

	// new feed
	if i == len(c.Feeds) {
		c.Feeds = append(c.Feeds, feed)
		c.itemSelected = append(c.itemSelected, 0)
	} else {
		if c.Feeds[i].Title != feed.Title {
			c.Feeds = append(c.Feeds, nil)
			copy(c.Feeds[i+1:], c.Feeds[i:])
			c.Feeds[i] = feed

			c.itemSelected = append(c.itemSelected, 0)
			copy(c.itemSelected[i+1:], c.itemSelected[i:])
			c.itemSelected[i] = 0
		}
	}
	c.Unlock()
}
