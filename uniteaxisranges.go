// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotext

import (
	"math"

	"gonum.org/v1/plot"
)

// UniteAxisRanges sets the range of all axises to the minimum and the maximum of all axises.
func UniteAxisRanges(axises []*plot.Axis) {
	min := math.MaxFloat64
	max := -math.MaxFloat64

	for _, axis := range axises {
		min = math.Min(axis.Min, min)
		max = math.Max(axis.Max, max)
	}

	for _, axis := range axises {
		axis.Min = min
		axis.Max = max
	}
}
