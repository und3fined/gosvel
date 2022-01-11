/**
 * File: config.go
 * Project: config
 * File Created: 11 Jan 2022 20:58:42
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 11 Jan 2022 21:47:09
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package config

type KitFiles struct {
	Assets string `json:"assets"`
	Routes string `json:"routes"`
}

type KitConfig struct {
	Files KitFiles `json:"files"`
}

type Config struct {
	Extensions []string  `json:"extensions"`
	Kit        KitConfig `json:"kit"`
}
