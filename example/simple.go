package main

import (
	"fmt"
	"github.com/xackery/discord"
)

func main() {
	c, e := discord.Create("email", "password")

	if e != nil {
		panic(e.Error())
	}

	fmt.Println(c.Session.Token)
}
