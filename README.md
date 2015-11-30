# conflag [![Build Status](https://travis-ci.org/monochromegane/conflag.svg?branch=master)](https://travis-ci.org/monochromegane/conflag)

A combination command-line flag and configuration file library for Go.

## Usage

```go
// define your flags.
var procs int
flag.IntVar(&procs, "procs", runtime.NumCPU(), "GOMAXPROCS")

// set flags from configuration before parse command-line flags.
if args, err := conflag.ArgsFrom("/path/to/config.toml"); err == nil {
	flag.CommandLine.Parse(args)
}

// parse command-line flags.
flag.Parse()
```

and you create `/path/to/config.toml`

```toml
procs = 2
```

and run your app without option, `procs` flag will be set in `2` that is defined at configration file.

### Priority of flag

A priority of flag is

`command-line flag` > `configration file` > `flag default value`

In the above case,

| run                         | procs                              |
| --------------------------- | ---------------------------------- |
| myapp -procs 3              | 3                                  |
| myapp (with config-file)    | 2                                  |
| myapp (without config-file) | runtime.NumCPU() (default of flag) |

### Position

You can specify `positions` arguments to `ArgsFrom` function.

```toml
[options]
flag = "value"

[other settings]
hoge = "fuga"
```

```go
// parse configration only under the options section.
conflag.ArgsFrom("/path/to/config.toml", "options")
```

### List

You can use list for multiple parameters.
The following toml makes `-flag value1 -flag value2` arguments.

```toml
flag = [ "value1", "value2" ]
```


### go-flags

If you use [go-flags](https://github.com/jessevdk/go-flags) package, you can specify options like the following.

```go
parser := go-flags.NewParser(&opts, go-flags.Default)

conflag.LongHyphen = true
conflag.BoolValue = false
if args, err := conflag.ArgsFrom("/path/to/config.toml"); err == nil {
        parser.ParseArgs(args)
}
```

## Features

- Combine command-line flag and configuration file.
- Specify configration section.
- Specify list parameters.
- Support TOML configuration file.
- Support JSON configuration file.
- Support YAML configuration file.

## Installation

```sh
$ go get github.com/monochromegane/conflag
```

## Contribution

1. Fork it
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

## License

[MIT](https://github.com/monochromegane/conflag/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)

