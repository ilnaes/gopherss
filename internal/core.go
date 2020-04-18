package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

func Run(f string) {
	cl := newClient()

	if strings.Contains(f, "xml") {
		feed, err := feedFromFile(f)
		if err != nil {
			panic(err)
		}

		cl.Feeds = append(cl.Feeds, feed)

	} else if strings.Contains(f, "json") {
		dat, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}

		json.Unmarshal(dat, &cl.Feeds)

		for _, f := range cl.Feeds {
			f.mu = &sync.Mutex{}
			cl.itemSelected = append(cl.itemSelected, 0)
		}
	}

	t := Tui{
		model: &cl,
	}

	go cl.reload()
	t.start()

	file, err := os.Create("fm.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ppjs, err := json.MarshalIndent(cl.Feeds, "", "  ")
	if err != nil {
		panic(err)
	}

	file.Write(ppjs)
}
