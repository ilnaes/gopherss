package internal

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Tui struct {
	model       *Client
	items       *widgets.List
	feeds       *widgets.List
	width       int
	height      int
	previousKey string
}

func (t *Tui) start() error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Articles"
	l.TextStyle = ui.NewStyle(7)
	l.SelectedRowStyle = ui.NewStyle(15)
	l.WrapText = false

	t.items = l

	l = widgets.NewList()
	l.Title = "Feeds"
	l.TextStyle = ui.NewStyle(7)
	l.SelectedRowStyle = ui.NewStyle(15)
	l.WrapText = false

	t.width, t.height = ui.TerminalDimensions()

	t.feeds = l

	t.render()
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.Type {
			case ui.ResizeEvent:
				t.width, t.height = ui.TerminalDimensions()
			case ui.KeyboardEvent:
				if t.handleKey(e.ID) {
					return nil
				}
			}

		case <-time.Tick(refreshTime):
		}

		t.render()
	}
}

func (t *Tui) handleKey(key string) bool {
	switch key {
	case "q", "<C-c>":
		return true
	case "j", "<Down>":
		t.model.scrollDown()
	case "k", "<Up>":
		t.model.scrollUp()
	case "<C-d>":
		t.items.ScrollHalfPageDown()
	case "<C-u>":
		t.items.ScrollHalfPageUp()
	case "<C-f>":
		t.items.ScrollPageDown()
	case "<C-b>":
		t.items.ScrollPageUp()
	case "g":
		if t.previousKey == "g" {
			t.items.ScrollTop()
		}
	case "<Home>":
		t.items.ScrollTop()
	case "G", "<End>":
		t.items.ScrollBottom()
	case "<Right>", "<Left>":
		state := t.model.peekState()
		if state.route == main {
			if state.active == feeds {
				state.active = items
			} else if state.active == items {
				state.active = feeds
			}
		}
	}

	if t.previousKey == "g" {
		t.previousKey = ""
	} else {
		t.previousKey = key
	}

	return false
}

func (t *Tui) render() {
	state := t.model.peekState()

	it := t.model.getItems()
	for i := range it {
		it[i] = fmt.Sprintf("[%d] %s", i+1, it[i])
	}
	t.items.Rows = it
	t.items.SelectedRow = t.model.itemSelected
	if state.active == items {
		t.items.Block.BorderStyle = ui.NewStyle(ui.ColorCyan)
		t.items.TitleStyle = ui.NewStyle(ui.ColorCyan)
	} else {
		t.items.Block.BorderStyle = ui.NewStyle(ui.ColorWhite)
		t.items.TitleStyle = ui.NewStyle(ui.ColorWhite)
	}
	t.items.SetRect(t.width/2, 0, t.width, 15)

	fd := t.model.getFeeds()
	for i := range fd {
		fd[i] = fmt.Sprintf("[%d] %s", i+1, fd[i])
	}
	t.feeds.Rows = fd
	t.feeds.SelectedRow = t.model.feedSelected
	if state.active == feeds {
		t.feeds.Block.BorderStyle = ui.NewStyle(ui.ColorCyan)
		t.feeds.TitleStyle = ui.NewStyle(ui.ColorCyan)
	} else {
		t.feeds.Block.BorderStyle = ui.NewStyle(ui.ColorWhite)
		t.feeds.TitleStyle = ui.NewStyle(ui.ColorWhite)
	}
	t.feeds.SetRect(0, 0, t.width/2, 15)

	ui.Render(t.items)
	ui.Render(t.feeds)
}
