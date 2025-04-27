package main

import (
	fetch "fetchlinks/internal"
	"fmt"
)

func main() {
	typeValue := "top"
	resp, err := fetch.GetHackerNewsLinks(fetch.StoryType(typeValue))
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
