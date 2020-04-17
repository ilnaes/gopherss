package internal

import (
	"sync"
	"testing"
)

func TestFeedCreate(t *testing.T) {
	f := Feed{
		FeedLink: "",
		mu:       &sync.Mutex{},
	}

	err := f.updateFromStr(`<rss version="2.0">
				<channel>
				<title>What</title>
				<link>http://localhost</link>
				<description>Testing</description>
				</rss>
			   `)
	if err != nil {
		t.Error("Couldn't updateFromStr: ", err)
	}

	if f.Title != "What" {
		t.Error("Incorrect title")
	}

	if f.Description != "Testing" {
		t.Error("Incorrect description")
	}

	if len(f.Items) != 0 {
		t.Error("Incorrect items")
	}
}

func TestFeedItems(t *testing.T) {
	f := Feed{
		FeedLink: "",
		mu:       &sync.Mutex{},
	}

	err := f.updateFromStr(`<rss version="2.0">
				<channel>
				<title>What</title>
				<link>http://localhost</link>
				<description>Testing</description>
				<item>
					<title>I1</title>
					<pubDate>Thu, 09 Apr 2020 06:40:37 +0000</pubDate>
				</item>
				</rss>
			   `)
	if err != nil {
		t.Error("Couldn't updateFromStr: ", err)
	}

	if len(f.Items) != 1 {
		t.Error("Didn't get items")
	}
	if f.Items[0].PubDate == nil {
		t.Error("Didn't get date")
	}
	f.Items[0].read()

	err = f.updateFromStr(`<rss version="2.0">
				<channel>
				<title>What</title>
				<link>http://localhost</link>
				<description>Testing</description>
				<item>
					<title>I2</title>
					<pubDate>Wed, 08 Apr 2020 06:40:37 +0000</pubDate>
				</item>
				</rss>
			   `)
	if err != nil {
		t.Error("Couldn't updateFromStr: ", err)
	}

	if len(f.Items) != 1 {
		t.Error("Items didn't merge correctly")
	}

	err = f.updateFromStr(`<rss version="2.0">
				<channel>
				<title>What</title>
				<link>http://localhost</link>
				<description>Testing</description>
				<item>
					<title>I3</title>
					<pubDate>Fri, 10 Apr 2020 06:40:37 +0000</pubDate>
				</item>
				</rss>
			   `)
	if err != nil {
		t.Error("Couldn't updateFromStr: ", err)
	}

	if len(f.Items) != 2 {
		t.Error("Items didn't merge correctly")
	}

	if !f.Items[0].Read {
		t.Error("trampled on flag")
	}

	if f.Items[1].Title != "I3" {
		t.Error("Didn't read second item correctly")
	}
}
