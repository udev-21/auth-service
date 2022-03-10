package main

import (
	"fmt"
	"strings"
)

func main1() {
	fmt.Println(strings.TrimSpace("asdf"))
	fmt.Println(strings.TrimSpace("asdf "))
	fmt.Println(strings.TrimSpace(" asdf"))
	fmt.Println(strings.TrimSpace(" asdf "))

	fmt.Println(strings.TrimSpace("\tasdf "))
	fmt.Println(strings.TrimSpace("asdf\t"))
	fmt.Println(strings.TrimSpace("\tasdf\t"))

	fmt.Println(strings.TrimSpace("\nasdf"))
	fmt.Println(strings.TrimSpace("asdf\n"))
	fmt.Println(strings.TrimSpace("\nasdf\n"))

	fmt.Println(strings.TrimSpace("\t\nasdf\n"))
	fmt.Println(strings.TrimSpace("\nasdf\n"))
}
