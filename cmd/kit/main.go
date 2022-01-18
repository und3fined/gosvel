/**
 * File: main.go
 * Project: kit
 * File Created: 12 Jan 2022 10:47:26
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 18 Jan 2022 16:19:36
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package main

import (
	"log"
	"os"
	"path"

	"gosvel/tools/kit/app"
	"gosvel/tools/kit/config"
	"gosvel/tools/kit/manifest"
)

func main() {
	cwd, _ := os.Getwd()
	buildDir := path.Join(".gsvel/build")

	log.Println("cwd", cwd)

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

	if data, err := manifestInstance.Create(); err != nil {
		log.Printf("Err : %+v", err)
	} else {
		err := app.Create(cwd, buildDir, data)
		log.Println("Err", err)
	}
}
