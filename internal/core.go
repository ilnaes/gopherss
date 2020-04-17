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
		feed, err := feedFromFile(f)
		if err != nil {
			panic(err)
		}

		cl := newClient()
		t := Tui{
			model: &cl,
		}

		cl.Feeds = append(cl.Feeds, feed)

		file, err := os.Create("fm.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		ppjs, err := json.MarshalIndent(cl, "", "  ")
		if err != nil {
			panic(err)
		}

		file.Write(ppjs)

		go cl.reload()
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
