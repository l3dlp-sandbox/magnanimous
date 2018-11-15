{{ define path "/sections/docs/get_started.html" }}
{{ define chapter 1 }}
{{ define title "Getting Started" }}
{{ include _docs_header.html }}

# 1. Getting Started

This section briefly describes how you can get Magnanimous and quickly get a simple static website built using it.

It's intentionally light on details! Check the next chapters to learn more. 

## Downloading

Magnanimous is a single executable file. To download it using the command line, type the following:

```
$ curl <url-will-go-here>
```

You can also download it manually from [this page](link-will-go-here).

## Creating the source directories

Magnanimous expects you to have a simple directory structure like this:

<img src="{{ eval baseURL + "/images/docs/my-website.jpg"  }}" width="500em;" alt="Magnanimous directories" />

The `README.md` file is not really needed, but you can document how your website works and things like that in this file.

All **source** files go, unsurprisingly, in the `source/` directory, under one of its sub-directories.

> A **source** file is just a file that you want to be in your website. We call it source to differentiate it from
  the generated files, which we are just about to meet!

The `source/` directory contains two sub-directories:

* `static/` can be used for files that should be just copied to your website without modifications.
* `processed/` is where **processed** files are placed.

Only files within these two directories will be present in the final website.

A **processed** file is one which contains [Magnanimous instructions]({{ eval INSTRUCTIONS_PATH }}), which are 
used to modify the actual contents of the file or provide metadata.

> **Metadata** is some information about a file, like its path, title, or the date it was created, which can be used 
  elsewhere (e.g. on the table of contents) or included in the visible file contents.
  
Any file whose name starts with an underscore, like `_header.html`, will **not** be present in the final website.
But they are useful to create [Components](components.html), or _fragments_ which can be included into other files.

<img src="{{ eval baseURL + "/images/docs/transformation.jpg" }}" width="500em;" alt="Magnanimous transformation" />

## Building the website

After creating a few source files, from inside the root directory (`my-website/`, in the above example), just run the
Magnanimous executable from the command line:

```
$ magnanimous
```

You can also run `magnanimous` from other directories, in which case you just need to give the website's directory
as an argument:

```
$ magnanimous path/to/my-website
```

This will create a static website in the `target/` directory.

## Testing the website

Now that your website is ready, you can run any web server to serve the `target/` directory so you can see what
the website looks like.

For example, if you have Python installed, simply run this command:

```
$ cd target/ && python -m SimpleHTTPServer 8082
```

> [Click here](https://gist.github.com/willurd/5720255) for many other web server one-liners.

Then, open `http://localhost:8082` on your browser. If you have a file called `index.html`, that page should be shown
in the browser... otherwise, just add the path to one of your files to the URL (e.g. if you have a source file under 
`source/processed/blog/blog1.html`, try opening `http://localhost:8082/blog/blog1.html`).

To publish the website publicly so everyone can visit it, you can use one of many available static website hosts:

* [Netlify](https://www.netlify.com/)
* [GitHub Pages](https://pages.github.com/)
* [Digital Ocean](https://www.digitalocean.com/)
* [BitBalloon](https://www.bitballoon.com/)
* [AWS (Amazon)](http://docs.aws.amazon.com/gettingstarted/latest/swh/website-hosting-intro.html)

The [Basic Tutorial](basic_tutorial.html) shows how to publish a website with GitHub Pages. 

{{ include _docs_footer.html }}