package main

import (
	"github.com/xackery/discord/example/message"
)

var email = "email"
var password = "password"

//Listens
func main() {
	message.SendMessageExample(email, password)    //An example of how to send a message to a channel
	message.ReplyOnMessageExample(email, password) //An example of how to reply to a message event
}
