package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-07-25 14:10:01.59973 +0200 CEST m=+0.021741616

import "math"

// Int32 is int32 Signed fixed-point signal.
type Int32 struct {
	buffer []int32
	channels
	bitDepth
}

// Int32 allocates a new sequential int32 signal buffer.
func (a Allocator) Int32(bd BitDepth) Signed {
	return Int32{
		buffer:   make([]int32, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, BitDepth32),
	}
}

func (s Int32) setBitDepth(bd BitDepth) Signed {
	s.bitDepth = limitBitDepth(bd, BitDepth32)
	return s
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
// Sample values are capped by maximum value of the buffer bit depth.
func (s Int32) AppendSample(value int64) Signed {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, int32(s.BitDepth().SignedValue(value)))
	return s
}

// SetSample sets sample value for provided index.
// Sample values are capped by maximum value of the buffer bit depth.
func (s Int32) SetSample(i int, value int64) {
	s.buffer[i] = int32(s.BitDepth().SignedValue(value))
}

// GetInt32 selects a new sequential int32 signal buffer.
// from the pool.
func (p *Pool) GetInt32(bd BitDepth) Signed {
	if p == nil {
		return nil
	}
	return p.i32.Get().(Signed).setBitDepth(bd)
}

// PutInt32 places signal buffer back to the pool. If a type of
// provided buffer isn't Int32 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (p *Pool) PutInt32(s Signed) {
	if p == nil {
		return
	}
	if sig, ok := s.(Int32); ok {
		capacity := cap(sig.buffer)
		for i, buf := 0, sig.buffer[:capacity]; i < capacity; i++ {
			buf[i] = 0
		}
	} else {
		panic("pool put int32 invalid type")
	}
	mustSameCapacity(s.Capacity(), p.allocator.Capacity)
	if s.Length() != p.allocator.Length {
		s = s.Slice(0, p.allocator.Length)
	}
	p.i32.Put(s)
}

// Capacity returns capacity of a single channel.
func (s Int32) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Int32) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s Int32) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Int32) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided channel and index.
func (s Int32) Sample(i int) int64 {
	return int64(s.buffer[i])
}

// Append appends [0:Length] samples from src to current buffer and returns new
// Signed buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s Int32) Append(src Signed) Signed {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]int32, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Signed(s)
	for i := 0; i < src.Len(); i++ {
		result = result.AppendSample(src.Sample(i))
	}
	return result
}

// Slice slices buffer with respect to channels.
func (s Int32) Slice(start, end int) Signed {
	start = BufferIndex(s.Channels(), 0, start)
	end = BufferIndex(s.Channels(), 0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// ReadInt32 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadInt32(src Signed, dst []int32) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = int32(BitDepth32.SignedValue(src.Sample(i)))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedInt32 reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedInt32(src Signed, dst [][]int32) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = int32(BitDepth32.SignedValue(src.Sample(BufferIndex(src.Channels(), c, i))))
		}
	}
	return
}

// WriteInt32 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteInt32(src []int32, dst Signed) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, int64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedInt32 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedInt32(src [][]int32, dst Signed) (written int) {
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
				dst.SetSample(BufferIndex(dst.Channels(), c, i), int64(src[c][i]))
			} else {
				dst.SetSample(BufferIndex(dst.Channels(), c, i), 0)
			}
		}
	}
	return
}
