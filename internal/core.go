package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type client struct {
	Feeds []*feed
}

func Run(f string) {
	if strings.Contains(f, "xml") {
		c := client{
			Feeds: make([]*feed, 0),
		}

		feed, err := feedFromURL(f)
		if err != nil {
			panic(err)
		}

		c.Feeds = append(c.Feeds, feed)

		ppjs, err := json.MarshalIndent(c, "", "  ")
		fmt.Println(string(ppjs))

		file, err := os.Create("fm.json")
		defer file.Close()

		file.Write(ppjs)
	} else {
		dat, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}

		var c client
		json.Unmarshal(dat, &c)

		fmt.Println(c)
	}
}
