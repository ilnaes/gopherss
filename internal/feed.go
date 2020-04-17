package internal

import (
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
	fp := gofeed.NewParser()

	gf, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
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
	return &f
}

// merges the new feed nf into current feed
// note: assumes nf is transitory so no need to lock it
func (f *Feed) merge(nf *Feed) {
	Items := make([]*Item, 0)

	l := len(nf.Items)
	// Items are in chronological ascending order
	for i := range nf.Items {
		Items = append(Items, nf.Items[l-i-1])
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.Title = nf.Title
	f.Link = nf.Link
	f.Description = nf.Description
	f.LastDate = nf.LastDate

	// append only Items that come later
	if f.Items == nil || len(f.Items) == 0 {
		f.Items = Items
	} else {
		last := f.Items[len(f.Items)-1].PubDate

		for i := range Items {
			if last.Before(*Items[i].PubDate) {
				f.Items = append(f.Items, Items[i:]...)
				break
			}
		}
	}
}

// get new items
func (f *Feed) update() error {
	nf, err := feedFromURL(f.FeedLink)
	if err != nil {
		return err
	}

	f.merge(nf)
	return nil
}

func (f *Feed) updateFromStr(s string) error {
	nf, err := feedFromStr(s)
	if err != nil {
		return err
	}

	f.merge(nf)
	return nil
}
