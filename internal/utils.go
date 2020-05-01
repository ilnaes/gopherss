package internal

import (
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/net/html"
)

func htmlParse(s string) string {
	z := html.NewTokenizer(strings.NewReader(s))
	res := []byte{}

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return string(res)
		case html.StartTagToken:
			tn, _ := z.TagName()
			if len(tn) == 1 && tn[0] == 'p' {
				res = append(res, '\n')
			}
		case html.TextToken:
			res = append(res, z.Text()...)
		}
	}
}

func openURL(s string) {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("open", s)
		cmd.Run()
	}
}
