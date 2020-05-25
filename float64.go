// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-25 03:55:17.436942 +0200 CEST m=+0.016198025
package signal

import "math"

// Float64 is a sequential float64 floating-point signal.
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
	if _, ok := s.(Float64); !ok {
		panic("pool put float64 invalid type")
	}
	mustSameCapacity(s.Capacity(), p.allocator.Capacity)
	p.f64.Put(s.Slice(0, p.allocator.Length))
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

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s Float64) AppendSample(value float64) Floating {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, float64(value))
	return s
}

// Sample returns signal value for provided channel and position.
func (s Float64) Sample(pos int) float64 {
	return float64(s.buffer[pos])
}

// SetSample sets sample value for provided position.
func (s Float64) SetSample(pos int, value float64) {
	s.buffer[pos] = float64(value)
}

// Slice slices buffer with respect to channels.
func (s Float64) Slice(start, end int) Floating {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Reset sets length of the buffer to zero.
func (s Float64) Reset() Floating {
	return s.Slice(0, 0)
}

// Append appends data from src buffer to the end of the buffer.
func (s Float64) Append(src Floating) Floating {
	mustSameChannels(s.Channels(), src.Channels())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]float64, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Floating(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// ReadFloat64 reads values from the buffer into provided slice.
func ReadFloat64(src Floating, dst []float64) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = float64(src.Sample(pos))
	}
}

// ReadStripedFloat64 reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStripedFloat64(src Floating, dst [][]float64) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = float64(src.Sample(src.ChannelPos(channel, pos)))
		}
	}
}

// WriteFloat64 writes values from provided slice into the buffer.
func WriteFloat64(src []float64, dst Floating) Floating {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(float64(src[pos]))
	}
	return dst
}

// WriteStripedFloat64 appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be appended.
func WriteStripedFloat64(src [][]float64, dst Floating) Floating {
	mustSameChannels(dst.Channels(), len(src))
	var length int
	for i := range src {
		if len(src[i]) > length {
			length = len(src[i])
		}
	}
	length = min(length, dst.Capacity()-dst.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < dst.Channels(); channel++ {
			if pos < len(src[channel]) {
				dst = dst.AppendSample(float64(src[channel][pos]))
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}