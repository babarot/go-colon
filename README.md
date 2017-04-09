go-colon [![][travis-badge]][travis-link] [![][license-badge]][license-link] [![][godoc-badge]][godoc-link]
========

Package colon for parsing colon-separated strings like PATH

## Installation

```console
$ go get github.com/b4b4r07/go-colon
```

## Usage

## Example

Set the variables as you like [^1]:

```bash
export ENHANCD_FILTER="fzy:/usr/local/bin/peco --select-1:fzf --multi --ansi:zaw"
```

Run this script like [`_example`](_example/main.go).

```go
result, err := colon.Parse(os.Getenv("ENHANCD_FILTER"))
if err != nil {
    panic(err)
}
pp.Println(result.Executable())
```

Result:

```go
&colon.Result{
  colon.Object{
    Index: 1,
    Attr:  colon.Attribute{
      First: "fzy",
      Other: []string{},
      Args:  []string{
        "fzy",
      },
      Base:    "fzy",
      Dir:     "",
      IsDir:   false,
      Command: "/Users/b4b4r07/.zplug/bin/fzy",
    },
    Errors: []error{},
  },
  colon.Object{
    Index: 2,
    Attr:  colon.Attribute{
      First: "/usr/local/bin/peco",
      Other: []string{
        "--select-1",
      },
      Args: []string{
        "/usr/local/bin/peco",
        "--select-1",
      },
      Base:    "peco",
      Dir:     "/usr/local/bin",
      IsDir:   false,
      Command: "/usr/local/bin/peco",
    },
    Errors: []error{},
  },
  colon.Object{
    Index: 3,
    Attr:  colon.Attribute{
      First: "fzf",
      Other: []string{
        "--multi",
        "--ansi",
      },
      Args: []string{
        "fzf",
        "--multi",
        "--ansi",
      },
      Base:    "fzf",
      Dir:     "",
      IsDir:   false,
      Command: "/Users/b4b4r07/.zplug/bin/fzf",
    },
    Errors: []error{},
  },
}
```

For more info, see [godoc][godoc-link].

## License

MIT

## Author

b4b4r07

[^1]: https://github.com/b4b4r07/enhancd

[travis-badge]: http://img.shields.io/travis/b4b4r07/go-colon.svg?style=flat-square
[travis-link]: https://travis-ci.org/b4b4r07/go-colon

[license-badge]: http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square
[license-link]: https://github.com/b4b4r07/go-colon/blob/master/LICENSE

[godoc-badge]: http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square
[godoc-link]: http://godoc.org/github.com/b4b4r07/go-colon
