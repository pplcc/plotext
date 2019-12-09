// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package custplotter

import "gonum.org/v1/plot/plotter"

// TOHLCVMAer wraps the Len and TOHLCVMA methods. In addition to the all data contained by
// the basic TOHLCV types it provides moving average value as well
type TOHLCVMAer interface {
	// Len returns the number of time, open, high, low, close, volume tuples.
	Len() int

	// TOHLCVMA returns an time, open, high, low, close, volume and moving average tuple.
	TOHLCVMA(int) (float64, float64, float64, float64, float64, float64, float64)
}

// TOHLCVMAs implements the TOHLCVMAer interface using a slice.
type TOHLCVMAs []struct{ T, O, H, L, C, V, MA float64 }

// Len implements the Len method of the TOHLCVMAer interface.
func (TOHLCVMA TOHLCVMAs) Len() int {
	return len(TOHLCVMA)
}

// TOHLCVMA implements the TOHLCVMA method of the TOHLCVMAer interface.
func (TOHLCVMA TOHLCVMAs) TOHLCVMA(i int) (float64, float64, float64, float64, float64, float64, float64) {
	return TOHLCVMA[i].T, TOHLCVMA[i].O, TOHLCVMA[i].H, TOHLCVMA[i].L, TOHLCVMA[i].C, TOHLCVMA[i].V, TOHLCVMA[i].MA
}

// CopyTOHLCVMAs copies an TOHLCVMAer.
func CopyTOHLCVMAs(data TOHLCVMAer) (TOHLCVMAs, error) {
	cpy := make(TOHLCVMAs, data.Len())
	for i := range cpy {
		cpy[i].T, cpy[i].O, cpy[i].H, cpy[i].L, cpy[i].C, cpy[i].V, cpy[i].MA = data.TOHLCVMA(i)
		if err := plotter.CheckFloats(cpy[i].O, cpy[i].H, cpy[i].L, cpy[i].C, cpy[i].V); err != nil {
			return nil, err
		}
	}
	return cpy, nil
}
