package internal

func handleInput(k string, t *Tui) {
	switch t.model.peekState().active {
	case feeds, items:
		handleList(k, t.model)
	case search:
		handleSearch(k, t.model)
	}
}

func handleSearch(k string, model *Client) {
	switch k {
	case "<Escape>":
		model.popState()
	case "<Backspace>":
		n := len(model.input)
		if n > 0 {
			model.input = model.input[:n-1]
			model.cursor -= 1
		}
	case "<Enter>":
		go model.addFeed(model.input)
		model.input = ""
		model.cursor = 0
		model.popState()
	default:
		if len(k) == 1 {
			model.input += k
			model.cursor += 1
		}
	}
}

func handleList(k string, model *Client) {
	switch k {
	case "j", "<Down>":
		model.scrollDown()
	case "k", "<Up>":
		model.scrollUp()
	case "h", "<Left>", "l", "<Right>":
		if model.peekState().active == items {
			model.peekState().active = feeds
		} else {
			model.peekState().active = items
		}
	case "<Enter>":
		if model.peekState().active == items {
			model.openBrowser()
		}
	}
}
