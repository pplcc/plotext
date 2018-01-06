// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/pplcc/plotext/custplotter"
	"github.com/pplcc/plotext/examples"
	"gonum.org/v1/plot"
)

func main() {
	n := 60
	fakeTOHLCVs := examples.CreateTOHLCVExampleData(n)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	p.Title.Text = "Candlesticks"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Price"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	bars, err := custplotter.NewCandlesticks(fakeTOHLCVs)
	if err != nil {
		log.Panic(err)
	}

	p.Add(bars)

	err = p.Save(450, 200, "candlesticks.png")
	if err != nil {
		log.Panic(err)
	}
}
