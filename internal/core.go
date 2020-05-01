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
	def := dir + "/.gopherss.json"

	if len(f) == 0 {
		if _, err := os.Stat(def); err == nil {
			f = def
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

	for _, f := range cl.Feeds {
		f.mu.Lock()
		f.prune()
		f.mu.Unlock()
	}

	file, err := os.Create(def)
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
