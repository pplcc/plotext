// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/pplcc/plotext/custplotter"
	"github.com/pplcc/plotext/examples"
)

func CreateTOHLCVTestData() custplotter.TOHLCVs {
	return examples.CreateTOHLCVExampleData(20)
}
