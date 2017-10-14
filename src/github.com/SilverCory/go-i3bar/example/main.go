package main

import (
	"github.com/SilverCory/go-i3bar"
	"os"
	"syscall"
	"time"
)

type MarqueeHelloWorldHandler struct {
	itr      int
	text     string
	reverse  bool
	position int
	go_i3bar.Handler
}

func (m *MarqueeHelloWorldHandler) GetMessage() *go_i3bar.Message {
	m.itr += 1
	itrMaxed := m.itr

	if m.itr > len(m.text)+12 {
		m.itr = -1
	}

	if itrMaxed > len(m.text)-1 {
		itrMaxed = len(m.text) - 1
	}

	ret := &go_i3bar.Message{
		Position: m.position,
		MinWidth: m.text,
		Urgent:   true,
	}

	if m.reverse {
		ret.FullText = m.text[itrMaxed:]
	} else {
		ret.FullText = m.text[:itrMaxed]
	}

	return ret

}

func (m *MarqueeHelloWorldHandler) Click(click *go_i3bar.Click) {
	m.reverse = !m.reverse
}

func main() {
	duration, err := time.ParseDuration("60ms")
	if err != nil {
		panic(err)
	}

	bar := go_i3bar.New(syscall.Signal(10), syscall.Signal(12), true, duration, os.Stdout, os.Stdin)

	bar.RegisterHandler("meme", "1", &MarqueeHelloWorldHandler{
		itr:      0,
		position: 2,
		text:     "Hello World!",
	})

	bar.RegisterHandler("meme", "2", &MarqueeHelloWorldHandler{
		itr:      3,
		position: 1,
		text:     "Goodbye Moon!",
	})

	panic(bar.Start())

}
