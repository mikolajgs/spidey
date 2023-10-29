# spidey
Tiny static website generator

Spidey is a tiny tool that generates static website from HTML and markdown snippets.  The generated 
website can contain pages as well as posts.  HTML code is defined by layouts and snippets, such as 
header, footer or list of posts etc.  Pages and posts can be written in either pure HTML or markdown.
These guys are glued together to form a static websites.

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
