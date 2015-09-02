package main

import (
	"fmt"
	"github.com/xackery/discord"
)

func main() {
	c, e := discord.Create("test@here.com", "password")
	if e != nil {
		panic(e.Error())
	}

	fmt.Println("Simple example", c)

}
