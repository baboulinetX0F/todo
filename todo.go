package main

import (
	"fmt"
)

const (
	RedColor = "\033[1;31m%s\033[0m"
)

type Task struct {
	id          uint16
	description string
	status      bool
}

func main() {
	fmt.Printf(RedColor, "☑ Test")
}
