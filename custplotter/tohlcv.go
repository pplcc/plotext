// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package custplotter

import "gonum.org/v1/plot/plotter"

// TOHLCV wraps the Len and TOHLCV methods.
type TOHLCVer interface {
	// Len returns the number of time, open, high, low, close, volume tuples.
	Len() int

	// TOHLCV returns an time, open, high, low, close, volume tuple.
	TOHLCV(int) (float64, float64, float64, float64, float64, float64)
}

// TOHLCVs implements the TOHLCVer interface using a slice.
type TOHLCVs []struct{ T, O, H, L, C, V float64 }

// Len implements the Len method of the TOHLCVer interface.
func (TOHLCV TOHLCVs) Len() int {
	return len(TOHLCV)
}

// TOHLCV implements the TOHLCV method of the TOHLCVer interface.
func (TOHLCV TOHLCVs) TOHLCV(i int) (float64, float64, float64, float64, float64, float64) {
	return TOHLCV[i].T, TOHLCV[i].O, TOHLCV[i].H, TOHLCV[i].L, TOHLCV[i].C, TOHLCV[i].V
}

// CopyTOHLCVs copies an TOHLCVer.
func CopyTOHLCVs(data TOHLCVer) (TOHLCVs, error) {
	cpy := make(TOHLCVs, data.Len())
	for i := range cpy {
		cpy[i].T, cpy[i].O, cpy[i].H, cpy[i].L, cpy[i].C, cpy[i].V = data.TOHLCV(i)
		if err := plotter.CheckFloats(cpy[i].O, cpy[i].H, cpy[i].L, cpy[i].C, cpy[i].V); err != nil {
			return nil, err
		}
	}
	return cpy, nil
}
