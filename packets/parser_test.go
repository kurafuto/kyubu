package packets

import (
	"bytes"
	"io"
	"testing"
)

/*
// TODO: Cyclic import errors on NewIdentification. Hmm
func TestParseToEOF(t *testing.T) {
	ident, err := NewIdentification("test", "test")
	if err != nil {
		t.Error("packet init error:", err)
	}

	level, err := NewLevelInitialize()
	if err != nil {
		t.Error("packet init error", err)
	}

	b := append(ident.Bytes(), level.Bytes()...)
	p := NewParser(bytes.NewBuffer(b), ClientBound)

	n, err := p.Recv()
	if err != nil {
		t.Fatal("round 1: parser.Recv() error:", err)
	}
	if n.Id() != 0x00 {
		t.Fatalf("round 1: expected packet 0x00, got %#.2x", n.Id())
	}

	n, err = p.Recv()
	if err != nil {
		t.Fatal("round 2: parser.Recv() error:", err)
	}
	if n.Id() != 0x01 {
		t.Fatalf("round 2: expected packet 0x01, got %#.2x", n.Id())
	}

	_, err = p.Recv()
	if err != io.EOF {
		t.Fatal("round 3: expected EOF, got", err)
	}
}
*/

func TestParseFewBytes(t *testing.T) {
	// 0x03 -> Time Update
	b := bytes.NewBuffer([]byte{0x03, 0x00, 0x00, 0x00})
	p := NewParser(b, Play, ClientBound)

	if _, err := p.Recv(); err == nil {
		t.Fatal("expected EOF, got nil")
	}
}

func TestInvalidPacketId(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xff, 0x00, 0x00, 0x00})
	p := NewParser(b, Play, ClientBound)

	if n, err := p.Recv(); err == nil {
		t.Fatalf("expected err for 0xff, got nil and packet: %#v", n)
	}
}

type errReaderWriter struct{ Err error }

func (r *errReaderWriter) Read(b []byte) (int, error) {
	return 0, r.Err
}

func (r *errReaderWriter) Write(p []byte) (int, error) {
	return 0, r.Err
}

func TestParseIdErr(t *testing.T) {
	p := NewParser(&errReaderWriter{io.ErrClosedPipe}, Play, ClientBound)

	if _, err := p.Recv(); err != io.ErrClosedPipe {
		t.Fatalf("expected io.ErrClosedPipe, got %#v", err)
	}
}

// This will always return one less byte than they want.
type lessReaderWriter struct{}

func (r *lessReaderWriter) Read(b []byte) (int, error) {
	return len(b) - 1, nil
}

func (r *lessReaderWriter) Write(p []byte) (int, error) {
	return len(p) - 1, nil
}

func TestParseIdLessBytes(t *testing.T) {
	p := NewParser(&lessReaderWriter{}, Play, ClientBound)

	_, err := p.Recv()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %s", err)
	}
}
