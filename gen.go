// +build ignore

// This program generates all signal types.
package main

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

type typeGenerator struct {
	InterfaceProps
	Timestamp   time.Time
	Builtin     string
	Name        string
	Pool        string
	MaxBitDepth string // used for fixed-point types only
}

type InterfaceProps struct {
	Interface  string
	SampleType string
}

func main() {
	var (
		signedProps = InterfaceProps{
			Interface:  "Signed",
			SampleType: "int64",
		}
		unsignedProps = InterfaceProps{
			Interface:  "Unsigned",
			SampleType: "uint64",
		}
		floatingProps = InterfaceProps{
			Interface:  "Floating",
			SampleType: "float64",
		}
	)
	types := map[typeGenerator]templates{
		{
			InterfaceProps: signedProps,
			Builtin:        "int8",
			Name:           "Int8",
			Pool:           "i8",
			MaxBitDepth:    "BitDepth8",
		}: fixedTemplates,
		{
			InterfaceProps: signedProps,
			Builtin:        "int16",
			Name:           "Int16",
			Pool:           "i16",
			MaxBitDepth:    "BitDepth16",
		}: fixedTemplates,
		{
			InterfaceProps: signedProps,
			Builtin:        "int32",
			Name:           "Int32",
			Pool:           "i32",
			MaxBitDepth:    "BitDepth32",
		}: fixedTemplates,
		{
			InterfaceProps: signedProps,
			Builtin:        "int64",
			Name:           "Int64",
			Pool:           "i64",
			MaxBitDepth:    "BitDepth64",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint8",
			Name:           "Uint8",
			Pool:           "u8",
			MaxBitDepth:    "BitDepth8",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint16",
			Name:           "Uint16",
			Pool:           "u16",
			MaxBitDepth:    "BitDepth16",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint32",
			Name:           "Uint32",
			Pool:           "u32",
			MaxBitDepth:    "BitDepth32",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint64",
			Name:           "Uint64",
			Pool:           "u64",
			MaxBitDepth:    "BitDepth64",
		}: fixedTemplates,
		{
			InterfaceProps: floatingProps,
			Builtin:        "float32",
			Name:           "Float32",
			Pool:           "f32",
		}: floatingTemplates,
		{
			InterfaceProps: floatingProps,
			Builtin:        "float64",
			Name:           "Float64",
			Pool:           "f64",
		}: floatingTemplates,
	}

	for gen, template := range types {
		generate(gen, template)
	}
}

func generate(gen typeGenerator, t templates) {
	gen.Timestamp = time.Now()

	generateFile(fmt.Sprintf("%s.go", gen.Builtin), gen, t.types)
	// generateFile(fmt.Sprintf("%s_test.go", gen.Builtin), gen, t.tests)

	// err = t.tests.Execute(f, gen)
	// die(fmt.Sprintf("execute %s tests template for %s type", t.tests.Name(), gen.Name), err)
}

func generateFile(fileName string, gen typeGenerator, t *template.Template) {
	if t == nil {
		return
	}
	f, err := os.Create(fileName)
	die(fmt.Sprintf("create %s file", fileName), err)
	defer f.Close()

	err = t.ExecuteTemplate(f, "signal", gen)
	die(fmt.Sprintf("execute %s template for %s type", t.Name(), gen.Name), err)
}

func die(reason string, err error) {
	if err != nil {
		panic(fmt.Sprintf("failed %s: %v", reason, err))
	}
}

type templates struct {
	types *template.Template
	tests *template.Template
}

var (
	floatingTemplates = templates{
		types: template.Must(template.New("floating").Parse(floating + base)),
		tests: template.Must(template.New("floating tests").Parse(floatingTests)),
	}
	fixedTemplates = templates{
		types: template.Must(template.New("fixed").Parse(fixed + base)),
		tests: template.Must(template.New("fixed tests").Parse(fixedTests)),
	}
)

