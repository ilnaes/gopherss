package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Run(f string) {
	if strings.Contains(f, "xml") {
		c := newClient()

		feed, err := feedFromFile(f)
		if err != nil {
			panic(err)
		}

		c.Feeds = append(c.Feeds, feed)

		file, err := os.Create("fm.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		ppjs, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			panic(err)
		}

		file.Write(ppjs)

		t := tui{}
		t.start()

	} else {
		dat, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}

		var c Client
		json.Unmarshal(dat, &c)

		fmt.Println(c.Feeds[0].Items[0].Description)
	}
}
