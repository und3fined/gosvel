/**
 * File: manifest.go
 * Project: create_manifest
 * File Created: 11 Jan 2022 20:39:41
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 18 Jan 2022 14:29:39
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package manifest

import (
	"encoding/json"
	"fmt"
	"gosvel/tools/utils/filepath"
	"log"
	"os"
)

type Manifest struct {
	opts *options

	defaultLayout string
	defaultError  string

	components []string
	routes     []RouteData
}

func (m *Manifest) init(opts ...Option) {
	if m.opts == nil {
		m.opts = &options{}
	}

	for _, op := range opts {
		op(m.opts)
	}
}

func (m *Manifest) initDefault() {
	m.defaultLayout = m.defaultComp("components/layout.svelte")
	m.defaultError = m.defaultComp("components/error.svelte")
}

func (m *Manifest) Create(opts ...Option) error {
	m.init(opts...)

	m.initDefault()

	routes := m.opts.Conf.Kit.Files.Routes

	base := filepath.Relative(m.opts.Cwd, m.opts.Cwd, routes)
	layoutPage := m.findLayout("__layout", base, m.defaultLayout)
	layoutError := m.findLayout("__error", base, m.defaultError)

	m.components = append(m.components, layoutPage, layoutError)

	if err := m.walk(routes, [][]*RouteSegment{}, []string{}, []string{layoutPage}, []string{layoutError}); err != nil {
		return err
	}

	manifestContent, _ := json.MarshalIndent(m.routes, "", "  ")

	manifestPath := fmt.Sprintf("%s/manifest_app.json", m.opts.Cwd)
	log.Println("manifestPath", manifestPath)
	f, _ := os.Create(manifestPath)
	defer f.Close()

	n, _ := f.Write(manifestContent)
	log.Printf("wrote %d bytes\n", n)

	return nil
}

func New(option ...Option) *Manifest {
	manifest := &Manifest{}
	manifest.init(option...)

	return manifest
}
