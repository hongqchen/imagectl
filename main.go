package main

import (
	"github.com/hongqchen/imagectl/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
