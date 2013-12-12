# Configuration File Parser for Go

## Syntax

Parser supports a subset of version
[TOML v0.2.0](https://github.com/mojombo/toml/blob/master/versions/toml-v0.2.0.md).

Main differences are:

* section and option names are case insensitive;
* sub-sections are not supported;
* multi-dimensional arrays are not supported;
* tables are not supported; and
* array of tables are not supported.

Why? I just don't need multi-dimensional arrays or tables in my configuration files.

Parser implements following grammar ([EBNF](http://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_Form) style):

```
	config = section | options
	section = '[' IDENTIFIER ']' EOL options
	options = option {option}
	option =  IDENTIFIER '=' (value | array) EOL
	value = BOOL | INT | FLOAT | DATE | STRING
	array = '[' {EOF} value {EOF} {, {EOF} value} {EOF} ']'
```

Parser can load as many configurations as you want, and each load can update existing options.

## Installation

    go get github.com/cbonello/gp-config

## Testing

`gp-config` uses [gocheck](http://labix.org/gocheck) for testing.

To install gocheck, just run the following command:

	go get launchpad.net/gocheck

To run the test suite, use the following command:

	go test github.com/cbonello/gp-config -v

You can also check test coverage if you're using go 1.2 or later:

	go test github.com/cbonello/gp-config -v -coverprofile=test.out
	go tool cover -html=test.out

## Usage

### Loading Configuration Files

Configuration can either be stored in a string or a file.

```go
import (
	"flag"
	"fmt"
	"github.com/cbonello/gp-config"
	"os"
)
const deflt = `version = [1, 0, 10]
			   [database]
				   dbname = "mydb"
				   user = "foo"
				   password = "bar"`
var dev bool = false
func main() {
	flag.BoolVar(&dev, "dev", false, "Runs application in debug mode, default is production.")
	flag.Parse()
	// Set default options (Production mode).
	cfg := config.NewConfiguration()
	if err := cfg.LoadString(deflt); err != nil {
		fmt.Printf("error: default config: %d:%d: %s\n",
			err.Line, err.Column, err)
		os.Exit(1)
	}
	if dev {
		// Override default options with debug mode settings.
		if err := cfg.LoadFile("debug.cfg"); err != nil {
			fmt.Printf("error: %s: %d:%d: %s\n",
				err.Filename, err.Line, err.Column, err)
			os.Exit(1)
		}
	}
	...
}
```

Contents of `debug.cfg` may for instance be:

```toml
[database]
	dbname = "mydb_test"
```

### Reading Configuration Files

#### Basic API

```go
	if dbname, err := cfg.GetString("database.dbname"); err != nil {
		fmt.Printf("error: configuration: missing database name\n")
		os.Exit(1)
	}
	user := cfg.GetStringDefault("database.dbname", "user")
	password := cfg.GetStringDefault("database.dbname", "foobar")
```

See `examples/demo/` for a complete demo application.

Full API documentation is available at [godoc.org](http://godoc.org/github.com/cbonello/gp-config).

#### Reflection API

Following structure may be declared to record options of `[database]` section defined above.

```go
	type database struct {
		Name     string `option:"dbname"`
		Username string `option:"user"`
		Password string
	}
```

A structure field annotation may be used if there is no direct match between a field name and an option name.

And finally, to decode:

```go
	var db database
	if err := cfg.Decode("server", &db); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
```

See `examples/demo-decode/` for a complete demo application.

Full API documentation is available at [godoc.org](http://godoc.org/github.com/cbonello/gp-config).

## Examples

Demo applications are provided in the `examples/` directory. To launch them:

    go run github.com/cbonello/gp-config/examples/demo/demo.go
    go run github.com/cbonello/gp-config/examples/demo-decode/demo.go
