package main

import (
	"fmt"
	"os"
	"time"
)

var version = "dev"

func main() {
	fmt.Println(version)

	if len(os.Args) > 1 && os.Args[1] == "wait" {
		time.Sleep(time.Minute)
	}
}
