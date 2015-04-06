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

var statusTable = []struct {
	code     int
	expected uint64
	actual   func(*RuntimeStats) uint64
}{
	{100, 1, func(r *RuntimeStats) uint64 { return r.Status1xx() }},
	{102, 1, func(r *RuntimeStats) uint64 { return r.Status1xx() }},
	{200, 1, func(r *RuntimeStats) uint64 { return r.Status2xx() }},
	{226, 1, func(r *RuntimeStats) uint64 { return r.Status2xx() }},
	{300, 1, func(r *RuntimeStats) uint64 { return r.Status3xx() }},
	{308, 1, func(r *RuntimeStats) uint64 { return r.Status3xx() }},
	{400, 1, func(r *RuntimeStats) uint64 { return r.Status4xx() }},
	{499, 1, func(r *RuntimeStats) uint64 { return r.Status4xx() }},
	{500, 1, func(r *RuntimeStats) uint64 { return r.Status5xx() }},
	{599, 1, func(r *RuntimeStats) uint64 { return r.Status5xx() }},
}

func Test_IncStatus(t *testing.T) {
	for _, tt := range statusTable {
		rs := &RuntimeStats{}
		err := rs.IncStatus(tt.code)

		if err != nil {
			t.Fatalf("rs.IncStatus(%v) = %v", tt.code, err)
		}

		if tt.expected != tt.actual(rs) {
			t.Fatalf("rs.IncStatus(%v) incremented value to %v, want %v", tt.code, tt.actual(rs), tt.expected)
		}
	}
}

var invalidStatusTable = []struct {
	code   int
	errMsg string
}{
	{99, "Unexpected response code 99."},
	{600, "Unexpected response code 600."},
}

func Test_IncStatus_with_invalid_status_codes(t *testing.T) {
	for _, tt := range invalidStatusTable {
		rs := &RuntimeStats{}
		err := rs.IncStatus(tt.code)
		if err == nil {
			t.Fatalf("rs.IncStatus(%v) = nil, want error", tt.code)
		}

		if err.Error() != tt.errMsg {
			t.Fatalf("err.Error() = %v, want %v", err.Error(), tt.errMsg)
		}
	}
}
