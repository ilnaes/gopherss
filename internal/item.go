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
	Content     string
	PubDate     *time.Time
	Save        bool
	Read        bool
	discard     bool
}

func itemFrom(gi *gofeed.Item) *Item {
	i := Item{
		Title:       gi.Title,
		Link:        gi.Link,
		Description: gi.Description,
		Content:     gi.Content,
	}

	if len(i.Content) == 0 {
		i.Content = gi.Description
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

func (it *Item) setUnread() {
	it.Read = false
}

func (it *Item) setRead() {
	it.Read = true
}

func (it *Item) Discard() {
	it.discard = true
	it.Save = false
}

func (it *Item) setSave() {
	it.discard = false
	it.Save = true
}

func (it *Item) getDescription() string {
	return htmlParse(it.Description)
}

func (it *Item) getContent() string {
	if len(it.Content) > 0 {
		return htmlParse(it.Content)
	} else {
		return htmlParse(it.Description)
	}
}
