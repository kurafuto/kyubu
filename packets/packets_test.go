package packets

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

type packetReader func([]byte) (Packet, error)

func testPacketReencode(creator interface{}, reader packetReader, t *testing.T, expectErr bool, a ...interface{}) {
	// Forgive me, reflection overlords, for my misuse of your package.
	vcreator := reflect.ValueOf(creator)
	args := []reflect.Value{}
	for _, val := range a {
		args = append(args, reflect.ValueOf(val))
	}
	res := vcreator.Call(args)
	if len(res) != 2 {
		t.Errorf("expected 2 return values, got %d", len(res))
	}

	if expectErr && res[1].Interface() == nil {
		t.Fatal("expected err, got nil")
	} else if expectErr {
		// We're done here.
		return
	} else if !expectErr && res[1].Interface() != nil {
		var e error
		e = res[1].Interface().(error)
		t.Fatal("instanciate error:", e)
	}

	var p Packet // god tier reflection cheats
	p = res[0].Interface().(Packet)

	read, err := reader(p.Bytes())
	if expectErr && err == nil {
		t.Fatal("expected err, got nil")
	} else if !expectErr && err != nil {
		t.Fatal("error recreating packet:", err)
	}

	// Hopefully these two /never/ happen, though. Probably means tests are borke.
	if p.Id() != read.Id() {
		t.Fatalf("original packet id %#.2x != recreated id %#.2x", p.Id(), read.Id())
	}
	if p.Size() != read.Size() {
		t.Fatalf("original packet size %d != recreated size %d", p.Size(), read.Size())
	}
	t.Logf("ReflectSize(p) = %d %#v", ReflectSize(p), p)
	if ReflectSize(p) != p.Size() {
		t.Fatalf("calculated ReflectSize %d != original size %d", ReflectSize(p), p.Size())
	}

	if !reflect.DeepEqual(p, read) {
		t.Fatalf("original packet does not equal recreated packet: %#v %#v", p, read)
	}
}

func TestIdentification(t *testing.T) {
	testPacketReencode(NewIdentification, ReadIdentification, t, false, "test", "test")
	testPacketReencode(NewIdentification, ReadIdentification, t, true, strings.Repeat("t", StringSize+1), "test")
}

func TestPing(t *testing.T) {
	testPacketReencode(NewPing, ReadPing, t, false)
}

func TestLevelInitialize(t *testing.T) {
	testPacketReencode(NewLevelInitialize, ReadLevelInitialize, t, false)
}

func TestLevelDataChunk(t *testing.T) {
	data := bytes.Repeat([]byte{0x6e, 0x6f}, BytesSize/2)
	testPacketReencode(NewLevelDataChunk, ReadLevelDataChunk, t, false, data, byte(25))
	data = bytes.Repeat([]byte{0x01}, BytesSize+1)
	testPacketReencode(NewLevelDataChunk, ReadLevelDataChunk, t, true, data, byte(25))
}

func TestLevelFinalize(t *testing.T) {
	testPacketReencode(NewLevelFinalize, ReadLevelFinalize, t, false, int16(256), int16(64), int16(256))
}

func TestSetBlock5(t *testing.T) {
	testPacketReencode(NewSetBlock5, ReadSetBlock5, t, false, int16(0), int16(0), int16(0), byte(0), byte(0))
	testPacketReencode(NewSetBlock5, ReadSetBlock5, t, true, int16(0), int16(0), int16(0), byte(64), byte(0))
}

func TestSetBlock6(t *testing.T) {
	testPacketReencode(NewSetBlock6, ReadSetBlock6, t, false, int16(0), int16(0), int16(0), byte(0))
}

func TestSpawnPlayer(t *testing.T) {
	v := []interface{}{int8(1), "Notch", int16(64), int16(32), int16(64), byte(0), byte(0)}
	testPacketReencode(NewSpawnPlayer, ReadSpawnPlayer, t, false, v...)
	v = []interface{}{int8(1), strings.Repeat("t", StringSize+1), int16(64), int16(32), int16(64), byte(0), byte(0)}
	testPacketReencode(NewSpawnPlayer, ReadSpawnPlayer, t, true, v...)
}

func TestPositionOrientation(t *testing.T) {
	v := []interface{}{int8(1), int16(64), int16(32), int16(64), byte(0), byte(0)}
	testPacketReencode(NewPositionOrientation, ReadPositionOrientation, t, false, v...)
}

func TestPositionOrientationUpdate(t *testing.T) {
	v := []interface{}{int8(1), int8(4), int8(2), int8(4), byte(0), byte(0)}
	testPacketReencode(NewPositionOrientationUpdate, ReadPositionOrientationUpdate, t, false, v...)
}

func TestPositionUpdate(t *testing.T) {
	v := []interface{}{int8(1), int8(4), int8(2), int8(4)}
	testPacketReencode(NewPositionUpdate, ReadPositionUpdate, t, false, v...)
}

func TestOrientationUpdate(t *testing.T) {
	v := []interface{}{int8(1), byte(0), byte(0)}
	testPacketReencode(NewOrientationUpdate, ReadOrientationUpdate, t, false, v...)
}

func TestDespawnPlayer(t *testing.T) {
	testPacketReencode(NewDespawnPlayer, ReadDespawnPlayer, t, false, int8(1))
}

func TestMessage(t *testing.T) {
	testPacketReencode(NewMessage, ReadMessage, t, false, int8(1), "Hello, world!")
	testPacketReencode(NewMessage, ReadMessage, t, true, int8(1), strings.Repeat("t", StringSize+1))
}

func TestDisconnectPlayer(t *testing.T) {
	testPacketReencode(NewDisconnectPlayer, ReadDisconnectPlayer, t, false, "Disconnected!")
}

func TestUpdateUserType(t *testing.T) {
	testPacketReencode(NewUpdateUserType, ReadUpdateUserType, t, false, byte(0x64))
}
