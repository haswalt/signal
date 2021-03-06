package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-07-25 14:10:01.591695 +0200 CEST m=+0.013706761

import "math"

// Float64 is float64 floating-point signal.
type Float64 struct {
	buffer []float64
	channels
}

// Float64 allocates a new sequential float64 signal buffer.
func (a Allocator) Float64() Floating {
	return Float64{
		buffer:   make([]float64, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
	}
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s Float64) AppendSample(value float64) Floating {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, float64(value))
	return s
}

// SetSample sets sample value for provided index.
func (s Float64) SetSample(i int, value float64) {
	s.buffer[i] = float64(value)
}

// GetFloat64 selects a new sequential float64 signal buffer.
// from the pool.
func (p *Pool) GetFloat64() Floating {
	if p == nil {
		return nil
	}
	return p.f64.Get().(Floating)
}

// PutFloat64 places signal buffer back to the pool. If a type of
// provided buffer isn't Float64 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (p *Pool) PutFloat64(s Floating) {
	if p == nil {
		return
	}
	if sig, ok := s.(Float64); ok {
		capacity := cap(sig.buffer)
		for i, buf := 0, sig.buffer[:capacity]; i < capacity; i++ {
			buf[i] = 0
		}
	} else {
		panic("pool put float64 invalid type")
	}
	mustSameCapacity(s.Capacity(), p.allocator.Capacity)
	if s.Length() != p.allocator.Length {
		s = s.Slice(0, p.allocator.Length)
	}
	p.f64.Put(s)
}

// Capacity returns capacity of a single channel.
func (s Float64) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Float64) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s Float64) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Float64) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided channel and index.
func (s Float64) Sample(i int) float64 {
	return float64(s.buffer[i])
}

// Append appends [0:Length] samples from src to current buffer and returns new
// Floating buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s Float64) Append(src Floating) Floating {
	mustSameChannels(s.Channels(), src.Channels())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]float64, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Floating(s)
	for i := 0; i < src.Len(); i++ {
		result = result.AppendSample(src.Sample(i))
	}
	return result
}

// Slice slices buffer with respect to channels.
func (s Float64) Slice(start, end int) Floating {
	start = BufferIndex(s.Channels(), 0, start)
	end = BufferIndex(s.Channels(), 0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// ReadFloat64 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadFloat64(src Floating, dst []float64) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = float64(src.Sample(i))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedFloat64 reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedFloat64(src Floating, dst [][]float64) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = float64(src.Sample(BufferIndex(src.Channels(), c, i)))
		}
	}
	return
}

// WriteFloat64 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteFloat64(src []float64, dst Floating) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, float64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedFloat64 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedFloat64(src [][]float64, dst Floating) (written int) {
	mustSameChannels(dst.Channels(), len(src))
	// determine the length of longest nested slice
	for i := range src {
		if len(src[i]) > written {
			written = len(src[i])
		}
	}
	// limit a number of writes to the length of the buffer
	written = min(written, dst.Length())
	for c := 0; c < dst.Channels(); c++ {
		for i := 0; i < written; i++ {
			if i < len(src[c]) {
				dst.SetSample(BufferIndex(dst.Channels(), c, i), float64(src[c][i]))
			} else {
				dst.SetSample(BufferIndex(dst.Channels(), c, i), 0)
			}
		}
	}
	return
}
