package main

import (
	"net/http/httptest"
	"testing"
)

func Test_ensure_WriteHeader_writes_to_embedded_struct(t *testing.T) {
	rec := httptest.NewRecorder()
	bw := &ByteWriter{rec, 0, 0}
	bw.WriteHeader(404)

	expectedCode := 404
	if bw.Status != expectedCode {
		t.Fatalf("bw.Status = %v, want %v", bw.Status, expectedCode)
	}

	if rec.Code != expectedCode {
		t.Fatalf("rec.Code = %v, want %v", rec.Code, expectedCode)
	}
}
