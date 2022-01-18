/**
 * File: types.go
 * Project: manifest
 * File Created: 11 Jan 2022 21:26:12
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 18 Jan 2022 15:33:36
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package manifest

type RouteDataType string

var (
	RouteDataEndpoint RouteDataType = "endpoint"
	RouteDataPage     RouteDataType = "page"
)

type RouteSegment struct {
	Content string `json:"content"` // content: string;
	Dynamic bool   `json:"dynamic"` // dynamic: boolean;
	Rest    bool   `json:"rest"`    // rest: boolean;
}

type RouteData struct {
	Type     RouteDataType   `json:"type"`
	Segments []*RouteSegment `json:"segments"`
	Pattern  string          `json:"pattern"`
	Params   []string        `json:"params"`
	Path     string          `json:"path,omitempty"`
	File     string          `json:"file,omitempty"`
	C        []string        `json:"c,omitempty"`
	E        []string        `json:"e,omitempty"`
}

type Asset struct {
	File string `json:"file"`             // file: string;
	Size int64  `json:"size"`             //size: number;
	Type string `json:"string,omitempty"` //type: string | null;
}

type ManifestData struct {
	Assets     []Asset     `json:"assets"`     // assets: Asset[];
	Layout     string      `json:"layout"`     // layout: string;
	Error      string      `json:"error"`      // error: string;
	Components []string    `json:"components"` //components: string[];
	Routes     []RouteData `json:"routes"`     //routes RouteData[];
}

type WalkItem struct {
	Basename    string
	Ext         string
	Parts       []*RouteSegment
	File        string
	IsDir       bool
	IsIndex     bool
	IsPage      bool
	RouteSuffix string

	// basename: string;
	//   ext: string;
	//   parts: Part[];
	//   file: string;
	//   is_dir: boolean;
	//   is_index: boolean;
	//   is_page: boolean;
	//   route_suffix: string;
}
