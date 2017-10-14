package go_i3bar

type Handler interface {
	GetMessage() *Message
	Click()
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

	if val, ok := b.handlers[click.name]; ok {
		if handler, ok := val[click.instance]; ok {
			return handler
		}
	}

	return nil

}
