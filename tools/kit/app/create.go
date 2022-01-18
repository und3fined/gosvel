/**
 * File: create.go
 * Project: app
 * File Created: 18 Jan 2022 15:06:47
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 18 Jan 2022 17:18:40
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package app

import (
	"gosvel/tools/kit/manifest"
	"gosvel/tools/utils/filepath"
	"io/ioutil"
	"path"
)

func Create(cwd, output string, manifestData *manifest.ManifestData) error {
	base := filepath.Relative(cwd, cwd, output)

	result := generateClientManifest(cwd, manifestData, base)

	manifestJS := path.Join(cwd, output, "manifest.js")
	if err := ioutil.WriteFile(manifestJS, []byte(result), 0644); err != nil {
		return err
	}

	return nil
}
