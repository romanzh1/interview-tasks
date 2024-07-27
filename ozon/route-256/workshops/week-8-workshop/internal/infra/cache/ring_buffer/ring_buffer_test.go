package ring_buffer

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func vSlice(vv ...int) []*value {
	r := make([]*value, len(vv))
	for i, v := range vv {
		if v == 0 {
			r[i] = &value{v: ""}
		} else {
			r[i] = &value{v: strconv.FormatInt(int64(v), 10)}
		}
	}

	return r
}

func TestBuffer_Add(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		b := New(3)

		b.Add("1")
		assert.Equal(t, vSlice(1, 0, 0), b.Get())

		b.Add("2")
		assert.Equal(t, vSlice(1, 2, 0), b.Get())

		b.Add("3")
		assert.Equal(t, vSlice(1, 2, 3), b.Get())

		b.Add("4")
		assert.Equal(t, vSlice(4, 2, 3), b.Get())
	})

	t.Run("zero size window", func(t *testing.T) {
		b := New(0)

		b.Add("1")
		assert.Equal(t, []*value{}, b.Get())

		b.Add("2")
		assert.Equal(t, []*value{}, b.Get())

		b.Add("3")
		assert.Equal(t, []*value{}, b.Get())

		b.Add("4")
		assert.Equal(t, []*value{}, b.Get())
	})

	t.Run("one size window", func(t *testing.T) {
		b := New(1)

		b.Add("1")
		assert.Equal(t, vSlice(1), b.Get())

		b.Add("2")
		assert.Equal(t, vSlice(2), b.Get())

		b.Add("3")
		assert.Equal(t, vSlice(3), b.Get())

		b.Add("4")
		assert.Equal(t, vSlice(4), b.Get())
	})
}

func Benchmark_Add(b *testing.B) {
	buffer := New(1500)

	str := strings.Repeat("a", 1000)
	b.ReportAllocs()
	for i := 0; i < 1000000; i++ {
		buffer.Add(str)
	}
}
