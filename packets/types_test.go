package packets

import (
	"bytes"
	"testing"
)

func TestBooleanEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	var b = Boolean(true)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x01}, buf.Bytes()) {
		t.Fatalf("expected {0x01}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Boolean(false)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00}, got %#v", buf.Bytes())
	}
}

func TestBooleanDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00})

	var b Boolean
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != false {
		t.Fatalf("expected false, got %#v", b)
	}

	buf = bytes.NewBuffer([]byte{0x01})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != true {
		t.Fatalf("expected true, got %#v", b)
	}

	buf = bytes.NewBuffer([]byte{0x42})
	if err := b.Decode(buf); err == nil {
		t.Fatal("expected err, got nil")
	}
}

func TestByteEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	var b = Byte(0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Byte(-128)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x80}, buf.Bytes()) {
		t.Fatalf("expected {0x80}, got %#v", buf.Bytes())
	}
}

func TestByteDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00})

	var b Byte
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0x00 {
		t.Fatalf("expected 0x00, got % x", b)
	}

	buf = bytes.NewBuffer([]byte{0x7f})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0x7f {
		t.Fatalf("expected 0x7f, got % x", b)
	}
}

func TestUByteEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	var b = UByte(0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = UByte(255)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0xff}, buf.Bytes()) {
		t.Fatalf("expected {0xff}, got %#v", buf.Bytes())
	}
}

func TestUByteDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00})

	var b UByte
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0x00 {
		t.Fatalf("expected 0x00, got % x", b)
	}

	buf = bytes.NewBuffer([]byte{0xff})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0xff {
		t.Fatalf("expected 0xff, got % x", b)
	}
}

func TestShortEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	var b = Short(0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Short(-32768)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x80, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x80, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Short(32767)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x7f, 0xff}, buf.Bytes()) {
		t.Fatalf("expected {0x7f, 0xff}, got %#v", buf.Bytes())
	}

}

func TestShortDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00, 0x00})

	var b Short
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0 {
		t.Fatalf("expected 0, got %d", b)
	}

	buf = bytes.NewBuffer([]byte{0x80, 0x00})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != -32768 {
		t.Fatalf("expected -32768, got %d", b)
	}

	buf = bytes.NewBuffer([]byte{0x7f, 0xff})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 32767 {
		t.Fatalf("expected 32767, got %d", b)
	}
}
