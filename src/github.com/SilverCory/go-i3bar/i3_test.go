package go_i3bar_test

import (
	"github.com/SilverCory/go-i3bar"
	"os"
	"syscall"
	"testing"
	"time"
)

type HelloWorldHandler struct {
	go_i3bar.Handler
}

func (m *HelloWorldHandler) GetMessage() *go_i3bar.Message {
	return &go_i3bar.Message{
		Position: 1,
		FullText: "Hello World!",
	}
}

func (m *HelloWorldHandler) Click(click *go_i3bar.Click) {}

func TestStart(t *testing.T) {
	duration, err := time.ParseDuration("500ms")
	if err != nil {
		panic(err)
	}

	bar := go_i3bar.New(syscall.Signal(10), syscall.Signal(12), true, duration, os.Stdout, os.Stdin)

	meme := &HelloWorldHandler{}

	bar.RegisterHandler("meme", "", meme)
	bar.RegisterHandler("meme", "2", meme)

	go bar.Start()

	exitDuration, _ := time.ParseDuration("4s")
	select {
	case <-time.After(exitDuration):
		bar.Close()
	}

}
