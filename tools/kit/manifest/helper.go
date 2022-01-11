/**
 * File: helper.go
 * Project: manifest
 * File Created: 11 Jan 2022 21:48:31
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 11 Jan 2022 23:51:13
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package manifest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"gosvel/tools/kit/config"
)

var specials = []string{"__layout", "__layout.reset", "__error"}
var reTmpFile = regexp.MustCompile(`^(\.[a-z0-9]+)+$`)

func (m *Manifest) findLayout()

func findLayout(fileName, dir string, extensions []string, defaultLayout string) string {
	for _, ext := range extensions {
		filePath := path.Join(dir, fmt.Sprintf("%s.%s", fileName, ext))

		if _, err := os.Stat(filePath); err == nil {
			return filePath
		}
	}

	return defaultLayout
}

func walk(conf config.Config, cwd string, dir string, parentSegments [][]RouteSegment, parentParams []string, layoutStack []string, errorStack []string) error {
	items := []WalkItem{}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fInfo := range files {
		basename := fInfo.Name()
		file := strings.Replace(path.Join(dir, basename), cwd, "", 1)
		isDir := fInfo.IsDir()

		ext := path.Ext(basename)
		if !isDir {
			ext = findString(conf.Extensions, func(item string) bool {
				return strings.HasSuffix(basename, item)
			})
		}

		name := basename
		if ext != "" {
			name = strings.Replace(name, ext, "", 1)
			name = strings.TrimSuffix(name, ".")
		}

		// detect __ file
		if name[0:1] == "_" {
			if name[1:2] == "_" && !contains(specials, name) {
				return errors.New(fmt.Sprintf("Files and directories prefixed with __ are reserved (saw %s)", file.Name()))
			}

			return nil
		}

		// detect hidden file
		if basename[0:1] == "." && basename != ".well-known" {
			return nil
		}

		// filter out tmp files etc
		if !isDir && !reTmpFile.MatchString(ext) {
			return nil
		}

		segment := name
		if isDir {
			segment = basename
		}

		// parts := getParts(segment, file)
		isIndex := strings.HasPrefix(basename, "index.")
		if isDir {
			isIndex = false
		}

		dotIndex := strings.Index(basename, ".")
		endIndex := len(basename) - len(ext)

		items = append(items, WalkItem{
			Basename:    basename,
			Ext:         ext,
			Parts:       []RouteSegment{},
			File:        file,
			IsDir:       isDir,
			IsIndex:     isIndex,
			IsPage:      contains(conf.Extensions, ext),
			RouteSuffix: basename[dotIndex:endIndex],
		})
	}
}

func getParts(fileName, filePath) {

}

func contains(arr []string, item string) bool {
	state := false

	for _, i := range arr {
		if i == item {
			return true
		}
	}

	return state
}

func findString(arr []string, fn func(string) bool) string {
	for _, item := range arr {
		if fn(item) {
			return item
		}
	}

	return ""
}
