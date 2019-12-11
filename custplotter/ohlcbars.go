// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package custplotter

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// DefaultTickWidth is the default width of the open and close ticks.
var DefaultTickWidth = vg.Points(2)

// OHLCBars implements the Plotter interface, drawing
// a bar plot of time, open, high, low, close tuples.
type OHLCBars struct {
	TOHLCVs

	// ColorUp is the color of bars where C >= O
	ColorUp color.Color

	// ColorDown is the color of bars where C < O
	ColorDown color.Color

	// LineStyle is the style used to draw the bars.
	draw.LineStyle

	// CapWidth is the width of the caps drawn at the top
	// of each error bar.
	TickWidth vg.Length
}

// NewBars creates as new bar plotter for
// the given data.
func NewOHLCBars(TOHLCV TOHLCVer) (*OHLCBars, error) {
	cpy, err := CopyTOHLCVs(TOHLCV)
	if err != nil {
		return nil, err
	}

	return &OHLCBars{
		TOHLCVs:   cpy,
		ColorUp:   color.RGBA{R: 0, G: 128, B: 0, A: 255}, // eye is more sensible to green
		ColorDown: color.RGBA{R: 196, G: 0, B: 0, A: 255},
		LineStyle: plotter.DefaultLineStyle,
		TickWidth: DefaultTickWidth,
	}, nil
}

// Plot implements the Plot method of the plot.Plotter interface.
func (bars *OHLCBars) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	lineStyle := bars.LineStyle

	for _, TOHLCV := range bars.TOHLCVs {
		if TOHLCV.C >= TOHLCV.O {
			lineStyle.Color = bars.ColorUp
		} else {
			lineStyle.Color = bars.ColorDown
		}

		// Transform the data
		// to the corresponding drawing coordinate.
		x := trX(TOHLCV.T)
		yo := trY(TOHLCV.O)
		yh := trY(TOHLCV.H) // + vg.Length(bars.LineStyle.Width/2.0)
		yl := trY(TOHLCV.L) // - vg.Length(bars.LineStyle.Width/2.0)
		yc := trY(TOHLCV.C)

		bar := c.ClipLinesY([]vg.Point{{X: x, Y: yl}, {X: x, Y: yh}})
		c.StrokeLines(lineStyle, bar...)

		if c.Contains(vg.Point{X: x, Y: yo}) {
			c.StrokeLine2(lineStyle, x, yo, x-bars.TickWidth, yo)
		}

		if c.Contains(vg.Point{X: x, Y: yc}) {
			c.StrokeLine2(lineStyle, x, yc, x+bars.TickWidth, yc)
		}

	}
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (bars *OHLCBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymin = math.Inf(1)
	ymax = math.Inf(-1)
	for _, TOHLCV := range bars.TOHLCVs {
		xmin = math.Min(xmin, TOHLCV.T)
		xmax = math.Max(xmax, TOHLCV.T)
		ymin = math.Min(ymin, TOHLCV.L)
		ymax = math.Max(ymax, TOHLCV.H)
	}
	return
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
// We just return 2 glyph boxes at xmin, ymin and xmax, ymax
// Important is that they provide space for the first open tick and the last close tick
func (bars *OHLCBars) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	xmin, xmax, ymin, ymax := bars.DataRange()

	boxes[0].X = plt.X.Norm(xmin)
	boxes[0].Y = plt.Y.Norm(ymin)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -bars.TickWidth, Y: 0},
		Max: vg.Point{X: 0, Y: 0},
	}

	boxes[1].X = plt.X.Norm(xmax)
	boxes[1].Y = plt.Y.Norm(ymax)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{X: 0, Y: 0},
		Max: vg.Point{X: +bars.TickWidth, Y: 0},
	}

	return boxes
}
