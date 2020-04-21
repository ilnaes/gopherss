package internal

import (
	"errors"
	"io/ioutil"
	"sync"

	"github.com/mmcdole/gofeed"
)

type Feed struct {
	mu          *sync.Mutex
	FeedLink    string
	Title       string
	Link        string
	Description string
	LastDate    string
	Items       []*Item
}

// feed from url
func feedFromURL(url string) (*Feed, error) {
	n := len(url)

	if n == 0 {
		return nil, errors.New("Empty URL")
	}
	fp := gofeed.NewParser()

	gf, err := fp.ParseURL(url)
	if err != nil {
		gf, err = fp.ParseURL(url + "/rss")
		if err != nil {
			return nil, err
		}
	}

	return fromGofeed(gf), nil
}

func feedFromStr(s string) (*Feed, error) {
	fp := gofeed.NewParser()
	gf, err := fp.ParseString(s)
	if err != nil {
		return nil, err
	}

	return fromGofeed(gf), nil
}

// feed from xml file
func feedFromFile(file string) (*Feed, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return feedFromStr(string(dat))
}

// project gofeed.Feed onto Feed
func fromGofeed(gf *gofeed.Feed) *Feed {
	items := make([]*Item, 0)
	for _, i := range gf.Items {
		items = append(items, itemFrom(i))
	}

	f := Feed{
		mu:          &sync.Mutex{},
		FeedLink:    gf.FeedLink,
		Title:       gf.Title,
		Link:        gf.Link,
		Description: gf.Description,
		LastDate:    gf.Updated,
		Items:       items,
	}

	if len(gf.FeedLink) == 0 {
		// TODO: error catch all of this
		if gf.Extensions != nil {
			f.FeedLink = gf.Extensions["atom"]["link"][0].Attrs["href"]
		}
	}
	return &f
}

// removes feeds that are discarded
func (f *Feed) prune() {
	j := 0
	for i := range f.Items {
		if !f.Items[i].discard {
			f.Items[j] = f.Items[i]
			j++
		}
	}

	f.Items = f.Items[:j]
}

// merges the new feed nf into current feed
// note: assumes nf is transitory so no need to lock it
// note: also assumes feeds are in reverse chronological order
func (f *Feed) updateFrom(nf *Feed) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.Title = nf.Title
	f.Link = nf.Link
	f.Description = nf.Description
	f.LastDate = nf.LastDate

	if len(nf.FeedLink) != 0 {
		f.FeedLink = nf.FeedLink
	}

	f.Items = mergeItems(f.Items, nf.Items, true)
}

// get new items
func (f *Feed) update() error {
	nf, err := feedFromURL(f.FeedLink)
	if err != nil {
		return err
	}

	f.updateFrom(nf)
	return nil
}

func (f *Feed) updateFromStr(s string) error {
	nf, err := feedFromStr(s)
	if err != nil {
		return err
	}

	f.updateFrom(nf)
	return nil
}
