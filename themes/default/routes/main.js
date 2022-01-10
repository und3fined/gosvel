/**
 * File: main.js
 * Project: gosvel-default
 * File Created: 10 Jan 2022 19:10:15
 * Author: und3fined (me@und3fined.com)
 * -----
 * Last Modified: 10 Jan 2022 19:10:29
 * Modified By: und3fined (me@und3fined.com)
 * -----
 * Copyright (c) 2022 und3fined.com
 */
import App from "./index.svelte";

const app = new App({
  target: document.body,
  props: {
    name: "world",
  },
});

export default app;
