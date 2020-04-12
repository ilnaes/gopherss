package main

import (
	"fmt"
	c "github.com/ilnaes/gopherss/internal"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Need an argument!")
	} else {
		c.Run(os.Args[1])
	}
}
