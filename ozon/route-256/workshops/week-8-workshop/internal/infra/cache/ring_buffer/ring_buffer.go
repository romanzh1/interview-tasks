package ring_buffer

import (
	"sync"
)

type value struct {
	v string
}

type Buffer struct {
	size  int64
	index int64
	data  []value // для вытеснения из кэша
	mx    sync.Mutex
}

func New(size uint) *Buffer {
	data := make([]value, size)
	for i := range data {
		data[i] = value{}
	}

	return &Buffer{
		size:  int64(size),
		index: int64(size),
		data:  data,
	}
}

func (b *Buffer) Add(v string) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if b.size == 0 {
		// handle divide by zero error
		return
	}

	b.data[b.index%b.size].v = v
	b.index++
}

func (b *Buffer) Get() []value {
	b.mx.Lock()
	defer b.mx.Unlock()

	// copy slice
	return append([]value{}, b.data...)
}
