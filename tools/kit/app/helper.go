/**
 * File: helper.go
 * Project: app
 * File Created: 18 Jan 2022 15:07:54
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 18 Jan 2022 17:20:26
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package app

import (
	"fmt"
	"gosvel/tools/kit/manifest"
	"gosvel/tools/utils/filepath"
	"log"
	"regexp"
	"strings"
)

var reReplaceComp = regexp.MustCompile(`/^\t/`)

func generateClientManifest(cwd string, manifestData *manifest.ManifestData, base string) string {
	var compTmp []string
	var componentIndexes = make(map[string]int)

	for i, comp := range manifestData.Components {
		componentIndexes[comp] = i

		log.Printf("comp: %s", comp)

		genImport := fmt.Sprintf("() => import('%s')", getPath(cwd, base, comp))
		compTmp = append(compTmp, genImport)
	}

	components := strings.Join(compTmp, ",\n\t\t\t\t")
	components = fmt.Sprintf("[%s]", components)
	components = reReplaceComp.ReplaceAllString(components, "")

	var routeTmp []string
	for _, route := range manifestData.Routes {
		if route.Type == "page" {
			var paramsX []string
			if len(route.Params) > 0 {

				for i, param := range route.Params {
					str := fmt.Sprintf("%s: d(m[%d])", param, i+1)
					paramsX = append(paramsX, str)
				}
			}

			params := ""
			if len(paramsX) > 0 {
				params = fmt.Sprintf("(m) => ({%s})", strings.Join(paramsX, ", "))
			}

			var tuple []string
			tuple = append(tuple, route.Pattern)
			tuple = append(tuple, getIndices(route.C, componentIndexes))
			tuple = append(tuple, getIndices(route.E, componentIndexes))

			if params != "" {
				tuple = append(tuple, params)
			}

			cComp := route.C[len(route.C)-1]
			tupleStr := strings.Join(tuple, ", ")
			routeTmp = append(routeTmp, fmt.Sprintf("// %s\n\t\t[%s]", cComp, tupleStr))
		}
	}

	routes := strings.Join(routeTmp, ",\n\n\t\t")
	routes = fmt.Sprintf("[%s]", routes)

	return fmt.Sprintf(`
  const c = %s;
  const d = decodeURIComponent;
  export const routes = %s;

  // we import the root layout/error components eagerly, so that
	// connectivity errors after initialisation don't nuke the app
	export const fallback = [c[0](), c[1]()];
  `, components, routes)
}

func getPath(cwd, base, component string) string {
	return filepath.Relative(cwd, base, component)
}

func getIndices(parts []string, indexes map[string]int) string {
	var result []string

	for _, part := range parts {
		if part != "" {
			result = append(result, fmt.Sprintf("c[%d]", indexes[part]))
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(result, ", "))
}
