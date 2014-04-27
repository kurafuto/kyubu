package packets

import (
	"bytes"
	"io"
	"testing"
)

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
	p := NewParser(bytes.NewBuffer(b))

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

func TestParseFewBytes(t *testing.T) {
	// Packet 0x13, Message, 66 bytes
	b := bytes.NewBuffer([]byte{0x0d, 0x00, 0x00, 0x00})
	p := NewParser(b)

	if _, err := p.Next(); err == nil {
		t.Fatal("expected EOF, got nil")
	}
}

func TestInvalidPacketId(t *testing.T) {
	b := bytes.NewBuffer([]byte{0xff, 0x00, 0x00, 0x00})
	p := NewParser(b)

	if n, err := p.Next(); err == nil {
		t.Fatalf("expected err for 0xff, got nil and packet: %#v", n)
	}
}

// This will always return an ErrClosedPipe on Read().
type errReader struct{}

func (r *errReader) Read(b []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestParseIdErr(t *testing.T) {
	p := NewParser(&errReader{})

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
	p := NewParser(&lessReader{})

	_, err := p.Next()
	expected := "kyubu: Read failed for id, wanted 1 bytes, got 0"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

// This writes `Id` in ONCE if len(b) == 1, otherwise it returns given Err.
type idOnceReader struct {
	Id   byte
	done bool
	Err  error
}

func (r *idOnceReader) Read(b []byte) (int, error) {
	if r.done {
		return 0, r.Err
	}
	r.done = true
	b[0] = r.Id
	return 1, nil
}

func TestEOFReadingBody(t *testing.T) {
	p := NewParser(&idOnceReader{Id: 0x0d, Err: io.EOF})

	_, err := p.Next()
	expected := "kyubu: EOF reading packet 0x0d, got 0"
	if err == nil {
		t.Fatal("expected err, got nil")
	}
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestErrReadingBody(t *testing.T) {
	p := NewParser(&idOnceReader{Id: 0x0d, Err: io.ErrClosedPipe})

	_, err := p.Next()
	if err == nil {
		t.Fatal("expected err, got nil")
	}
	if err != io.ErrClosedPipe {
		t.Fatalf("expected io.ErrClosedPipe, got %q", err.Error())
	}
}
