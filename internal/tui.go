package internal

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func getSpinner() string {
	spins := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return spins[time.Now().Nanosecond()/100000000]
}

type Tui struct {
	model  *Client
	items  *widgets.List
	feeds  *widgets.List
	box    *widgets.Paragraph
	search *Textbox
	width  int
	height int
}

func (t *Tui) start() error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	tb := NewTextbox()
	tb.Title = "Search"
	tb.TextStyle = ui.NewStyle(7)
	t.search = tb

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

	p := widgets.NewParagraph()
	p.TextStyle = ui.NewStyle(7)
	t.box = p

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
		if t.model.peekState().active == search {
			t.model.popState()
		} else {
			return true
		}
	case "q":
		if t.model.peekState().active == search {
			handleSearch(key, t.model)
		} else {
			_, err := t.model.popState()
			if err != nil {
				return true
			}
		}
	case "/":
		if !t.model.searchOn {
			state := t.model.peekState()
			if state.route == main && state.active != search {
				t.model.pushState(State{
					route:  main,
					active: search,
				})
			} else {
				handleInput(key, t)
			}
		}
	default:
		handleInput(key, t)
	}

	return false
}

func (t *Tui) drawSearch() {
	if t.model.searchOn {
		t.search.Text = getSpinner() + " " + t.model.input
	} else {
		t.search.Text = "  " + t.model.input
	}
	if t.model.peekState().active == search {
		t.search.Block.BorderStyle = ui.NewStyle(ui.ColorCyan)
		t.search.TitleStyle = ui.NewStyle(ui.ColorCyan)
		t.search.Cursor = t.model.cursor + 2
		t.search.ShowCursor = true
	} else {
		t.search.Block.BorderStyle = ui.NewStyle(ui.ColorWhite)
		t.search.TitleStyle = ui.NewStyle(ui.ColorWhite)
		t.search.ShowCursor = false
	}
	t.search.SetRect(0, 0, t.width, 3)
}

func (t *Tui) drawFeeds() {
	fd := t.model.getFeeds()
	for i := range fd {
		fd[i] = fmt.Sprintf("[%d] %s", i+1, fd[i])
	}
	t.feeds.Rows = fd
	t.feeds.SelectedRow = t.model.feedSelected
	if t.model.peekState().active == feeds {
		t.feeds.Block.BorderStyle = ui.NewStyle(ui.ColorCyan)
		t.feeds.TitleStyle = ui.NewStyle(ui.ColorCyan)
	} else {
		t.feeds.Block.BorderStyle = ui.NewStyle(ui.ColorWhite)
		t.feeds.TitleStyle = ui.NewStyle(ui.ColorWhite)
	}
	t.feeds.SetRect(0, 3, t.width/3, 15)
}

func (t *Tui) drawItemsList() {
	itemList := t.model.getItems()
	for i := range itemList {
		itemList[i] = fmt.Sprintf("[%d] %s", i+1, itemList[i])
	}
	t.items.Rows = itemList
	if len(itemList) > 0 {
		t.items.SelectedRow = t.model.itemSelected[t.model.feedSelected]
	}
	if t.model.peekState().active == items {
		t.items.Block.BorderStyle = ui.NewStyle(ui.ColorCyan)
		t.items.TitleStyle = ui.NewStyle(ui.ColorCyan)
	} else {
		t.items.Block.BorderStyle = ui.NewStyle(ui.ColorWhite)
		t.items.TitleStyle = ui.NewStyle(ui.ColorWhite)
	}
	t.items.SetRect(t.width/3, 3, t.width, 15)
}

func (t *Tui) drawItemWindow() {
	if t.model.item != nil {
		t.box.Text = t.model.item.getContent()
	}

	t.box.Block.BorderStyle = ui.NewStyle(ui.ColorWhite)
	t.box.SetRect(0, 15, t.width, t.height)
}

func (t *Tui) render() {
	t.drawSearch()
	t.drawFeeds()
	t.drawItemsList()
	t.drawItemWindow()

	ui.Render(t.items, t.feeds, t.search, t.box)
}
