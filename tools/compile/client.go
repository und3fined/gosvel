/**
 * File: client.go
 * Project: compile
 * File Created: 31 Dec 2021 14:49:36
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 31 Dec 2021 15:40:26
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2021 und3fined.com
 */
package compile

import v8 "rogchap.com/v8go"

var CompileCtx *v8.Context

func New() *v8.Context {
	CompileCtx = v8.NewContext()

	return CompileCtx
}
