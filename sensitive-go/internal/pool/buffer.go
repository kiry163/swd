package pool

import (
	"bytes"
	"sync"
)

// BufferPool provides a pool of reusable byte buffers
var BufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GetBuffer retrieves a buffer from the pool
func GetBuffer() *bytes.Buffer {
	buf := BufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	BufferPool.Put(buf)
}

// StringBuilderPool provides a pool of reusable string builders
type StringBuilderPool struct {
	pool sync.Pool
}

// NewStringBuilderPool creates a new string builder pool
func NewStringBuilderPool(initSize int) *StringBuilderPool {
	return &StringBuilderPool{
		pool: sync.Pool{
			New: func() interface{} {
				b := new(bytes.Buffer)
				b.Grow(initSize)
				return b
			},
		},
	}
}

// Get retrieves a buffer from the pool
func (p *StringBuilderPool) Get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put returns a buffer to the pool
func (p *StringBuilderPool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	p.pool.Put(buf)
}


