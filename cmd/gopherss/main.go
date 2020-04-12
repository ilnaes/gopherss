package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
)

func main() {
	dat, e := ioutil.ReadFile("fm.xml")
	if e != nil {
		fmt.Println("Could not read file!")
		panic(e)
	}

	fp := gofeed.NewParser()

	feed, e := fp.ParseString(string(dat))

	fmt.Printf("%+v\n", feed)
}
