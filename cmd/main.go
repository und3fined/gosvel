/**
 * File: main.go
 * Project: cmd
 * File Created: 31 Dec 2021 15:38:08
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 10 Jan 2022 21:59:29
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2021 und3fined.com
 */
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/segmentio/encoding/json"
	"rogchap.com/v8go"
)

type CompiledContent struct {
	Code *string `json:"code"`
	Map  string  `json:"map"`
}

type CompileResult struct {
	JS  CompiledContent `json:"js"`
	CSS CompiledContent `json:"css"`
	// Warnings string `json:"warnings"`
}

func main() {

	sveltePlugin := api.Plugin{
		Name: "svelte",
		Setup: func(build api.PluginBuild) {
			svelte, _ := ioutil.ReadFile(filepath.Join("./", "internal/node_modules/svelte/compiler.js"))
			svelteStr := strings.Replace(string(svelte), "self.performance.now();", "{};", 1)
			svelteStr = strings.Replace(svelteStr, "const Url$1 = (typeof URL !== 'undefined' ? URL : require('url').URL);", "", 1)

			iso := v8go.NewIsolate()
			v8Ctx := v8go.NewContext(iso)

			build.OnStart(func() (api.OnStartResult, error) {
				if _, err := v8Ctx.RunScript(svelteStr, "svelte_compiler"); err != nil {
					log.Printf("[svelte_compiler] Error: %+v", err)
					return api.OnStartResult{}, err
				}

				svelteVersion, _ := v8Ctx.RunScript("svelte.VERSION", "svelte_version")
				fmt.Printf("svelte %s\n", svelteVersion.String())

				return api.OnStartResult{}, nil
			})

			build.OnLoad(api.OnLoadOptions{Filter: `\.svelte$`}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				text, err := ioutil.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				pwd, _ := filepath.Abs("./")
				filename := strings.Replace(args.Path, pwd+"/", "", 1)

				execCompileStr := fmt.Sprintf("svelte.compile(`%s`, { filename: '%s', generate: 'ssr', hydratable: true })", text, filename)
				scriptOrigin := fmt.Sprintf("compile_%s", filename)
				result, err := v8Ctx.RunScript(execCompileStr, scriptOrigin)
				if err != nil {
					return api.OnLoadResult{}, err
				} else {
					if result.IsObject() {
						jsonStr, _ := v8go.JSONStringify(v8Ctx, result)

						var compiledResult CompileResult
						json.NewDecoder(strings.NewReader(jsonStr)).Decode(&compiledResult)

						return api.OnLoadResult{
							Contents: compiledResult.JS.Code,
							Loader:   api.LoaderText,
						}, nil
					}
				}

				return api.OnLoadResult{}, err
			})
		},
	}

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{"./themes/default/routes/main.js"},
		Outfile:     "output.js",
		Bundle:      true,
		Write:       true,
		LogLevel:    api.LogLevelInfo,
		Plugins:     []api.Plugin{sveltePlugin},
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	// ctx := compile.New()

	// ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
	// ctx.RunScript("const result = add(3, 4)", "main.js")    // any functions previously added to the context can be called
	// val, _ := ctx.RunScript("result", "value.js")           // return a value in JavaScript back to Go
	// fmt.Printf("addition result: %s\n", val)
}
