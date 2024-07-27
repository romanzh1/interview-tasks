package window_buffer

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func vSlice(vv ...int) []value {
	r := make([]value, len(vv))
	for i, v := range vv {
		if v == 0 {
			continue
		} else {
			r[i] = value{v: strconv.FormatInt(int64(v), 10)}
		}
	}

	return r
}

func TestBuffer_Add(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		b := New(3)

		b.Add("1")
		assert.Equal(t, vSlice(0, 0, 1), b.Get())

		b.Add("2")
		assert.Equal(t, vSlice(0, 1, 2), b.Get())

		b.Add("3")
		assert.Equal(t, vSlice(1, 2, 3), b.Get())

		b.Add("4")
		assert.Equal(t, vSlice(2, 3, 4), b.Get())
	})

	t.Run("zero size window", func(t *testing.T) {
		b := New(0)

		b.Add("1")
		assert.Equal(t, []value{}, b.Get())

		b.Add("2")
		assert.Equal(t, []value{}, b.Get())

		b.Add("3")
		assert.Equal(t, []value{}, b.Get())

		b.Add("4")
		assert.Equal(t, []value{}, b.Get())
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

// @see also https://stackoverflow.com/questions/55045402/memory-leak-in-golang-slice
func Benchmark_Add(b *testing.B) {
	buffer := New(1500)

	str := strings.Repeat("a", 1000)

	for i := 0; i < 1000000; i++ {
		buffer.Add(str)
	}
}
