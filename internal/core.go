package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
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

	queued := []string{}

	for _, f := range cl.Feeds {
		f.mu.Lock()

		for _, i := range f.Items {
			if i.queued {
				queued = append(queued, i.Link)
			}
		}

		f.prune()
		f.mu.Unlock()
	}

	if len(queued) > 0 {
		cmd := exec.Command("open", queued...)
		_ = cmd.Run()
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
