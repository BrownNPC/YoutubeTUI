//go:build windows

package yt

// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

import (
	"debug/pe"
	"os"
)


func isExecutable(path string, _ os.FileInfo) bool {
	// Try to parse as PE (Portable Executable) format.
	f, err := pe.Open(path)
	if err != nil {
		return false
	}
	f.Close()
	return true
}
