// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file is based on gonum.org/v1/plot/align_test.go which is
// Copyright ©2017 The gonum Authors. All rights reserved.

package plotext

import (
	"math"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"github.com/pplcc/plotext/internal"
)

func TestAt(t *testing.T) {

	//    2     3    1
	// +----+------+--+
	// |    |      |  |
	// |    |      |  | 3
	// |    |      |  |
	// +----+------+--+
	// |    |      |  | 1
	// +----+------+--+
	// |    |      |  |
	// |    |      |  | 2
	// +----+------+--+
	// +----+------+--+ 0.5

	table := Table{
		ColWidths:  []float64{2, 3, 1},
		RowHeights: []float64{3, 1, 2, 0.5},
		PadTop:     1,
		PadBottom:  2,
		PadRight:   3,
		PadLeft:    4,
		PadX:       5,
		PadY:       6,
	}

	img := vgimg.New(vg.Points(150), vg.Points(175))
	dc := draw.New(img)

	for j := 0; j < len(table.RowHeights); j++ {
		for i := 0; i < len(table.ColWidths); i++ {
			c := table.At(dc, i, j)
			// frame
			c.StrokeLines(plotter.DefaultLineStyle, []vg.Point{c.Rectangle.Min, {X: c.Rectangle.Min.X, Y: c.Rectangle.Max.Y}, c.Rectangle.Max, {X: c.Rectangle.Max.X, Y: c.Rectangle.Min.Y}, c.Rectangle.Min})
			// cross
			c.StrokeLine2(plotter.DefaultLineStyle, c.Rectangle.Min.X, c.Rectangle.Min.Y, c.Rectangle.Max.X, c.Rectangle.Max.Y)
			c.StrokeLine2(plotter.DefaultLineStyle, c.Rectangle.Min.X, c.Rectangle.Max.Y, c.Rectangle.Max.X, c.Rectangle.Min.Y)
		}
	}
	testFile := "testdata/tableat.png"
	w, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}

	internal.TestImage(t, testFile)
}
func TestAlign(t *testing.T) {

	//    2     3    1
	// +----+------+--+
	// |    |      |  |
	// |    |      |  | 3
	// |    |      |  |
	// +----+------+--+
	// |    |      |  | 1
	// +----+------+--+
	// |    |      |  |
	// |    |      |  | 2
	// +----+------+--+
	// +----+------+--+ 0.5

	table := Table{
		ColWidths:  []float64{2, 3, 1},
		RowHeights: []float64{3, 1, 2, 0.5},
		PadTop:     1,
		PadBottom:  2,
		PadRight:   3,
		PadLeft:    4,
		PadX:       5,
		PadY:       6,
	}

	plots := make([][]*plot.Plot, len(table.RowHeights))
	for j := 0; j < len(table.RowHeights); j++ {
		plots[j] = make([]*plot.Plot, len(table.ColWidths))
		for i := 0; i < len(table.ColWidths); i++ {
			if i == 0 && j == 2 {
				// This shows what happens when there are nil plots.
				continue
			}

			p, err := plot.New()
			if err != nil {
				panic(err)
			}

			if j == 0 && i == 2 {
				// This shows what happens when the axis padding
				// is different among plots.
				p.X.Padding, p.Y.Padding = 0, 0
			}

			if true && j == 1 && i == 1 {
				// To test the Align function, we make the axis labels
				// on one of the plots stick out.
				p.Y.Max = 1e9
				p.X.Max = 1e9
				p.X.Tick.Label.Rotation = math.Pi / 2
				p.X.Tick.Label.XAlign = draw.XRight
				p.X.Tick.Label.YAlign = draw.YCenter
				p.X.Tick.Label.Font.Size = 8
				p.Y.Tick.Label.Font.Size = 8
			} else {
				p.Y.Max = 1e9
				p.X.Max = 1e9
				p.X.Tick.Label.Font.Size = 1
				p.Y.Tick.Label.Font.Size = 1
			}

			plots[j][i] = p
		}
	}

	img := vgimg.New(vg.Points(300), vg.Points(350))
	dc := draw.New(img)

	canvases := table.Align(plots, dc)
	for j := 0; j < len(table.RowHeights); j++ {
		for i := 0; i < len(table.ColWidths); i++ {
			if plots[j][i] != nil {
				plots[j][i].Draw(canvases[j][i])
			}
		}
	}

	testFile := "testdata/tablealign.png"
	w, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}

	internal.TestImage(t, testFile)
}
