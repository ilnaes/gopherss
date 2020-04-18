package main

import (
	"os"

	c "github.com/ilnaes/gopherss/internal"
)

func main() {
	if len(os.Args) == 1 {
		c.Run("")
	} else {
		c.Run(os.Args[1])
	}
}
