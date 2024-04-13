# spidey
Tiny static website generator

Spidey is a tool designed to create static websites from HTML and Markdown snippets.  The tool
enables the construction of websites that include both pages and posts.  It utilizes predefined layouts
and snippets—like headers, footers, and post lists—to structure content.  You can write pages and posts
in either pure HTML or Markdown, which Spidey then seamlessly assembles into a complete static website.

Have a butcher's at `src` and `dist` in the `examples` directory to see it in action.

Spidey has been rapidly written to replace other tool.  It contains minimal functionality, just to fulfill
requirement of creating simple websites.  Hence, some bits in the code are hardcoded, and there might be
many TODOs.

### Building
Run `go build` in the root directory to build the binary.

### Running
Spidey has one command called `generate` which takes two arguments:
* source directory where the website configuration, layouts, pages, posts and other contents are located
* destination directory where HTML files should be generated, and this one has to be empty

#### Quick start
Create any empty directory where HTML files should be written, eg. `/tmp/spidey-generated-files` and run
the following command from root of this repository:

    spidey generate -s $(pwd)/examples/src -d /tmp/spidey-generated-files

#### Live preview
Spidey does not allow live preview of the website yet.  However, you can use nginx docker container that
would run in the background and serve the pages.

To set things up, open the `_config.yml` file in an editor and change the value of `baseurl` to 
`http://localhost:8080` (or any other port).
Next, start nginx container and mount destination directory to it (place where HTML files are generated):

    docker run --name some-nginx -p 8080:80 -v /tmp/spidey-generated-files:/usr/share/nginx/html:ro -d nginx
