package go_i3bar

import (
	"bufio"
	"encoding/json"
	"io"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
	"os"
)

type Bar struct {
	protocol Protocol
	duration time.Duration
	writer   io.Writer
	reader   io.Reader
	encoder *json.Encoder
	decoder *json.Decoder
	handlers map[string]map[string]Handler

	exit chan int
	mux  sync.Mutex
}

type Protocol struct {
	Version     int  `json:"version"`
	StopSignal  int  `json:"stop_signal"`
	ContSignal  int  `json:"cont_signal"`
	ClickEvents bool `json:"click_events"`
}

type Click struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Button   int    `json:"button"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type Message struct {
	Position       int    `json:"-"`
	FullText       string `json:"full_text"`
	ShortText      string `json:"short_text,omitempty"`
	Color          string `json:"color,omitempty"`
	Background     string `json:"background,omitempty"`
	Border         string `json:"border,omitempty"`
	MinWidth       string `json:"min_width,omitempty"`
	Align          Align  `json:"align,omitempty"`
	Name           string `json:"name,omitempty"`
	Instance       string `json:"instance,omitempty"`
	Urgent         bool   `json:"urgent,omitempty"`
	Separator      bool   `json:"separator,omitempty"`
	SeparatorWidth int    `json:"separator_block_width,omitempty"`
}

type Align string

const (
	LEFT   Align = "left"
	RIGHT  Align = "right"
	CENTER Align = "center"
)

func New(stopSignal syscall.Signal, continueSignal syscall.Signal, clickEvents bool, duration time.Duration, writer io.Writer, reader io.Reader) *Bar {

	return &Bar{
		protocol: Protocol{
			Version:     1,
			StopSignal:  int(stopSignal),
			ContSignal:  int(continueSignal),
			ClickEvents: clickEvents,
		},
		duration: duration,
		writer:   writer,
		reader:   reader,
		handlers: make(map[string]map[string]Handler),
	}

}

func (b *Bar) gatherMessages() []*Message {

	ret := []*Message{}
	for name, v := range b.handlers {
		for instance, handler := range v {
			msg := handler.GetMessage()
			msg.Instance = instance
			msg.Name = name
			ret = append(ret, msg)
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Position < ret[j].Position
	})

	return ret

}

func (b *Bar) Start() error {

	b.exit = make(chan int, 1)
	b.encoder = json.NewEncoder(b.writer)

	if err := b.encoder.Encode(b.protocol); err != nil {
		return err
	}

	running := make(chan error, 1)
	go func() {
		for {
			time.Sleep(b.duration)
			if ret := b.Draw(); ret != nil {
				running <- ret
			}
		}
	}()

	go func() {
		running <- b.Read()
	}()

	for {
		select {
		case ret := <-running:
			return ret
		case <-b.exit:
			return nil
		}
	}

	// TODO implement click via std in?

	return nil

}
func (b *Bar) Read() error {

	b.decoder = json.NewDecoder(b.reader)
	for {
		click := &Click{}
		err := b.decoder.Decode(click)
		if err != nil {
			return err
		}

		handler := b.FindHandler(click)
		if handler != nil {
			handler.Click(click)
		}

	}

}

func (b *Bar) Draw() error {
	b.mux.Lock()

	if status, err := json.Marshal(b.gatherMessages()); err != nil {
		return err
	} else {
		if err := b.encoder.Encode(status); err != nil {
			return err
		}
	}

	b.mux.Unlock()

	return nil
}

func (b *Bar) Close() {
	b.exit <- 0
}
