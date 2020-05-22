// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-23 00:11:27.63103 +0200 CEST m=+0.013658257
package signal

import "math"

// Int16 is int16 signed fixed signal.
type Int16 struct {
	buffer []int16
	channels
	bitDepth
}

// Int16 allocates a new sequential int16 signal buffer.
func (a Allocator) Int16(bd BitDepth) Signed {
	return Int16{
		buffer:   make([]int16, 0, a.Capacity*a.Channels),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, BitDepth16),
	}
}

// GetInt16 selects a new sequential int16 signal buffer.
// from the pool.
func (p *Pool) GetInt16(bd BitDepth) Signed {
	if p == nil {
		return nil
	}
	return p.i16.Get().(Signed).setBitDepth(bd)
}

// PutInt16 places signal buffer back to the pool. If a type of
// provided buffer isn't Int16 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (p *Pool) PutInt16(s Signed) {
	if p == nil {
		return
	}
	if _, ok := s.(Int16); !ok {
		panic("pool put int16 invalid type")
	}
	mustSameCapacity(s.Capacity(), p.allocator.Capacity)
	p.i16.Put(s.Reset())
}

func (s Int16) setBitDepth(bd BitDepth) Signed {
	s.bitDepth = limitBitDepth(bd, BitDepth16)
	return s
}

// Capacity returns capacity of a single channel.
func (s Int16) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Int16) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s Int16) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Int16) Len() int {
	return len(s.buffer)
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s Int16) AppendSample(value int64) Signed {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, int16(s.BitDepth().SignedValue(value)))
	return s
}

// Sample returns signal value for provided channel and position.
func (s Int16) Sample(pos int) int64 {
	return int64(s.buffer[pos])
}

// SetSample sets sample value for provided position.
func (s Int16) SetSample(pos int, value int64) {
	s.buffer[pos] = int16(s.BitDepth().SignedValue(value))
}

// Slice slices buffer with respect to channels.
func (s Int16) Slice(start, end int) Signed {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Reset sets length of the buffer to zero.
func (s Int16) Reset() Signed {
	return s.Slice(0, 0)
}

// Append appends [0:Length] data from src to current buffer and returns new
// Signed buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s Int16) Append(src Signed) Signed {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with sources cap
		s.buffer = append(make([]int16, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Signed(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// ReadInt16 reads values from the buffer into provided slice.
func ReadInt16(src Signed, dst []int16) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = int16(BitDepth16.SignedValue(src.Sample(pos)))
	}
}

// ReadStripedInt16 reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStripedInt16(src Signed, dst [][]int16) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = int16(BitDepth16.SignedValue(src.Sample(src.ChannelPos(channel, pos))))
		}
	}
}

// WriteInt16 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func WriteInt16(src []int16, dst Signed) Signed {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(int64(src[pos]))
	}
	return dst
}

// WriteStripedInt16 appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be appended. Sample values are capped by maximum value
// of the buffer bit depth.
func WriteStripedInt16(src [][]int16, dst Signed) Signed {
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
				dst = dst.AppendSample(int64(src[channel][pos]))
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}
