package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(hello())
}

func hello() string {
	return strings.Title("hello")
}
