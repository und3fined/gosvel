/**
 * File: main.go
 * Project: kit
 * File Created: 12 Jan 2022 10:47:26
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 12 Jan 2022 17:44:23
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package main

import (
	"gosvel/tools/kit/config"
	"gosvel/tools/kit/manifest"
	"gosvel/tools/kit/utils/filepath"
	"log"
	"path"
)

func main() {
	cwd := filepath.CWD()
	buildDir := path.Join(cwd, ".gsvel/build")

	conf := config.Config{
		Extensions: []string{".svelte"},
		Kit: config.KitConfig{
			Files: config.KitFiles{
				Assets: "./themes/default/assets",
				Routes: "./themes/default/routes",
			},
		},
	}

	manifestInstance := manifest.New(manifest.Config(conf), manifest.Output(buildDir), manifest.Cwd(cwd))

	if err := manifestInstance.Create(); err != nil {
		log.Printf("Err : %+v", err)
	}
}
