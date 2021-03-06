(() => {
  // themes/default/routes/index.svelte
  var routes_default = '/* themes/default/routes/index.svelte generated by Svelte v3.45.0 */\nimport { create_ssr_component, escape } from "svelte/internal";\n\nconst css = {\n	code: "main.svelte-1tky8bj{text-align:center;padding:1em;max-width:240px;margin:0 auto}h1.svelte-1tky8bj{color:#ff3e00;text-transform:uppercase;font-size:4em;font-weight:100}@media(min-width: 640px){main.svelte-1tky8bj{max-width:none}}",\n	map: "{\\"version\\":3,\\"file\\":\\"index.svelte\\",\\"sources\\":[\\"index.svelte\\"],\\"sourcesContent\\":[\\"<script>\\\\n\\\\texport let name;\\\\n<\/script>\\\\n\\\\n<main>\\\\n\\\\t<h1>Hello {name}!</h1>\\\\n\\\\t<p>Visit the <a href=\\\\\\"https://svelte.dev/tutorial\\\\\\">Svelte tutorial</a> to learn how to build Svelte apps.</p>\\\\n</main>\\\\n\\\\n<style>\\\\n\\\\tmain {\\\\n\\\\t\\\\ttext-align: center;\\\\n\\\\t\\\\tpadding: 1em;\\\\n\\\\t\\\\tmax-width: 240px;\\\\n\\\\t\\\\tmargin: 0 auto;\\\\n\\\\t}\\\\n\\\\n\\\\th1 {\\\\n\\\\t\\\\tcolor: #ff3e00;\\\\n\\\\t\\\\ttext-transform: uppercase;\\\\n\\\\t\\\\tfont-size: 4em;\\\\n\\\\t\\\\tfont-weight: 100;\\\\n\\\\t}\\\\n\\\\n\\\\t@media (min-width: 640px) {\\\\n\\\\t\\\\tmain {\\\\n\\\\t\\\\t\\\\tmax-width: none;\\\\n\\\\t\\\\t}\\\\n\\\\t}\\\\n</style>\\"],\\"names\\":[],\\"mappings\\":\\"AAUC,IAAI,eAAC,CAAC,AACL,UAAU,CAAE,MAAM,CAClB,OAAO,CAAE,GAAG,CACZ,SAAS,CAAE,KAAK,CAChB,MAAM,CAAE,CAAC,CAAC,IAAI,AACf,CAAC,AAED,EAAE,eAAC,CAAC,AACH,KAAK,CAAE,OAAO,CACd,cAAc,CAAE,SAAS,CACzB,SAAS,CAAE,GAAG,CACd,WAAW,CAAE,GAAG,AACjB,CAAC,AAED,MAAM,AAAC,YAAY,KAAK,CAAC,AAAC,CAAC,AAC1B,IAAI,eAAC,CAAC,AACL,SAAS,CAAE,IAAI,AAChB,CAAC,AACF,CAAC\\"}"\n};\n\nconst Routes = create_ssr_component(($$result, $$props, $$bindings, slots) => {\n	let { name } = $$props;\n	if ($$props.name === void 0 && $$bindings.name && name !== void 0) $$bindings.name(name);\n	$$result.css.add(css);\n\n	return `<main class="${"svelte-1tky8bj"}"><h1 class="${"svelte-1tky8bj"}">Hello ${escape(name)}!</h1>\n	<p>Visit the <a href="${"https://svelte.dev/tutorial"}">Svelte tutorial</a> to learn how to build Svelte apps.</p>\n</main>`;\n});\n\nexport default Routes;';

  // themes/default/routes/main.js
  var app = new routes_default({
    target: document.body,
    props: {
      name: "world"
    }
  });
  var main_default = app;
})();
