import golitSSRPlugin from './plugins/golit-ssr.mjs';

export default function(eleventyConfig) {
  eleventyConfig.addPassthroughCopy({ "static": "/" });

  eleventyConfig.addPairedShortcode("raw_html", function(content) {
    return content;
  });

  eleventyConfig.addPlugin(golitSSRPlugin, {
    binary: '../../dist/golit',
    defs: 'bundles/',
    port: 9777,
    ignore: ['rh-audio-player', 'rh-footer'],
    preload: ['prism-esm'],
    concurrency: 2,
  });

  return {
    templateFormats: ["md", "njk", "html"],
    markdownTemplateEngine: "njk",
    htmlTemplateEngine: "njk",
    dir: {
      input: "src",
      output: "_site",
      includes: "_includes",
    },
  };
}
