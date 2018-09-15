package plateform

import (
	"fmt"

	"github.com/gizak/termui"
	"github.com/pkg/errors"
)

const maxRowSize = 12

type termUI struct {
	body *termui.Grid
	row  []*termui.Row
}

// NewTermUI returns a new Terminal Interface object with a given output mode.
func NewTermUI() (*termUI, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	// set the basic properties
	body := termui.NewGrid()
	body.X = 0
	body.Y = 0
	body.BgColor = termui.ThemeAttr("bg")
	body.Width = termui.TermWidth()

	return &termUI{
		body: body,
		row:  []*termui.Row{},
	}, nil
}

func (termUI) Close() {
	termui.Close()
}

func (t *termUI) TextBox(
	data string,
	fg uint16,
	bd uint16,
	bdlabel string,
	h int,
	size int,
) {
	textBox := termui.NewPar(data)

	textBox.TextFgColor = termui.Attribute(fg)
	textBox.BorderFg = termui.Attribute(bd)
	textBox.BorderLabel = bdlabel
	textBox.Height = h

	t.row = append(t.row, termui.NewCol(size, 0, textBox))
}

func (t *termUI) BarChart(data []int, dimensions []string, barWidth int, bdLabel string, size int) {
	bc := termui.NewBarChart()
	bc.BorderLabel = bdLabel
	bc.Data = data
	bc.BarWidth = barWidth
	bc.BarGap = 0
	bc.DataLabels = dimensions
	bc.Width = 200
	bc.Height = 10
	bc.TextColor = termui.ColorGreen
	bc.BarColor = termui.ColorRed
	bc.NumColor = termui.ColorYellow

	t.row = append(t.row, termui.NewCol(size, 0, bc))
}

// KQuit set a key to quit the application.
func (termUI) KQuit(key string) {
	termui.Handle(fmt.Sprintf("/sys/kbd/%s", key), func(termui.Event) {
		termui.StopLoop()
	})
}

func (t *termUI) AddRow() error {
	err := t.validateRowSize()
	if err != nil {
		return err
	}

	t.body.AddRows(termui.NewRow(t.row...))
	// clean the internal row
	t.row = []*termui.Row{}

	return nil
}

func (t termUI) validateRowSize() error {
	var ts int
	for _, r := range t.row {
		for _, c := range r.Cols {
			ts += c.Offset
		}
	}

	if ts > maxRowSize {
		return errors.Errorf("could not create row: size %d too big", ts)
	}

	return nil
}

func (t *termUI) Render() {
	t.body.Align()
	termui.Render(t.body)
	termui.Loop()
}
