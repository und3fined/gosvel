/**
 * File: options.go
 * Project: manifest
 * File Created: 11 Jan 2022 23:41:01
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 11 Jan 2022 23:59:12
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
package manifest

import "gosvel/tools/kit/config"

type options struct {
	Conf   config.Config
	Cwd    string
	Output string
}

type Option func(*options)

func Config(conf config.Config) Option {
	return func(o *options) {
		o.Conf = conf
	}
}

func Cwd(cwd string) Option {
	return func(o *options) {
		o.Cwd = cwd
	}
}

func Output(output string) Option {
	return func(o *options) {
		o.Output = output
	}
}
