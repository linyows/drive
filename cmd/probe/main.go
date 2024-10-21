package main

import "os"

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	c := newCmd(os.Args)
	if c != nil {
		c.start()
	}
}
