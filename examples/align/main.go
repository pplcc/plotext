// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/pplcc/plotext"
	"github.com/pplcc/plotext/custplotter"
	"github.com/pplcc/plotext/examples"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func main() {
	// This simple example creates a candlestick plot above a volume plot.
	// The candlesticks plot gets 2/3 of the available height and the volume
	// bars plot gets 1/3.
	// The x-axises use the same scale and are aligned, even if the labels
	// on one of the y-axises requires more space.
	//
	// Note: AlignTest in align_test.go creates a more complex page
	// (see testdata/tablealign_golden.png)

	// create some fake data

	n := 60
	fakeTOHLCVs := examples.CreateTOHLCVExampleData(n)

	// create the candlesticks plot

	p1, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p1.Title.Text = "Candlesticks and Volume Bars"
	// p.X.Label.Text = "Time"
	p1.Y.Label.Text = "Price"
	p1.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	candlesticks, err := custplotter.NewCandlesticks(fakeTOHLCVs)
	if err != nil {
		log.Panic(err)
	}

	p1.Add(candlesticks)

	// create the volume bars plot

	p2, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p2.X.Label.Text = "Time"
	p2.Y.Label.Text = "Volume"
	p2.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	vBars, err := custplotter.NewVBars(fakeTOHLCVs)
	if err != nil {
		log.Panic(err)
	}

	// The x axis are perfectly aligned, but due to the different width of a
	// candle and a volume bar the horizontal position of the y-axis is
	// little different. If you want to compensate for this then you can add
	// some padding to the y-axis of the volume plot. To do so uncomment the
	// following line:
	//
	// p2.Y.Padding += (candlesticks.CandleWidth - vBars.LineStyle.Width) / 2

	p2.Add(vBars)

	// it is not really required with this example data, but let's
	// make sure that the x axises have the same range anyway
	plotext.UniteAxisRanges([]*plot.Axis{&p1.X, &p2.X})

	// create a table with one column and two rows
	table := plotext.Table{
		RowHeights: []float64{2, 1}, // 2/3 for candlesticks and 1/3 for volume bars
		ColWidths:  []float64{1},
	}

	// see align_test.go for another example on how to construct this structure using loops
	plots := [][]*plot.Plot{[]*plot.Plot{p1}, []*plot.Plot{p2}}

	img := vgimg.New(450, 300)
	dc := draw.New(img)

	canvases := table.Align(plots, dc)
	plots[0][0].Draw(canvases[0][0])
	plots[1][0].Draw(canvases[1][0])

	testFile := "align.png"
	w, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}
}
