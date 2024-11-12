# go-vanity-ssg

CLI tool allowing generation static-site pages generation enabling support for
golang's [vanity imports](https://pkg.go.dev/cmd/go#hdr-Remote_import_paths).

Main features of this implementation:

- **Emit index page with list of all packages**  
  Allows to showcase all packages that are supported by vanity import source.
- **Static pages**  
  No need in any sort of webserver - you can provide vanity imports even being hosted from GitHub pages.
- **Public/private packages split.**  
  In cases when you need to enable clients from local networks to be redirected to local repositories and documentation,
  this tool allows emitting static assets for such cases. With a little bit of webserver tweaking you will be able to
  serve different content depending on clint source.
- **Customizable templates**  
  Tool allows to load own templates, for the case you want to swag a bit.

![packages index example](.github/index-example.png)

## Installation

```shell
```

## Usage

### `go-vanity-ssg emit-config`

Emits example config file.

**_Flags:_**

- **`--config`** - path to config file. Default: `./.vanity.config.yaml`
- **`--overwrite`** - overwrite config file if it already exists.

### `go-vanity-ssg emit-templates`

Emits template files embedded into tool.

**_Flags:_**

- **`--dir`** - path to directory to emit templates to. Default: `./templates`
- **`--overwrite`** - overwrite template files if exists.

### `go-vanity-ssg`

Emit generated html files.

**_Flags:_**

- **`--config`** - path to config file to use. Default: `./.vanity.config.yaml`
- **`--out-dir`** - Directory to emit html files. Default: `./dist`
- **`--templates-dir`** - Directory to load custom templates from.
- **`--public`** - Emit public packages files (default).
- **`--private`** - Emit private packages files.
- **`--no-inherit-public`** - Do not include public packages to the list of private packages.