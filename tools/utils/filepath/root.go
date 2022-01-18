/**
 * File: root.go
 * Project: filepath
 * File Created: 11 Jan 2022 21:40:17
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 18 Jan 2022 14:26:26
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package filepath

import (
	"os"
	"path/filepath"
	"strings"
)

func CWD() string {
	ex, _ := os.Executable()

	if strings.Contains(ex, "go-build") {
		ex2, _ := os.Getwd()
		return filepath.Dir(ex2)
	}

	return filepath.Dir(ex)
}

func Relative(cwd, from, to string) string {
	to = filepath.Join(cwd, to)
	nextPath, _ := filepath.Rel(from, to)
	return nextPath
}
