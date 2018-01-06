// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

import (
	"math"
	"math/rand"
	"time"

	"github.com/pplcc/plotext/custplotter"
)

// CreateTOHLCVExampleData generates and returns some artificial TOHLCV data for testing and demo purpose
func CreateTOHLCVExampleData(n int) custplotter.TOHLCVs {
	rnd := rand.New(rand.NewSource(1))
	m := 4 * n
	fract := make([]float64, m)
	for i := 0; i < m; i++ {
		fract[i] = 100
	}
	stat1 := 0.0
	stat2 := 0.0
	for k := m; k > 0; k = k / 2 {
		j := 0
		for i := 0; i < m; i++ {
			if j == 0 {
				j = k
				stat2 = stat1
				stat1 = 10.0 * (float64(k)/float64(m) + 0.02) * (2.0*rnd.Float64() - 1.0)
			}
			fract[i] += float64(k-j)/float64(k)*stat1 + float64(j)/float64(k)*stat2
			j--
		}
	}

	data := make(custplotter.TOHLCVs, n)

	loc, _ := time.LoadLocation("America/New_York")
	for i := range data {
		data[i].T = float64(time.Date(2000, 01, 02, 03, 04, 05, 0, loc).Add(time.Duration(i) * time.Minute).Unix())
		data[i].O = fract[4*i]
		data[i].H = math.Max(math.Max(fract[4*i], fract[4*i+1]), math.Max(fract[4*i+2], fract[4*i+3]))
		data[i].L = math.Min(math.Min(fract[4*i], fract[4*i+1]), math.Min(fract[4*i+2], fract[4*i+3]))
		data[i].C = fract[4*i+3]

		data[i].V = (data[i].H - data[i].L + math.Abs(data[i].C-data[i].O)) * 100 // just use this as a fake volume
	}
	return data
}
