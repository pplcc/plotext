// Copyright ©2018 Peter Paolucci. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

func compareFiles(path1, path2 string) bool {
	file1, err := ioutil.ReadFile(path1)

	if err != nil {
		log.Fatal(err)
	}

	file2, err := ioutil.ReadFile(path2)

	if err != nil {
		log.Fatal(err)
	}

	return bytes.Equal(file1, file2)
}

// TestImage compares a <NAME>.png testFile with the <NAME>_goldend.png file
// and calls a test error if the do not have the same contents.
func TestImage(t *testing.T, testFile string) {
	if !compareFiles(testFile, strings.TrimSuffix(testFile, filepath.Ext(testFile))+"_golden.png") {
		t.Errorf("image mismatch for %s\n", testFile)
	}
}
