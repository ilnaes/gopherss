package internal

import (
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"time"
)

type feed struct {
	FeedLink    string
	Title       string
	Link        string
	Description string
	LastDate    string
	Items       []*item
}

type item struct {
	Title     string
	Link      string
	Author    string
	Published *time.Time
}

func itemFrom(gi *gofeed.Item) *item {
	return &item{
		Title:     gi.Title,
		Link:      gi.Link,
		Author:    gi.Author.Name,
		Published: gi.PublishedParsed,
	}
}

func (f *feed) refresh() error {
	fp := gofeed.NewParser()

	// gf, err := fp.ParseURL(f.FeedLink)
	dat, err := ioutil.ReadFile(f.FeedLink)
	gf, err := fp.ParseString(string(dat))
	if err != nil {
		return err
	}

	items := make([]*item, 0)

	for _, i := range gf.Items {
		items = append(items, itemFrom(i))
	}

	f.Title = gf.Title
	f.Link = gf.Link
	f.Description = gf.Description
	f.LastDate = gf.Updated
	f.Items = items

	return nil
}

func feedFromURL(s string) (*feed, error) {
	f := feed{
		FeedLink: s,
	}

	err := f.refresh()
	if err != nil {
		return nil, err
	}

	return &f, nil
}
