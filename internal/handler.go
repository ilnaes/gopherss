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
		model.searchOn = true
		model.popState()
		go func() {
			model.addFeed(model.input)
			model.searchOn = false
			model.input = ""
			model.cursor = 0
		}()
	default:
		if len(k) == 1 {
			model.input += k
			model.cursor += 1
		}
	}
}

func handleList(k string, model *Client) {
	active := &model.peekState().active
	switch k {
	case "j", "<Down>":
		model.scrollDown()
	case "k", "<Up>":
		model.scrollUp()
	case "h", "<Left>", "l", "<Right>":
		if *active == items {
			*active = feeds
		} else {
			*active = items
			model.updateItem()
		}
	case "<Enter>":
		if *active == items {
			// TODO
		}
	case "d":
		if *active == feeds {
			model.removeFeed()
		}
	}
}
