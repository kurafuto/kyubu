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

func TestUShortEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	var b = UShort(0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = UShort(65535)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0xff, 0xff}, buf.Bytes()) {
		t.Fatalf("expected {0xff, 0xff}, got %#v", buf.Bytes())
	}

}

func TestUShortDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00, 0x00})

	var b UShort
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0 {
		t.Fatalf("expected 0, got %d", b)
	}

	buf = bytes.NewBuffer([]byte{0xff, 0xff})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 65535 {
		t.Fatalf("expected 32767, got %d", b)
	}
}

func TestIntEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	var b = Int(0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00, 0x00, 0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00, 0x00, 0x00, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Int(-2147483648)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x80, 0x00, 0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x80, 0x00, 0x00, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Int(2147483647)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x7f, 0xff, 0xff, 0xff}, buf.Bytes()) {
		t.Fatalf("expected {0x7f, 0xff, 0xff, 0xff}, got %#v", buf.Bytes())
	}

}

func TestIntDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	var b Int
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0 {
		t.Fatalf("expected 0, got %d", b)
	}

	buf = bytes.NewBuffer([]byte{0x80, 0x00, 0x00, 0x00})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != -2147483648 {
		t.Fatalf("expected -2147483648, got %d", b)
	}

	buf = bytes.NewBuffer([]byte{0x7f, 0xff, 0xff, 0xff})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 2147483647 {
		t.Fatalf("expected 2147483647, got %d", b)
	}
}

func TestLongEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	var b = Long(0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Long(-9223372036854775808)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Long(9223372036854775807)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, buf.Bytes()) {
		t.Fatalf("expected {0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, got %#v", buf.Bytes())
	}

}

func TestLongDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	var b Long
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0 {
		t.Fatalf("expected 0, got %d", b)
	}

	buf = bytes.NewBuffer([]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != -9223372036854775808 {
		t.Fatalf("expected -9223372036854775808, got %d", b)
	}

	buf = bytes.NewBuffer([]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 9223372036854775807 {
		t.Fatalf("expected 9223372036854775807, got %d", b)
	}
}

func TestFloatEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	var b = Float(0.0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00, 0x00, 0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00, 0x00, 0x00, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Float(2139095039.0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x4e, 0xff, 0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x4e, 0xff, 0x00, 0x00}, got %#v", buf.Bytes())
	}
}

func TestFloatDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00})

	var b Float
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0.0 {
		t.Fatalf("expected 0.0, got %f", b)
	}

	buf = bytes.NewBuffer([]byte{0x4e, 0xff, 0x00, 0x00})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 2139095039.0 {
		t.Fatalf("expected 2139095039.0, got %f", b)
	}
}

func TestDoubleEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

	var b = Double(0.0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, got %#v", buf.Bytes())
	}

	buf.Reset()
	b = Double(3423266188008119638567638907689404624415750050196511731360612015004450816.0)
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte{0x4e, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, buf.Bytes()) {
		t.Fatalf("expected {0x4e, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, got %#v", buf.Bytes())
	}
}

func TestDoubleDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	var b Double
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 0.0 {
		t.Fatalf("expected 0.0, got %f", b)
	}

	buf = bytes.NewBuffer([]byte{0x4e, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	if b != 3423266188008119638567638907689404624415750050196511731360612015004450816.0 {
		t.Fatalf("expected 3423266188008119638567638907689404624415750050196511731360612015004450816.0, got %f", b)
	}
}

func TestStringEncode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

	var b = String("Hello, world!")
	if err := b.Encode(buf); err != nil {
		t.Fatal(err)
	}

	exp := []byte{0x1a, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21}
	if !bytes.Equal(exp, buf.Bytes()) {
		t.Fatalf("expected %v, got %#v", exp, buf.Bytes())
	}
}

func TestStringDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x1a, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21})

	var b String
	if err := b.Decode(buf); err != nil {
		t.Fatal(err)
	}

	exp := String("Hello, world!")
	if b != exp {
		t.Fatalf("expected %q, got %q", exp, b)
	}
}
