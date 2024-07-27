package window_buffer

import "sync"

type value struct {
	v string
}

type Buffer struct {
	data []value
	mx   sync.Mutex
}

func New(size uint) *Buffer {
	return &Buffer{
		data: make([]value, size),
	}
}

func (b *Buffer) Add(v string) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.data = append(b.data, value{v: v})
	b.data = b.data[1:]
}

func (b *Buffer) Get() []value {
	b.mx.Lock()
	defer b.mx.Unlock()

	// copy slice
	return append([]value{}, b.data...)
}
