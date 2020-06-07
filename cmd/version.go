package main

import "runtime"

var (
	// Version is the commit short version of this project
	Version string
	// Branch of the final compiled binary
	Branch string
	// GoVersion is the go compiler version that compiles this project
	GoVersion = runtime.Version()
)
