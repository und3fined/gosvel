/**
 * File: helper.go
 * Project: manifest
 * File Created: 11 Jan 2022 21:48:31
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 12 Jan 2022 22:42:17
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package manifest

import (
	"fmt"
	"gosvel/tools/kit/utils/filepath"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var specials = []string{"__layout", "__layout.reset", "__error"}
var reTmpFile = regexp.MustCompile(`^(\.[a-z0-9]+)+$`)
var rePart = regexp.MustCompile(`\[(.+?\(.+?\)|.+?)\]`)
var reRestContent = regexp.MustCompile(`^\.{3}.+$`)
var reParseContent = regexp.MustCompile(`([^(]+)$`)
var reValidContent = regexp.MustCompile(`^(\.\.\.)?[a-zA-Z0-9_$]+$`)

func (m *Manifest) defaultComp(fileComp string) string {
	if path.IsAbs(m.opts.Output) {
		return path.Join(m.opts.Output, fileComp)
	}

	return path.Join(m.opts.Cwd, m.opts.Output, fileComp)
}

func (m *Manifest) findLayout(fileName, dir string, defaultLayout string) string {
	for _, ext := range m.opts.Conf.Extensions {
		filePath := path.Join(dir, fmt.Sprintf("%s.%s", fileName, ext))
		if _, err := os.Stat(filePath); err == nil {
			return filePath
		}
	}

	return defaultLayout
}

func (m *Manifest) walk(dir string, parentSegments [][]*RouteSegment, parentParams []string, layoutStack []string, errorStack []string) error {
	log.Printf("\n\n---------\n\nStart walk: %s", dir)

	cwd := m.opts.Cwd
	extensions := m.opts.Conf.Extensions

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fInfo := range files {
		basename := fInfo.Name()

		log.Println("basename", basename)

		resolved := path.Join(dir, basename)
		file := filepath.Relative(cwd, resolved)
		isDir := fInfo.IsDir()

		ext := findString(extensions, func(item string) bool {
			return strings.HasSuffix(basename, item)
		})

		if ext == "" {
			ext = path.Ext(basename)
		}

		name := basename
		if ext != "" {
			name = strings.Replace(name, ext, "", 1)
		}

		// detect __ file
		if name[0:1] == "_" {
			if name[1:2] == "_" && !contains(specials, name) {
				return fmt.Errorf("Files and directories prefixed with __ are reserved (saw %s)", fInfo.Name())
			}

			continue
		}

		// detect hidden file
		if basename[0:1] == "." && basename != ".well-known" {
			continue
		}

		// filter out tmp files etc
		if !isDir && !reTmpFile.MatchString(ext) {
			continue
		}

		segment := basename
		if ext != "" {
			segment = basename[:len(basename)-len(ext)]
		}

		parts := m.getParts(segment, file)

		isIndex := strings.HasPrefix(basename, "index.")
		if isDir {
			isIndex = false
		}

		routeSuffix := ""
		dotIndex := strings.Index(basename, ".")
		if dotIndex > -1 {
			endIndex := len(basename) - len(ext)
			routeSuffix = basename[dotIndex:endIndex]
		}

		m.items = append(m.items, WalkItem{
			Basename:    basename,
			Ext:         ext,
			Parts:       parts,
			File:        file,
			IsDir:       isDir,
			IsIndex:     isIndex,
			IsPage:      contains(extensions, ext),
			RouteSuffix: routeSuffix,
		})
	}

	sort.Slice(m.items, m.comparator)

	for _, item := range m.items {
		segments := parentSegments

		log.Println("Item:", item.Basename, "index", item.IsIndex)

		if item.IsIndex {
			if item.RouteSuffix != "" {
				segmentLen := len(segments)
				if segmentLen > 0 {
					lastSegment := segments[segmentLen-1]
					lastPart := lastSegment[len(lastSegment)-1]

					if lastPart.Dynamic {
						lastSegment = append(lastSegment, &RouteSegment{
							Dynamic: false,
							Rest:    false,
							Content: item.RouteSuffix,
						})
					} else {
						lastSegment[len(lastSegment)-1] = &RouteSegment{
							Dynamic: false,
							Rest:    false,
							Content: fmt.Sprintf("%s%s", lastPart.Content, item.RouteSuffix),
						}
					}

					segments[segmentLen-1] = lastSegment
				} else {
					segments = append(segments, item.Parts)
				}
			}
		} else {
			segments = append(segments, item.Parts)
		}

		params := parentParams

		for _, part := range item.Parts {
			if part.Dynamic {
				params = append(params, part.Content)
			}
		}

		var simpleSegments []*RouteSegment

		for _, segment := range segments {
			dynamic := some(segment, func(part *RouteSegment) bool {
				return part.Dynamic
			})

			rest := some(segment, func(part *RouteSegment) bool {
				return part.Rest
			})

			simpleSegments = append(simpleSegments, &RouteSegment{
				Dynamic: dynamic,
				Rest:    rest,
				Content: generateContent(segment),
			})
		}

		if item.IsDir {
			layoutReset := m.findLayout("__layout.reset", item.File, "")
			layout := m.findLayout("__layout", item.File, "")
			layoutError := m.findLayout("__error", item.File, "")

			if layoutReset != "" && layout != "" {
				return fmt.Errorf("Cannot have __layout next to __layout.reset: %s", layoutReset)
			}

			if layout != "" {
				m.components = append(m.components, layout)
				layoutStack = append(layoutStack, layout)
			}

			if layoutError != "" {
				m.components = append(m.components, layoutError)
				errorStack = append(errorStack, layoutError)
			}

			if layoutReset != "" {
				m.components = append(m.components, layoutReset)
				layoutStack = []string{layoutReset}
				errorStack = []string{layoutError}
			}

			err := m.walk(
				path.Join(dir, item.Basename),
				segments,
				params,
				layoutStack,
				errorStack,
			)
			if err.Error() == "no such file or directory" {
				continue
			}
		} else if item.IsPage {
			m.components = append(m.components, item.File)

			concatenated := append(layoutStack, item.File)
			layoutErrors := errorStack

			rePattern := getPattern(segments, true)

			i := len(concatenated)

			log.Printf("I: %d", i)
			log.Printf("concatenated: %s", concatenated)
			log.Printf("layoutErrors: %s", layoutErrors)

			for ; i > 0; i-- {
				log.Printf("I2: %d", i)

				layoutErr := getItem(layoutErrors, i)
				layoutPage := getItem(concatenated, i)

				log.Printf("I2: %s %s", layoutErr, layoutPage)

				if layoutErr == "" && layoutPage == "" {
					layoutErrors = layoutErrors[:i-1]
					concatenated = concatenated[:i-1]
				}
			}

			i = len(layoutErrors)
			for ; i >= 0; i-- {
				if layoutErrors[i] != "" {
					break
				}
			}

			layoutErrors = layoutErrors[:i+1]

			path := ""
			if every(segments, func(segment []*RouteSegment) bool {
				return len(segment) == 1 && !segment[0].Dynamic
			}) {
				for _, segment := range segments {
					path += segment[0].Content
				}
			}

			m.routes = append(m.routes, RouteData{
				Type:     "page",
				Segments: simpleSegments,
				Pattern:  rePattern,
				Params:   params,
				Path:     path,
				A:        concatenated,
				B:        layoutErrors,
			})
		} else {
			pattern := getPattern(segments, item.RouteSuffix == "")

			m.routes = append(m.routes, RouteData{
				Type:     "endpoint",
				Segments: simpleSegments,
				Pattern:  pattern,
				File:     item.File,
				Params:   params,
			})
		}

	}

	return nil
}

func (m *Manifest) getParts(part, file string) []*RouteSegment {
	var result []*RouteSegment

	parts := rePart.FindStringSubmatch(part)

	if len(parts) == 0 {
		result = append(result, &RouteSegment{
			Content: part,
			Dynamic: false,
			Rest:    false,
		})

		return result
	}

	for i, str := range parts {
		dynamic := i%2 == 0
		content := str
		validContent := true

		if dynamic {
			tmpContent := reParseContent.FindStringSubmatch(str)
			if len(tmpContent) > 0 {
				content = tmpContent[0]
			}
			validContent = reValidContent.MatchString(content)
		}

		if content == "" || (dynamic && !validContent) {
			log.Fatalln(fmt.Errorf("Invalid route %s — parameter name must match /^[a-zA-Z0-9_$]+$/", file))
		}

		log.Println("Get partsconen", content, dynamic)

		result = append(result, &RouteSegment{
			Content: content,
			Dynamic: dynamic,
			Rest:    dynamic && reRestContent.MatchString(content),
		})
	}

	return result
}

func (m *Manifest) comparator(i, j int) bool {
	a := m.items[i]
	b := m.items[j]

	if a.IsIndex != b.IsIndex {
		if a.IsIndex {
			return isSpread(a.File)
		}
		return !isSpread(b.File)
	}

	max := math.Max(float64(len(a.Parts)), float64(len(b.Parts)))

	for i := 0; i < int(max); i++ {
		a_sub_part := a.Parts[i]
		b_sub_part := b.Parts[i]

		if a_sub_part == nil {
			return true
		}

		if b_sub_part == nil {
			return false
		}

		if a_sub_part.Rest && b_sub_part.Rest {
			if a.IsPage != b.IsPage {
				return a.IsPage
			}

			return !(a_sub_part.Content < b_sub_part.Content)
		}

		if a_sub_part.Rest != b_sub_part.Rest {
			return a_sub_part.Rest
		}

		if a_sub_part.Dynamic != b_sub_part.Dynamic {
			return a_sub_part.Dynamic
		}

		if !a_sub_part.Dynamic && a_sub_part.Content != b_sub_part.Content {
			contentLen := len(b_sub_part.Content) - len(a_sub_part.Content)
			if contentLen == 0 {
				return !(a_sub_part.Content < b_sub_part.Content)
			}
			return contentLen > 0
		}
	}

	if a.IsPage != b.IsPage {
		return a.IsPage
	}

	return !(a.File < b.File)
}

func isSpread(path string) bool {
	var rePath = regexp.MustCompile(`\[\.{3}`)
	return rePath.MatchString(path)
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

func some(arr []*RouteSegment, fn func(*RouteSegment) bool) bool {
	for _, item := range arr {
		if fn(item) {
			return true
		}
	}

	return false
}

func every(arr [][]*RouteSegment, fn func([]*RouteSegment) bool) bool {
	result := true
	for _, item := range arr {
		result = fn(item)
		if !result {
			break
		}
	}

	return result
}

func generateContent(segment []*RouteSegment) string {
	var content []string
	for _, part := range segment {
		if part.Dynamic {
			content = append(content, fmt.Sprintf("[%s]", part.Content))
		} else {
			content = append(content, part.Content)
		}
	}

	return strings.Join(content, "")
}

func getPattern(segments [][]*RouteSegment, addTrailingSlash bool) *regexp.Regexp {
	var pattern []string

	for _, segment := range segments {
		rest := segment[0].Rest
		if rest {
			pattern = append(pattern, "(?:\\/(.*))?")
		} else {
			pattern = append(pattern, "\\/")

			for _, part := range segment {
				if part.Dynamic {
					pattern = append(pattern, "([^/]+?)")
				} else {
					content := normalize(part.Content)

					content = regexp.MustCompile(`%5[Bb]`).ReplaceAllString(content, "[")
					content = regexp.MustCompile(`%5[Bd]`).ReplaceAllString(content, "]")
					content = regexp.MustCompile(`#`).ReplaceAllString(content, "%23")
					content = regexp.MustCompile(`\?`).ReplaceAllString(content, "%3F")
					content = regexp.MustCompile(`[.*+?^${}()|[\]\\]`).ReplaceAllString(content, "\\$&")
					pattern = append(pattern, content)
				}
			}
		}
	}

	log.Printf("pattern: %s", strings.Join(pattern, ""))

	return regexp.MustCompile(strings.Join(pattern, ""))
}

func normalize(str string) string {
	trans := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(trans, str)
	return result
}

func getItem(arr []string, i int) string {
	if len(arr) < i || i > len(arr) {
		return ""
	}

	return arr[i]
}
