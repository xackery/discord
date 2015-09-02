package storage

import (
	"fmt"
	"github.com/xackery/discord"
	"testing"
)

func TestClientConnect(t *testing.T) {
	c, e := discord.Create("test@here.com", "password")
	if e != nil {
		t.Error(fmt.Println("Error:", e.Error()))
	}
	fmt.Println("C: ", c)
	return
}
