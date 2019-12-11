// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package custplotter_test

import (
	"log"
	"testing"

	"gonum.org/v1/plot"

	"github.com/pplcc/plotext/custplotter"
	"github.com/pplcc/plotext/custplotter/internal"
)

func TestNewCandlesticks(t *testing.T) {
	t.SkipNow() // test is broken in upstream

	testTOHLCVs := internal.CreateTOHLCVTestData(20)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	bars, err := custplotter.NewCandlesticks(testTOHLCVs)
	if err != nil {
		log.Panic(err)
	}

	p.Add(bars)

	testFile := "testdata/candlesticks.png"
	err = p.Save(180, 100, testFile)
	if err != nil {
		log.Panic(err)
	}

	internal.TestImage(t, testFile)
}

func TestNewCandlesticksWithMovingAverage(t *testing.T) {
	window := 5
	testTOHLCVMAs := internal.CreateTOHLCVMATestData(100, window)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}

	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	bars, err := custplotter.NewCandlesticksWithMovingAverage(testTOHLCVMAs, window)
	if err != nil {
		log.Panic(err)
	}

	p.Add(bars)

	testFile := "testdata/candlesticks_with_moving_average.png"
	err = p.Save(640, 480, testFile)
	if err != nil {
		log.Panic(err)
	}
}