const (
	base = `{{define "signal"}}package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}

import "math"
{{template "structures" .}}
// Put{{ .Name }} places signal buffer back to the pool. If a type of
// provided buffer isn't {{ .Name }} or its capacity doesn't equal
// allocator capacity, the function will panic.
func (p *Pool) Put{{ .Name }}(s {{ .Interface }}) {
	if p == nil {
		return
	}
	if _, ok := s.({{ .Name }}); !ok {
		panic("pool put {{ .Builtin }} invalid type")
	}
	mustSameCapacity(s.Capacity(), p.allocator.Capacity)
	p.{{ .Pool }}.Put(s.Slice(0, p.allocator.Length))
}

// Capacity returns capacity of a single channel.
func (s {{ .Name }}) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s {{ .Name }}) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s {{ .Name }}) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s {{ .Name }}) Len() int {
	return len(s.buffer)
}
{{template "appends" .}}
// Sample returns signal value for provided channel and position.
func (s {{ .Name }}) Sample(pos int) {{ .SampleType }} {
	return {{ .SampleType }}(s.buffer[pos])
}

// Slice slices buffer with respect to channels.
func (s {{ .Name }}) Slice(start, end int) {{ .Interface }} {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}
{{template "reads" .}}
// Write{{ .Name }} writes values from provided slice into the buffer.
func Write{{ .Name }}(src []{{ .Builtin }}, dst {{ .Interface }}) {{ .Interface }} {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample({{ .SampleType }}(src[pos]))
	}
	return dst
}

// WriteStriped{{ .Name }} appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// those channels will be appended.
func WriteStriped{{ .Name }}(src [][]{{ .Builtin }}, dst {{ .Interface }}) {{ .Interface }} {
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
				dst = dst.AppendSample({{ .SampleType }}(src[channel][pos]))
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}
{{end}}
`

	fixed = `{{define "structures"}}
// {{ .Name }} is {{ .Builtin }} {{ .Interface }} fixed-point signal.
type {{ .Name }} struct {
	buffer []{{ .Builtin }}
	channels
	bitDepth
}

// {{ .Name }} allocates a new sequential {{ .Builtin }} signal buffer.
func (a Allocator) {{ .Name }}(bd BitDepth) {{ .Interface }} {
	return {{ .Name }}{
		buffer:   make([]{{ .Builtin }}, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, {{ .MaxBitDepth }}),
	}
}

// Get{{ .Name }} selects a new sequential {{ .Builtin }} signal buffer.
// from the pool.
func (p *Pool) Get{{ .Name }}(bd BitDepth) {{ .Interface }} {
	if p == nil {
		return nil
	}
	return p.{{ .Pool }}.Get().({{ .Interface }}).setBitDepth(bd)
}

func (s {{ .Name }}) setBitDepth(bd BitDepth) {{ .Interface }} {
	s.bitDepth = limitBitDepth(bd, {{ .MaxBitDepth }})
	return s
}
{{end}}

{{define "appends"}}
// Append appends [0:Length] samples from src to current buffer and returns new
// {{ .Interface }} buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s {{ .Name }}) Append(src {{ .Interface }}) {{ .Interface }} {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]{{ .Builtin }}, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := {{ .Interface }}(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
// Sample values are capped by maximum value of the buffer bit depth.
func (s {{ .Name }}) AppendSample(value {{ .SampleType }}) {{ .Interface }} {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, {{ .Builtin }}(s.BitDepth().{{ .Interface }}Value(value)))
	return s
}

// SetSample sets sample value for provided position.
// Sample values are capped by maximum value of the buffer bit depth.
func (s {{ .Name }}) SetSample(pos int, value {{ .SampleType }}) {
	s.buffer[pos] = {{ .Builtin }}(s.BitDepth().{{ .Interface }}Value(value))
}
{{end}}

{{define "reads"}}
// Read{{ .Name }} reads values from the buffer into provided slice.
func Read{{ .Name }}(src {{ .Interface }}, dst []{{ .Builtin }}) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = {{ .Builtin }}({{ .MaxBitDepth }}.{{ .Interface }}Value(src.Sample(pos)))
	}
}

// ReadStriped{{ .Name }} reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStriped{{ .Name }}(src {{ .Interface }}, dst [][]{{ .Builtin }}) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = {{ .Builtin }}({{ .MaxBitDepth }}.{{ .Interface }}Value(src.Sample(src.ChannelPos(channel, pos))))
		}
	}
}
{{end}}
`

	floating = `{{define "structures"}}
// {{ .Name }} is {{ .Builtin }} floating-point signal.
type {{ .Name }} struct {
	buffer []{{ .Builtin }}
	channels
}

// {{ .Name }} allocates a new sequential {{ .Builtin }} signal buffer.
func (a Allocator) {{ .Name }}() {{ .Interface }} {
	return {{ .Name }}{
		buffer:   make([]{{ .Builtin }}, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
	}
}

// Get{{ .Name }} selects a new sequential {{ .Builtin }} signal buffer.
// from the pool.
func (p *Pool) Get{{ .Name }}() {{ .Interface }} {
	if p == nil {
		return nil
	}
	return p.{{ .Pool }}.Get().({{ .Interface }})
}
{{end}}

{{define "appends"}}
// Append appends [0:Length] samples from src to current buffer and returns new
// {{ .Interface }} buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s {{ .Name }}) Append(src {{ .Interface }}) {{ .Interface }} {
	mustSameChannels(s.Channels(), src.Channels())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]{{ .Builtin }}, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := {{ .Interface }}(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s {{ .Name }}) AppendSample(value {{ .SampleType }}) {{ .Interface }} {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, {{ .Builtin }}(value))
	return s
}

// SetSample sets sample value for provided position.
func (s {{ .Name }}) SetSample(pos int, value {{ .SampleType }}) {
	s.buffer[pos] = {{ .Builtin }}(value)
}
{{end}}

{{define "reads"}}
// Read{{ .Name }} reads values from the buffer into provided slice.
func Read{{ .Name }}(src {{ .Interface }}, dst []{{ .Builtin }}) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = {{ .Builtin }}(src.Sample(pos))
	}
}

// ReadStriped{{ .Name }} reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStriped{{ .Name }}(src {{ .Interface }}, dst [][]{{ .Builtin }}) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = {{ .Builtin }}(src.Sample(src.ChannelPos(channel, pos)))
		}
	}
}
{{end}}
`

	fixedTests = `package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}

import (
	"testing"

	"pipelined.dev/signal"
)

func Test{{ .Name }}(t *testing.T) {
	t.Run("{{ .Builtin }}", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.{{ .Name }}(signal.{{ .MaxBitDepth }}).
			Append(signal.WriteStriped{{ .Name }}(
				[][]{{ .Builtin }}{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.{{ .Name }}(signal.{{ .MaxBitDepth }})),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]{{ .Builtin }}{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}
`
	floatingTests = `package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}

import (
	"testing"

	"pipelined.dev/signal"
)

func Test{{ .Name }}(t *testing.T) {
	t.Run("{{ .Builtin }}", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.{{ .Name }}().
			Append(signal.WriteStriped{{ .Name }}(
				[][]{{ .Builtin }}{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.{{ .Name }}()),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]{{ .Builtin }}{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}
`
)
