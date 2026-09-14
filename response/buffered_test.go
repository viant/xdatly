package response

import (
	"io"
	"net/http"
	"testing"
)

func TestBufferedImplementsPublicTransportContracts(t *testing.T) {
	var _ Response = (*Buffered)(nil)
	var _ Compressed = (*Buffered)(nil)
	var _ StatusCoder = (*Buffered)(nil)
}

func TestNewBuffered(t *testing.T) {
	actual := NewBuffered(
		WithStatusCode(201),
		WithBytes([]byte("payload")),
		WithHeader("Content-Type", "text/plain"),
		WithCompressionType("gzip"),
	)

	if actual.StatusCode() != 201 {
		t.Fatalf("expected status 201, got %d", actual.StatusCode())
	}
	if actual.Size() != len("payload") {
		t.Fatalf("expected size %d, got %d", len("payload"), actual.Size())
	}
	if actual.Headers().Get("Content-Type") != "text/plain" {
		t.Fatalf("expected content type header")
	}
	if actual.CompressionType() != "gzip" {
		t.Fatalf("expected compression type gzip, got %q", actual.CompressionType())
	}
	body, err := io.ReadAll(actual.Body())
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if string(body) != "payload" {
		t.Fatalf("expected payload body, got %q", string(body))
	}
}

func TestBufferedWrite(t *testing.T) {
	actual := NewBuffered(WithHeaders(http.Header{"X-Test": []string{"true"}}))

	if _, err := actual.Write([]byte("abc")); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if actual.Size() != 3 {
		t.Fatalf("expected size 3, got %d", actual.Size())
	}
	body, err := io.ReadAll(actual.Body())
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if string(body) != "abc" {
		t.Fatalf("expected abc body, got %q", string(body))
	}
	if actual.Headers().Get("X-Test") != "true" {
		t.Fatalf("expected header copy")
	}
}
