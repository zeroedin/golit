---
type: patch
---

Remove TypeScript support from golit bundle. golit is a post-build tool that accepts compiled JavaScript from node_modules/, not TypeScript source. The directory walker no longer collects .ts/.tsx files, and the inline BundleSource loader uses LoaderJS instead of LoaderTS.
