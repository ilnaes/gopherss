package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"sync"
)

func Run(f string) {
	cl := newClient()

	usr, _ := user.Current()
	dir := usr.HomeDir

	if len(f) == 0 {
		if _, err := os.Stat(dir + "/.gopherss.json"); err == nil {
			f = dir + "/.gopherss.json"
		}
	}

	if len(f) > 0 {
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

	file, err := os.Create(dir + "/.gopherss.json")
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
