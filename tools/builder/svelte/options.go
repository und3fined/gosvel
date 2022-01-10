/**
 * File: options.go
 * Project: svelte
 * File Created: 10 Jan 2022 22:13:16
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 10 Jan 2022 22:19:22
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package svelte

import "github.com/evanw/esbuild/pkg/api"

type EsOptions struct {
	api.BuildOptions
}

type CompileOption struct {
	Format     string `json:"format,omitempty"`
	Name       string `json:"name,omitempty"`
	Filename   string `json:"filename,omitempty"`
	Generate   string `json:"generate,omitempty"`
	ErrorMode  string `json:"errorMode,omitempty"`
	Dev        bool   `json:"dev,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Hydratable bool   `json:"hydratable"`

	// format?: ModuleFormat;
	//   name?: string;
	//   filename?: string;
	//   generate?: 'dom' | 'ssr' | false;
	//   errorMode?: 'throw' | 'warn';
	//   varsReport?: 'full' | 'strict' | false;
	//   sourcemap?: object | string;
	//   enableSourcemap?: EnableSourcemap;
	//   outputFilename?: string;
	//   cssOutputFilename?: string;
	//   sveltePath?: string;
	//   dev?: boolean;
	//   accessors?: boolean;
	//   immutable?: boolean;
	//   hydratable?: boolean;
	//   legacy?: boolean;
	//   customElement?: boolean;
	//   tag?: string;
	//   css?: boolean;
	//   loopGuardTimeout?: number;
	//   namespace?: string;
	//   cssHash?: CssHashGetter;
	//   preserveComments?: boolean;
	//   preserveWhitespace?: boolean;
}
