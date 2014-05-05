package packets

import (
	"fmt"
	"reflect"
)

// ReflectBytes uses reflection to create a byte array out of a packet.
// This can then be written down a net.Conn, or into a reader for later replay.
func ReflectBytes(p Packet) []byte {
	b := NewPacketWrapper([]byte{})
	vp := reflect.ValueOf(p)
	for i := 0; i < vp.NumField(); i++ {
		v := vp.Field(i)
		switch v.Kind() {
		case reflect.Uint8: // byte
			b.WriteByte(byte(v.Uint()))
		case reflect.Int8: // sbyte
			b.WriteSByte(int8(v.Int()))
		case reflect.Int16: // short
			b.WriteShort(int16(v.Int()))
		case reflect.String: // string
			b.WriteString(v.String())
		case reflect.Slice: // []byte
			b.WriteBytes(v.Bytes())
		}
	}
	return b.Buffer.Bytes()
}

// ReflectRead unmarshals the bytes from b into v. Similar to encoding/json.Umarshal
//   var p MyPacket
//   err := ReflectRead(b, &p)
func ReflectRead(b []byte, v Packet) error {
	buf := NewPacketWrapper(b)
	vp := reflect.ValueOf(v).Elem()
	for i := 0; i < vp.NumField(); i++ {
		v := vp.Field(i)
		switch v.Kind() {
		case reflect.Uint8: // byte
			if val, err := buf.ReadByte(); err != nil {
				return err
			} else {
				v.SetUint(uint64(val))
			}
		case reflect.Int8: // sbyte
			if val, err := buf.ReadSByte(); err != nil {
				return err
			} else {
				v.SetInt(int64(val))
			}
		case reflect.Int16: // short
			if val, err := buf.ReadShort(); err != nil {
				return err
			} else {
				v.SetInt(int64(val))
			}
		case reflect.String: // string
			if val, err := buf.ReadString(); err != nil {
				return err
			} else {
				v.SetString(val)
			}
		case reflect.Slice:
			if val, err := buf.ReadBytes(); err != nil {
				return err
			} else {
				v.SetBytes(val)
			}
		}
	}
	if buf.Buffer.Len() != 0 {
		fmt.Printf("expected 0 leftover, had %d\n", buf.Buffer.Len())
		//return fmt.Errorf("kyubu: expected 0 leftover bytes, but had %d", buf.Buffer.Len())
	}
	return nil
}

// ReflectSize uses reflection to figure out the size of a given packet based
// on the fields it defines.
func ReflectSize(p Packet) (size int) {
	vp := reflect.ValueOf(p).Elem()
	for i := 0; i < vp.NumField(); i++ {
		v := vp.Field(i)
		switch v.Kind() {
		case reflect.Uint8: // byte
			size += ByteSize
		case reflect.Int8: // sbyte
			size += SByteSize
		case reflect.Int16: // short
			size += ShortSize
		case reflect.String: // string
			size += StringSize
		case reflect.Slice: // []byte
			size += BytesSize
		}
	}
	return
}
