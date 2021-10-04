package main

import (
	"fmt"
	"log"

	"github.com/tylerchambers/goapt/pkg/parse"
)

func main() {
	src, err := parse.AllSources()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range src {
		fmt.Println(*v)
	}
}
