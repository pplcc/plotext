// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file is based on gonum.org/v1/plot/align.go which is
// Copyright ©2017 The gonum Authors. All rights reserved.

package plotext

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Table creates a table of subcanvases from a Canvas. In contrast to tiles of
// gonum.org/v1/plot the columns and rows of a table can have different widths
// and heights respectively.
type Table struct {
	// RowHeights specifies the number of rows their and relative heights
	// E. g.: When {2, 1} are used then the first row will get 66.7% and the
	// second row 33.3% of available height.
	RowHeights []float64
	// ColWidths specifies the number of columns and their relative widths
	ColWidths []float64
	// PadTop, PadBottom, PadRight, and PadLeft specify the padding
	// on the corresponding side of the table.
	PadTop, PadBottom, PadRight, PadLeft vg.Length
	// PadX and PadY specify the padding between columns and rows
	// of tiles respectively..
	PadX, PadY vg.Length
}

// At returns the subcanvas within c that corresponds to the
// cell at column x, row y, where 0, 0 is the upper, right corner
func (tab Table) At(c draw.Canvas, x, y int) draw.Canvas {
	// Canvas origin is left, bottom. Positive directions are right, up
	var sumColWidths float64
	for _, relColWidth := range tab.ColWidths {
		sumColWidths += relColWidth
	}

	var sumColWidthsRight float64
	for i, relColWidth := range tab.ColWidths {
		if i == x {
			break
		}
		sumColWidthsRight += relColWidth
	}

	var sumRowHeights float64
	for _, RowHeights := range tab.RowHeights {
		sumRowHeights += RowHeights
	}

	var sumRowHeightsAbove float64
	for i, RowHeights := range tab.RowHeights {
		if i == y {
			break
		}
		sumRowHeightsAbove += RowHeights
	}

	heightPerRelUnit := (c.Max.Y - c.Min.Y - tab.PadTop - tab.PadBottom -
		vg.Length(len(tab.RowHeights)-1)*tab.PadY) / vg.Length(sumRowHeights)
	widthPerRelUnit := (c.Max.X - c.Min.X - tab.PadLeft - tab.PadRight -
		vg.Length(len(tab.ColWidths)-1)*tab.PadX) / vg.Length(sumColWidths)

	ymax := c.Max.Y - tab.PadTop - vg.Length(y)*(tab.PadY) - vg.Length(sumRowHeightsAbove)*heightPerRelUnit
	ymin := ymax - vg.Length(tab.RowHeights[y])*heightPerRelUnit

	xmin := c.Min.X + tab.PadLeft + vg.Length(x)*(tab.PadX) + vg.Length(sumColWidthsRight)*widthPerRelUnit
	xmax := xmin + vg.Length(tab.ColWidths[x])*widthPerRelUnit

	return draw.Canvas{
		Canvas: vg.Canvas(c),
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: xmin, Y: ymin},
			Max: vg.Point{X: xmax, Y: ymax},
		},
	}
}

