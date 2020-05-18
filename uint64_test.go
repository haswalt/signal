// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-18 20:18:05.384035 +0200 CEST m=+0.010622843
package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func TestUint64(t *testing.T) {
	t.Run("uint64", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Uint64(signal.BitDepth64).
			Append(signal.WriteStripedUint64(
				[][]uint64{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.Uint64(signal.BitDepth64)),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]uint64{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}