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
	queued      bool
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

func (it *Item) queue() {
	it.queued = true
}

func (it *Item) dequeue() {
	it.queued = false
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

// merges item list i2 into i1, prefering i1
// mark determines if older items in i1 get marked for deleted
func mergeItems(i1, i2 []*Item, mark bool) []*Item {
	// trivial cases
	if i1 == nil || len(i1) == 0 {
		return i2
	}
	if i2 == nil || len(i2) == 0 {
		return i1
	}

	items := make([]*Item, 0)

	n := len(i2) + len(i1)

	j := 0
	k := 0
	for i := 0; i < n; i++ {
		if j == len(i1) {
			items = append(items, i2[k])
			k++
		} else if k == len(i2) {
			if mark && !i1[j].Save {
				i1[j].Discard()
			}
			items = append(items, i1[j])
			j++
		} else {
			if i1[j].PubDate.After(*i2[k].PubDate) {
				items = append(items, i1[j])
				j++
			} else if i1[j].PubDate.Before(*i2[k].PubDate) {
				items = append(items, i2[k])
				k++
			} else {
				items = append(items, i1[j])
				j++
				k++
				i++
			}
		}
	}
	return items
}
