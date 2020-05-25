// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-25 04:38:56.136346 +0200 CEST m=+0.017784743
package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func TestFloat64(t *testing.T) {
	t.Run("float64", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Float64().
			Append(signal.WriteStripedFloat64(
				[][]float64{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.Float64()),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]float64{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}