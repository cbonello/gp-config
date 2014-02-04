package main

import (
	"fmt"
	"github.com/cbonello/gp-config"
	"math/rand"
	"os"
	"time"
)

const (
	deflt = `
version = [1, 0, 0]

[server]
	URL = "www.myurl.com"

[database]
	dbname = "mydb"
	user = "foo"
	password = "bar"
`
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func main() {
	// Set default options (Production mode).
	cfg := config.NewConfiguration()
	if err := cfg.LoadString(deflt); err != nil {
		fmt.Printf("error: default config: %d:%d: %s\n",
			err.Line, err.Column, err)
		os.Exit(1)
	}

	// Simulate debug or production run.
	debug := (random(1, 10) < 6)
	if debug {
		fmt.Println("DEBUG MODE")
		// Override default options with debug mode settings.
		if err := cfg.LoadFile("debug.cfg"); err != nil {
			fmt.Printf("error: %s: %d:%d: %s\n",
				err.Filename, err.Line, err.Column, err)
			os.Exit(1)
		}
	} else {
		fmt.Println("PRODUCTION MODE")
	}

	version := cfg.GetIntArrayDefault("version", []int64{0, 0, 1})

	server := cfg.GetStringDefault("server.URL", "www.myurl.com")
	port := cfg.GetIntDefault("server.port", 80)

	dbname := cfg.GetStringDefault("database.dbname", "mydb")
	user, err := cfg.GetString("database.user")
	if err != nil {
		fmt.Printf("error: %s: database: missing user name\n", dbname)
		os.Exit(1)
	}
	password, err := cfg.GetString("database.password")
	if err != nil {
		fmt.Printf("error: %s: missing password\n", dbname)
		os.Exit(1)
	}

	fmt.Printf("version  = %d.%d.%d\n", version[0], version[1], version[2])
	fmt.Printf("server   = '%s'\n", server)
	fmt.Printf("port     = %d\n", port)
	fmt.Printf("dbname   = '%s'\n", dbname)
	fmt.Printf("user     = '%s'\n", user)
	fmt.Printf("password = '%s'\n", password)
}
