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

// VBars implements the Plotter interface, drawing
// a volume bar plot.
type VBars struct {
	TOHLCVs

	// ColorUp is the color of bars where C >= O
	ColorUp color.Color

	// ColorDown is the color of bars where C < O
	ColorDown color.Color

	// LineStyle is the style used to draw the bars.
	draw.LineStyle
}

// NewBars creates as new bar plotter for
// the given data.
func NewVBars(TOHLCV TOHLCVer) (*VBars, error) {
	cpy, err := CopyTOHLCVs(TOHLCV)
	if err != nil {
		return nil, err
	}

	minDeltaT := 0.0
	if len(cpy) > 1 {
		minDeltaT = math.MaxFloat64
		for i := 0; i < len(cpy)-1; i++ {
			minDeltaT = math.Min(minDeltaT, cpy[i+1].T-cpy[i].T)
		}
	}

	return &VBars{
		TOHLCVs:   cpy,
		ColorUp:   color.RGBA{R: 0, G: 128, B: 0, A: 255}, // eye is more sensible to green
		ColorDown: color.RGBA{R: 196, G: 0, B: 0, A: 255},
		LineStyle: plotter.DefaultLineStyle,
	}, nil
}

// Plot implements the Plot method of the plot.Plotter interface.
func (bars *VBars) Plot(c draw.Canvas, plt *plot.Plot) {
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
		y0 := trY(0)
		y := trY(TOHLCV.V)

		bar := c.ClipLinesY([]vg.Point{{X: x, Y: y0}, {X: x, Y: y}})
		c.StrokeLines(lineStyle, bar...)

	}
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (bars *VBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymin = 0
	ymax = math.Inf(-1)
	for _, TOHLCV := range bars.TOHLCVs {
		xmin = math.Min(xmin, TOHLCV.T)
		xmax = math.Max(xmax, TOHLCV.T)
		ymax = math.Max(ymax, TOHLCV.V)
	}
	return
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
// We just return 2 glyph boxes at xmin, ymin and xmax, ymax
// Important is that they provide space for the left part of the first candle's body and for the right part of the last candle's body
func (bars *VBars) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	xmin, xmax, ymin, ymax := bars.DataRange()

	boxes[0].X = plt.X.Norm(xmin)
	boxes[0].Y = plt.Y.Norm(ymin)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -bars.LineStyle.Width / 2, Y: 0},
		Max: vg.Point{X: 0, Y: 0},
	}

	boxes[1].X = plt.X.Norm(xmax)
	boxes[1].Y = plt.Y.Norm(ymax)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{X: 0, Y: 0},
		Max: vg.Point{X: +bars.LineStyle.Width / 2, Y: 0},
	}

	return boxes
}
