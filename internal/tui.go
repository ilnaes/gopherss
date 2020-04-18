package internal

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Tui struct {
	model  *Client
	items  *widgets.List
	feeds  *widgets.List
	search *Textbox
	width  int
	height int
}

func (t *Tui) start() error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	p := NewTextbox()
	p.Title = "Search"
	p.PaddingLeft = 1
	p.TextStyle = ui.NewStyle(7)
	t.search = p

	l1 := widgets.NewList()
	l1.Title = "Articles"
	l1.TextStyle = ui.NewStyle(7)
	l1.SelectedRowStyle = ui.NewStyle(15)
	l1.WrapText = false
	t.items = l1

	l2 := widgets.NewList()
	l2.Title = "Feeds"
	l2.TextStyle = ui.NewStyle(7)
	l2.SelectedRowStyle = ui.NewStyle(15)
	l2.WrapText = false
	t.feeds = l2

	t.width, t.height = ui.TerminalDimensions()

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
	case "<C-c>":
		return true
	case "/":
		state := t.model.peekState()
		if state.route == main && state.active != search {
			t.model.pushState(State{
				route:  main,
				active: search,
			})
		} else {
			handleInput(key, t)
		}
	default:
		handleInput(key, t)
	}

	return false
}

func (t *Tui) render() {
	state := t.model.peekState()

	// draw search
	t.search.Text = t.model.input
	if state.active == search {
		t.search.Block.BorderStyle = ui.NewStyle(ui.ColorCyan)
		t.search.TitleStyle = ui.NewStyle(ui.ColorCyan)
		t.search.Cursor = t.model.cursor
		t.search.ShowCursor = true
	} else {
		t.search.Block.BorderStyle = ui.NewStyle(ui.ColorWhite)
		t.search.TitleStyle = ui.NewStyle(ui.ColorWhite)
		t.search.ShowCursor = false
	}
	t.search.SetRect(0, 0, t.width, 3)

	// draw feeds
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
	t.feeds.SetRect(0, 3, t.width/2, 15)

	// draw items
	it := t.model.getItems()
	for i := range it {
		it[i] = fmt.Sprintf("[%d] %s", i+1, it[i])
	}
	t.items.Rows = it

	if len(it) > 0 {
		t.items.SelectedRow = t.model.itemSelected[t.model.feedSelected]
	}
	if state.active == items {
		t.items.Block.BorderStyle = ui.NewStyle(ui.ColorCyan)
		t.items.TitleStyle = ui.NewStyle(ui.ColorCyan)
	} else {
		t.items.Block.BorderStyle = ui.NewStyle(ui.ColorWhite)
		t.items.TitleStyle = ui.NewStyle(ui.ColorWhite)
	}
	t.items.SetRect(t.width/2, 3, t.width, 15)

	ui.Render(t.items, t.feeds, t.search)
}
