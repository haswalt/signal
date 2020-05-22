// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-22 19:17:55.874948 +0200 CEST m=+0.009245976
package signal

// Float32 is a sequential float32 floating-point signal.
type Float32 struct {
	buffer []float32
	channels
}

// Float32 allocates new sequential float32 signal buffer.
func (a Allocator) Float32() Floating {
	return Float32{
		buffer:   make([]float32, 0, a.Channels*a.Capacity),
		channels: channels(a.Channels),
	}
}

// Capacity returns capacity of a single channel.
func (s Float32) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Float32) Length() int {
	return len(s.buffer) / int(s.channels)
}

// Cap returns capacity of whole buffer.
func (s Float32) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Float32) Len() int {
	return len(s.buffer)
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s Float32) AppendSample(value float64) Floating {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, float32(value))
	return s
}

// Sample returns signal value for provided channel and position.
func (s Float32) Sample(pos int) float64 {
	return float64(s.buffer[pos])
}

// SetSample sets sample value for provided position.
func (s Float32) SetSample(pos int, value float64) {
	s.buffer[pos] = float32(value)
}

// Slice slices buffer with respect to channels.
func (s Float32) Slice(start, end int) Floating {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Reset sets length of the buffer to zero.
func (s Float32) Reset() Floating {
	return s.Slice(0, 0)
}

// Append appends data from src buffer to the end of the buffer.
func (s Float32) Append(src Floating) Floating {
	mustSameChannels(s.Channels(), src.Channels())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]float32, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Floating(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// ReadFloat32 reads values from the buffer into provided slice.
func ReadFloat32(src Floating, dst []float32) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = float32(src.Sample(pos))
	}
}

// ReadStripedFloat32 reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStripedFloat32(src Floating, dst [][]float32) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = float32(src.Sample(src.ChannelPos(channel, pos)))
		}
	}
}

// WriteFloat32 writes values from provided slice into the buffer.
func WriteFloat32(src []float32, dst Floating) Floating {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(float64(src[pos]))
	}
	return dst
}

// WriteStripedFloat32 appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be appended.
func WriteStripedFloat32(src [][]float32, dst Floating) Floating {
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