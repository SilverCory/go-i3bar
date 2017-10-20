package go_i3bar

import (
	"fmt"
	"os"
)

type Handler interface {
	GetMessage() *Message
	Click(click *Click)
}

func (b *Bar) RegisterHandler(name, instance string, handler Handler) {

	v, ok := b.handlers[name]
	if !ok {
		b.handlers[name] = make(map[string]Handler)
		v = b.handlers[name]
	}

	v[instance] = handler

}

func (b *Bar) FindHandler(click *Click) Handler {

	if val, ok := b.handlers[click.Name]; ok {
		if handler, ok := val[click.Instance]; ok {
			return handler
		}
	}

	return nil

}
