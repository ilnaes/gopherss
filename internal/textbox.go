package internal

import (
	"image"

	. "github.com/gizak/termui/v3"
)

type Textbox struct {
	Block
	Text       string
	TextStyle  Style
	Top        int
	WrapText   bool
	Cursor     int
	ShowCursor bool
}

func NewTextbox() *Textbox {
	return &Textbox{
		Block:     *NewBlock(),
		TextStyle: Theme.Paragraph.Text,
		WrapText:  true,
	}
}

func (self *Textbox) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	cells := make([]Cell, 0)
	for _, r := range []rune(self.Text) {
		cells = append(cells, Cell{Rune: r, Style: self.TextStyle})
	}
	if self.WrapText {
		cells = WrapCells(cells, uint(self.Inner.Dx()))
	}

	if self.ShowCursor {
		if self.Cursor >= len(self.Text) {
			cells = append(cells, Cell{Style: Style{Bg: ColorYellow}})
		} else {
			cells[self.Cursor].Style.Bg = ColorYellow
		}
	}

	rows := SplitCells(cells, '\n')[self.Top:]

	for y, row := range rows {
		if y+self.Inner.Min.Y >= self.Inner.Max.Y {
			break
		}
		row = TrimCells(row, self.Inner.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(self.Inner.Min))
		}
	}
}
