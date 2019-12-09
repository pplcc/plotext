// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/pplcc/plotext/custplotter"
	"github.com/pplcc/plotext/examples"
)

func CreateTOHLCVTestData(n int) custplotter.TOHLCVs {
	return examples.CreateTOHLCVExampleData(n)
}

func CreateTOHLCVMATestData(n, window int) custplotter.TOHLCVMAs {
	src := examples.CreateTOHLCVExampleData(n)
	dst := make(custplotter.TOHLCVMAs, len(src))
	for i := 0; i < len(src); i++ {
		dst[i].T = src[i].T
		dst[i].O = src[i].O
		dst[i].H = src[i].H
		dst[i].L = src[i].L
		dst[i].C = src[i].C
		dst[i].V = src[i].V

		if i < window {
			dst[i].MA = src[i].C
		} else {
			// simple non-weighted averaging
			var acc float64
			for j := i - window; j < i; j++ {
				acc += src[j].C
			}
			dst[i].MA = acc / float64(window)
		}
	}
	return dst
}
