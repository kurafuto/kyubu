package chunk

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"testing"
)

func compressBytes(data []byte) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	gz := gzip.NewWriter(b)

	if err := binary.Write(gz, binary.BigEndian, int32(len(data))); err != nil {
		return nil, fmt.Errorf("binary write error: %s", err.Error())
	}

	n, err := gz.Write(data)
	if err != nil {
		return nil, fmt.Errorf("gzip write error: %s", err.Error())
	}
	if n != len(data) {
		return nil, fmt.Errorf("gzip write: wrote %d bytes, expected %d bytes", n, len(data))
	}

	gz.Flush()
	return b.Bytes(), nil
}

func TestDecompress(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	compressed, err := compressBytes(data)
	if err != nil {
		t.Error(err)
	}

	ch, err := Decompress(compressed, 2, 2, 2)
	if err != nil {
		t.Fatal("error decompressing chunk:", err)
	}

	if ch.Length != len(data) {
		t.Fatalf("expected length of %d, got %d", len(data), ch.Length)
	}

	if !bytes.Equal(ch.Data, data) {
		t.Fatalf("expected data: %#v, got %#v", data, ch.Data)
	}
}

func TestCompress(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	compressed, err := compressBytes(data)
	if err != nil {
		t.Error(err)
	}

	ch := Chunk{
		Length: len(data),
		Data:   data,
		X:      2,
		Y:      2,
		Z:      2,
	}

	comp, err := ch.Compress()
	if err != nil {
		t.Fatal("Chunk compress error:", err)
	}

	if !bytes.Equal(compressed, comp) {
		t.Fatalf("invalid compressed data; expected %#v, got %#v", compressed, comp)
	}
}

func TestIncorrectLength(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}
	compressed, err := compressBytes(data)
	if err != nil {
		t.Error(err)
	}

	_, err = Decompress(compressed, 2, 2, 2)
	if err == nil {
		t.Fatal("Chunk decompressed 8 byte level from %d bytes", len(data))
	}
}
