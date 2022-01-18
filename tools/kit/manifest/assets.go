/**
 * File: assets.go
 * Project: manifest
 * File Created: 18 Jan 2022 15:15:34
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 18 Jan 2022 16:14:27
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package manifest

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (m *Manifest) getAssets(cwd, dir string) []Asset {
	var assets []Asset

	assetsFolder := path.Join(cwd, dir)
	if _, err := os.Stat(assetsFolder); os.IsNotExist(err) {
		return assets
	}

	if err := filepath.Walk(assetsFolder, func(curr string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			assetPath := strings.Replace(curr, cwd, "./", -1)
			assets = append(assets, Asset{
				File: path.Join(assetPath),
				Size: info.Size(),
				Type: "",
			})
		}
		return nil
	}); err != nil {
		return assets
	}

	return assets
}
