package internal

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type Item struct {
	Title       string
	Link        string
	Author      string
	Description string
	PubDate     *time.Time
	Save        bool
	Deleted     bool
	Read        bool
}

func itemFrom(gi *gofeed.Item) *Item {
	i := Item{
		Title:       gi.Title,
		Link:        gi.Link,
		Description: gi.Description,
	}

	if gi.Author != nil {
		i.Author = gi.Author.Name
	}

	if gi.PublishedParsed == nil {
		t := time.Now()
		i.PubDate = &t
	} else {
		i.PubDate = gi.PublishedParsed
	}

	return &i
}

func (it *Item) read() {
	it.Read = true
}
