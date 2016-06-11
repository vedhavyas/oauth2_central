package main

import (
	"fmt"
	"runtime"
)

var versionName = "1.0"
var versionCode = 1

func printVersion() {
	fmt.Printf("oauth2_central %v(%v) compiled with %v\n", versionName, versionCode, runtime.Version())
}
