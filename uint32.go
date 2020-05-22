// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-22 19:53:07.286878 +0200 CEST m=+0.009779530
package signal

// Uint32 is uint32 signed fixed signal.
type Uint32 struct {
	buffer []uint32
	channels
	bitDepth
}

// Uint32 allocates new sequential uint32 signal buffer.
func (a Allocator) Uint32(bd BitDepth) Unsigned {
	return Uint32{
		buffer:   make([]uint32, 0, a.Capacity*a.Channels),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, BitDepth64),
	}
}

func (s Uint32) setBitDepth(bd BitDepth) Unsigned {
	s.bitDepth = limitBitDepth(bd, BitDepth32)
	return s
}

// Capacity returns capacity of a single channel.
func (s Uint32) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Uint32) Length() int {
	return len(s.buffer) / int(s.channels)
}

// Cap returns capacity of whole buffer.
func (s Uint32) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Uint32) Len() int {
	return len(s.buffer)
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s Uint32) AppendSample(value uint64) Unsigned {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, uint32(s.BitDepth().UnsignedValue(value)))
	return s
}

// Sample returns signal value for provided channel and position.
func (s Uint32) Sample(pos int) uint64 {
	return uint64(s.buffer[pos])
}

// SetSample sets sample value for provided position.
func (s Uint32) SetSample(pos int, value uint64) {
	s.buffer[pos] = uint32(s.BitDepth().UnsignedValue(value))
}

// Slice slices buffer with respect to channels.
func (s Uint32) Slice(start, end int) Unsigned {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Reset sets length of the buffer to zero.
func (s Uint32) Reset() Unsigned {
	return s.Slice(0, 0)
}

// Append appends data from src to current buffer and returns new
// Unsigned buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic.  If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s Uint32) Append(src Unsigned) Unsigned {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with sources cap
		s.buffer = append(make([]uint32, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Unsigned(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// ReadUint32 reads values from the buffer into provided slice.
func ReadUint32(src Unsigned, dst []uint32) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = uint32(BitDepth32.UnsignedValue(src.Sample(pos)))
	}
}

// ReadStripedUint32 reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStripedUint32(src Unsigned, dst [][]uint32) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = uint32(BitDepth32.UnsignedValue(src.Sample(src.ChannelPos(channel, pos))))
		}
	}
}

// WriteUint32 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func WriteUint32(src []uint32, dst Unsigned) Unsigned {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(uint64(src[pos]))
	}
	return dst
}

// WriteStripedUint32 appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be appended. Sample values are capped by maximum value
// of the buffer bit depth.
func WriteStripedUint32(src [][]uint32, dst Unsigned) Unsigned {
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
				dst = dst.AppendSample(uint64(src[channel][pos]))
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}
