package main

import "testing"

var byteTable = []struct {
	b        uint64
	expected string
}{
	{1023, "1023 B"},
	{1024, "1 KB"},
	{1024 * 1024, "1 MB"},
	{1024 * 1024 * 1024, "1 GB"},
}

func Test_toByte(t *testing.T) {
	rs := &RuntimeStats{}

	for _, tt := range byteTable {
		actual := rs.toBytes(tt.b)

		if actual != tt.expected {
			t.Fatalf("rs.toBytes(%v) = %v, want %v", tt.b, actual, tt.expected)
		}
	}
}
