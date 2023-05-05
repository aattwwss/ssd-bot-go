package search

import (
	"fmt"
	"regexp"
)

func main() {
	str := "Golang@%Programs#"
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "")
	fmt.Println(str)
	fmt.Println("vim-go")
}
