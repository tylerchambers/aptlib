package main

import (
	"github.com/tylerchambers/aptlib"
)

func main() {
	c := new(aptlib.Client)
	c.Init()
	c.BuildRepoURIs()
	for i, v := range c.RepoURIs {
		print(i, v)
	}
}
