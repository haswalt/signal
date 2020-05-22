// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-22 19:17:55.874948 +0200 CEST m=+0.009245976
package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func TestFloat32(t *testing.T) {
	t.Run("float32", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Float32().
			Append(signal.WriteStripedFloat32(
				[][]float32{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.Float32()),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]float32{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}