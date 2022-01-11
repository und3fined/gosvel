/**
 * File: manifest.go
 * Project: create_manifest
 * File Created: 11 Jan 2022 20:39:41
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 11 Jan 2022 23:50:13
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package manifest

import (
	"path"

	"gosvel/tools/kit/config"
	"gosvel/tools/kit/utils/filepath"
)

type Manifest struct {
	opts *options

	components []string
	routes     []string
}

func (m *Manifest) init(opts ...Option) {
	if m.opts == nil {
		m.opts = &options{}
	}

	for _, op := range opts {
		op(m.opts)
	}
}

func (m *Manifest) Create(opts ...Option) {
	m.init(opts...)

}

func New(option ...Option) *Manifest {
	manifest := &Manifest{}
	manifest.init(option...)

	return manifest
}

// CreateManifest
func CreateManifest(conf config.Config, output, cwd string) ManifestData {
	if cwd == "" {
		cwd = filepath.CWD()
	}

	var components []string
	var routes []RouteData

	defaultLayout := path.Join(cwd, output, "components/layout.svelte")
	defaultError := path.Join(cwd, output, "components/error.svelte")

	base := path.Join(cwd, config.Kit.Files.Routes)
	layoutPage := findLayout("__layout", base, config.Extensions, defaultLayout)
	layoutError := findLayout("__error", base, config.Extensions, defaultError)

	components = append(components, layoutPage, layoutError)

	walk(base, [][]RouteSegment{}, []string{}, []string{layoutPage}, []string{layoutError})

	return ManifestData{}
}
