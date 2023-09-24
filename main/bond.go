/*
 * Copyright 2023-present Kuei-chun Chen. All rights reserved.
 * main/bond.go
 */
package main

import (
	"fmt"
	"time"

	"github.com/simagix/bond"
)

var repo = "simagix/bond"
var version = "devel-xxxxxx"

func main() {
	if version == "devel-xxxxxx" {
		version = "devel-" + time.Now().Format("20060102")
	}
	fullVersion := fmt.Sprintf(`%v %v`, repo, version)
	bond.Run(fullVersion)
}
