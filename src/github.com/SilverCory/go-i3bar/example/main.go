package main

import (
	"github.com/SilverCory/go-i3bar"
	"os"
	"syscall"
	"time"
)

var itr = 0
var text = "Hello World!"

type MarqueeHelloWorldHandler struct {
	go_i3bar.Handler
}

func (m *MarqueeHelloWorldHandler) GetMessage() *go_i3bar.Message {
	if itr > len(text)-1 {
		itr = 0
	} else {
		itr += 1
	}
	return &go_i3bar.Message{
		Position: 1,
		FullText: text[:itr],
		MinWidth: text,
		Urgent:   true,
	}
}

func (m *MarqueeHelloWorldHandler) Click() {}

func main() {
	duration, err := time.ParseDuration("60ms")
	if err != nil {
		panic(err)
	}

	bar := go_i3bar.New(syscall.Signal(10), syscall.Signal(12), true, duration, os.Stdout, os.Stdin)

	meme := &MarqueeHelloWorldHandler{}

	bar.RegisterHandler("meme", "", meme)
	bar.RegisterHandler("meme", "2", meme)

	bar.Start()

}
