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

// DefaultCandleWidthFactor is the default width of the candle relative to the DefaultLineStyle.Width.
var DefaultCandleWidthFactor = 3

// Candlesticks implements the Plotter interface, drawing
// a bar plot of time, open, high, low, close tuples.
type Candlesticks struct {
	TOHLCVs

	// ColorUp is the color of sticks where C >= O
	ColorUp color.Color

	// ColorDown is the color of sticks where C < O
	ColorDown color.Color

	// LineStyle is the style used to draw the sticks.
	draw.LineStyle

	// CandleWidth is the width of a candlestick
	CandleWidth vg.Length

	// FixedLineColor determines if a fixed line color can be used for up and down bars.
	// When set to true the color of LineStyle is used to draw the sticks and
	// the borders of the candle. If set to false then ColorUp or ColorDown are used to
	// draw the sticks and the borders of the candle. Thus a candle's fill color is also
	// used for the borders and sticks.
	FixedLineColor bool
}

// NewCandlesticks creates as new candlestick plotter for
// the given data.
func NewCandlesticks(TOHLCV TOHLCVer) (*Candlesticks, error) {
	cpy, err := CopyTOHLCVs(TOHLCV)
	if err != nil {
		return nil, err
	}

	return &Candlesticks{
		TOHLCVs:        cpy,
		FixedLineColor: true,
		ColorUp:        color.RGBA{R: 128, G: 192, B: 128, A: 255}, // eye is more sensible to green
		ColorDown:      color.RGBA{R: 255, G: 128, B: 128, A: 255},
		LineStyle:      plotter.DefaultLineStyle,
		CandleWidth:    vg.Length(DefaultCandleWidthFactor) * plotter.DefaultLineStyle.Width,
	}, nil
}

// Plot implements the Plot method of the plot.Plotter interface.
func (sticks *Candlesticks) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	lineStyle := sticks.LineStyle

	for _, TOHLCV := range sticks.TOHLCVs {
		var fillColor color.Color
		if TOHLCV.C >= TOHLCV.O {
			fillColor = sticks.ColorUp
		} else {
			fillColor = sticks.ColorDown
		}

		if !sticks.FixedLineColor {
			lineStyle.Color = fillColor
		}
		// Transform the data
		// to the corresponding drawing coordinate.
		x := trX(TOHLCV.T)
		yh := trY(TOHLCV.H)
		yl := trY(TOHLCV.L)
		ymaxoc := trY(math.Max(TOHLCV.O, TOHLCV.C))
		yminoc := trY(math.Min(TOHLCV.O, TOHLCV.C))

		// top stick
		line := c.ClipLinesY([]vg.Point{{x, yh}, {x, ymaxoc}})
		c.StrokeLines(lineStyle, line...)

		// bottom stick
		line = c.ClipLinesY([]vg.Point{{x, yl}, {x, yminoc}})
		c.StrokeLines(lineStyle, line...)

		// body
		poly := c.ClipPolygonY([]vg.Point{{x - sticks.CandleWidth/2, ymaxoc}, {x + sticks.CandleWidth/2, ymaxoc}, {x + sticks.CandleWidth/2, yminoc}, {x - sticks.CandleWidth/2, yminoc}, {x - sticks.CandleWidth/2, ymaxoc}})
		c.FillPolygon(fillColor, poly)
		c.StrokeLines(lineStyle, poly)
	}
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (sticks *Candlesticks) DataRange() (xMin, xMax, yMin, yMax float64) {
	xMin = math.Inf(1)
	xMax = math.Inf(-1)
	yMin = math.Inf(1)
	yMax = math.Inf(-1)
	for _, TOHLCV := range sticks.TOHLCVs {
		xMin = math.Min(xMin, TOHLCV.T)
		xMax = math.Max(xMax, TOHLCV.T)
		yMin = math.Min(yMin, TOHLCV.L)
		yMax = math.Max(yMax, TOHLCV.H)
	}
	return
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
// We just return 2 glyph boxes at xmin, ymin and xmax, ymax
// Important is that they provide space for the left part of the first candle's body and for the right part of the last candle's body
func (sticks *Candlesticks) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, 2)

	xmin, xmax, ymin, ymax := sticks.DataRange()

	boxes[0].X = plt.X.Norm(xmin)
	boxes[0].Y = plt.Y.Norm(ymin)
	boxes[0].Rectangle = vg.Rectangle{
		Min: vg.Point{X: -(sticks.CandleWidth + sticks.LineStyle.Width) / 2, Y: 0},
		Max: vg.Point{X: 0, Y: 0},
	}

	boxes[1].X = plt.X.Norm(xmax)
	boxes[1].Y = plt.Y.Norm(ymax)
	boxes[1].Rectangle = vg.Rectangle{
		Min: vg.Point{X: 0, Y: 0},
		Max: vg.Point{X: +(sticks.CandleWidth + sticks.LineStyle.Width) / 2, Y: 0},
	}

	return boxes
}
