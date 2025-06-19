// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//go:build unix

package yt

import (
	"os"
)


func isExecutable(_ string, stat os.FileInfo) bool {
	// On Unix systems, check if executable bit is set (user, group, or others).
	return stat.Mode().Perm()&0o100 != 0 || stat.Mode().Perm()&0o010 != 0 || stat.Mode().Perm()&0o001 != 0
}
