package internal

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Tui struct {
	model *Client
	l     *widgets.List
}

func (t *Tui) start() error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	items := t.model.getItems()
	for i := range items {
		items[i] = fmt.Sprintf("[%d] %s", i, items[i])
	}

	w, _ := ui.TerminalDimensions()

	l := widgets.NewList()
	l.Title = "List"
	l.Rows = items
	l.TextStyle = ui.NewStyle(7)
	l.SelectedRowStyle = ui.NewStyle(15)
	l.WrapText = false
	l.Block.BorderStyle = ui.NewStyle(ui.ColorCyan)
	l.SetRect(0, 0, w, 15)

	t.l = l

	t.render()
	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "j", "<Down>":
				t.l.ScrollDown()
			case "k", "<Up>":
				t.l.ScrollUp()
			case "<C-d>":
				t.l.ScrollHalfPageDown()
			case "<C-u>":
				t.l.ScrollHalfPageUp()
			case "<C-f>":
				t.l.ScrollPageDown()
			case "<C-b>":
				t.l.ScrollPageUp()
			case "g":
				if previousKey == "g" {
					t.l.ScrollTop()
				}
			case "<Home>":
				t.l.ScrollTop()
			case "G", "<End>":
				t.l.ScrollBottom()
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}
		case <-time.Tick(refreshTime):
		}

		t.render()
	}
}

func (t *Tui) render() {
	items := t.model.getItems()
	for i := range items {
		items[i] = fmt.Sprintf("[%d] %s", i, items[i])
	}
	t.l.Rows = items

	ui.Render(t.l)
}
