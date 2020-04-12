package internal

import (
	"fmt"
)

type client struct {
	Feeds []feed
}

func Run(f string) {
	feed, err := feedFromURL(f)
	if err != nil {
		panic(err)
	}

	for _, i := range feed.Items {
		fmt.Printf("%+v\n", i)
	}
}