// Align returns a two-dimensional row-major array of Canvases which will
// produce plots with DataCanvases that are neatly aligned.
// The arguments to the function are a two-dimensional row-major array
// of plots and the canvas to which the plots are to be drawn.
func (tab Table) Align(plots [][]*plot.Plot, dc draw.Canvas) [][]draw.Canvas {
	o := make([][]draw.Canvas, len(plots))

	if len(plots) != len(tab.RowHeights) {
		panic(fmt.Sprintf("plot: plots rows (%v) != tiles rows (%v)", len(plots), tab.RowHeights))
	}

	// Create the initial tiles.
	for j := 0; j < len(tab.RowHeights); j++ {
		if len(plots[j]) != len(tab.ColWidths) {
			panic(fmt.Sprintf("plot: plots row %v columns (%v) != tiles columns (%v)", j, len(plots[j]), tab.RowHeights))
		}

		o[j] = make([]draw.Canvas, len(plots[j]))
		for i := 0; i < len(tab.ColWidths); i++ {
			o[j][i] = tab.At(dc, i, j)
		}
	}

	type posNeg struct {
		p, n float64 // x: n = left, p = right; y: n = bottom; p = top
	}
	xSpacing := make([]posNeg, len(tab.ColWidths))
	ySpacing := make([]posNeg, len(tab.RowHeights))

	// Calculate the maximum spacing between data canvases
	// for each row and column.
	for j, row := range plots {
		for i, p := range row {
			if p == nil {
				continue
			}
			c := o[j][i]
			dataC := p.DataCanvas(o[j][i])
			xSpacing[i].n = math.Max(float64(dataC.Min.X-c.Min.X), xSpacing[i].n)
			xSpacing[i].p = math.Max(float64(c.Max.X-dataC.Max.X), xSpacing[i].p)
			ySpacing[j].n = math.Max(float64(dataC.Min.Y-c.Min.Y), ySpacing[j].n)
			ySpacing[j].p = math.Max(float64(c.Max.Y-dataC.Max.Y), ySpacing[j].p)
		}
	}

	// Calculate the total row and column spacing.
	var xTotalSpace float64
	xTotalSpace = float64(tab.PadLeft+tab.PadRight) + float64(len(tab.ColWidths)-1)*float64(tab.PadX)
	for _, s := range xSpacing {
		xTotalSpace += s.n + s.p
	}
	var yTotalSpace float64
	yTotalSpace = float64(tab.PadTop+tab.PadBottom) + float64(len(tab.RowHeights)-1)*float64(tab.PadY)
	for _, s := range ySpacing {
		yTotalSpace += s.n + s.p
	}

	var sumColWidths float64
	for _, colWidth := range tab.ColWidths {
		sumColWidths += colWidth
	}

	var sumRowHeights float64
	for _, rowHeight := range tab.RowHeights {
		sumRowHeights += rowHeight
	}

	avgWidthPerUnit := vg.Length((float64(dc.Max.X-dc.Min.X) - xTotalSpace) / sumColWidths)
	avgHeightPerUnit := vg.Length((float64(dc.Max.Y-dc.Min.Y) - yTotalSpace) / sumRowHeights)

	moveVertical := make([]vg.Length, len(tab.ColWidths))
	for j := len(tab.RowHeights) - 1; j >= 0; j-- {
		row := plots[j]
		var moveHorizontal vg.Length
		for i, p := range row {
			c := o[j][i]

			if p != nil {
				dataC := p.DataCanvas(c)
				// Adjust the horizontal and vertical spacing between
				// canvases to match the maximum for each column and row,
				// respectively.
				c = draw.Crop(c,
					vg.Length(xSpacing[i].n)-(dataC.Min.X-c.Min.X),
					c.Max.X-dataC.Max.X-vg.Length(xSpacing[i].p),
					vg.Length(ySpacing[j].n)-(dataC.Min.Y-c.Min.Y),
					c.Max.Y-dataC.Max.Y-vg.Length(ySpacing[j].p),
				)
			}

			var width, height vg.Length
			if p == nil {
				width = c.Max.X - c.Min.X - vg.Length(xSpacing[i].p+xSpacing[i].n)
				height = c.Max.Y - c.Min.Y - vg.Length(ySpacing[j].p+ySpacing[j].n)
			} else {
				dataC := p.DataCanvas(c)
				width = dataC.Max.X - dataC.Min.X
				height = dataC.Max.Y - dataC.Min.Y
			}
			// Adjust the canvas so that the height of the DataCanvas
			// is the same for all plots in a row and the width of the
			// DataCanvas is the same for all plots in a column.
			o[j][i] = draw.Crop(c,
				moveHorizontal,
				moveHorizontal+avgWidthPerUnit*vg.Length(tab.ColWidths[i])-width,
				moveVertical[i],
				moveVertical[i]+avgHeightPerUnit*vg.Length(tab.RowHeights[j])-height,
			)
			moveHorizontal += avgWidthPerUnit*vg.Length(tab.ColWidths[i]) - width
			moveVertical[i] += avgHeightPerUnit*vg.Length(tab.RowHeights[j]) - height
		}
	}

	return o
}
