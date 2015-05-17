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

	n, err := p.Next()
	if err != nil {
		t.Fatal("round 1: parser.Next() error:", err)
	}
	if n.Id() != 0x00 {
		t.Fatalf("round 1: expected packet 0x00, got %#.2x", n.Id())
	}

	n, err = p.Next()
	if err != nil {
		t.Fatal("round 2: parser.Next() error:", err)
	}
	if n.Id() != 0x01 {
		t.Fatalf("round 2: expected packet 0x01, got %#.2x", n.Id())
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Fatal("round 3: expected EOF, got", err)
	}
}
*/

func TestParseFewBytes(t *testing.T) {
	// Packet 0x13, Message, 66 bytes
	b := bytes.NewBuffer([]byte{0x0d, 0x00, 0x00, 0x00})
	p := NewParser(b, Anomalous)

	if _, err := p.Next(); err == nil {
		t.Fatal("expected EOF, got nil")
	}
}

func TestInvalidPacketId(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xff, 0x00, 0x00, 0x00})
	p := NewParser(b, Anomalous)

	if n, err := p.Next(); err == nil {
		t.Fatalf("expected err for 0xff, got nil and packet: %#v", n)
	}
}

// This will always return an ErrClosedPipe on Read().
type errReader struct{ Err error }

func (r *errReader) Read(b []byte) (int, error) {
	return 0, r.Err
}

func TestParseIdErr(t *testing.T) {
	p := NewParser(&errReader{io.ErrClosedPipe}, Anomalous)

	if _, err := p.Next(); err != io.ErrClosedPipe {
		t.Fatalf("expected io.ErrClosedPipe, got %#v", err)
	}
}

// This will always return one less byte than they want.
type lessReader struct{}

func (r *lessReader) Read(b []byte) (int, error) {
	return len(b) - 1, nil
}

func TestParseIdLessBytes(t *testing.T) {
	p := NewParser(&lessReader{}, Anomalous)

	_, err := p.Next()
	expected := "kyubu: Read failed for id, wanted 1 bytes, got 0"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}
