package cpe

import (
	//"bytes"
	"github.com/sysr-q/kyubu/packets"
	"reflect"
	"strings"
	"testing"
)

type packetReader func([]byte) (packets.Packet, error)

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

	var p ExtPacket // god tier reflection cheats
	p = res[0].Interface().(ExtPacket)
	t.Logf("p.String() == %q", p.String())

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

	t.Logf("ReflectSize(p) = %d %#v", packets.ReflectSize(p), p)
	if packets.ReflectSize(p) != p.Size() {
		t.Fatalf("calculated ReflectSize %d != original size %d", packets.ReflectSize(p), p.Size())
	}

	if !reflect.DeepEqual(p, read) {
		t.Fatalf("original packet does not equal recreated packet: %#v %#v", p, read)
	}
}

func TestExtInfo(t *testing.T) {
	testPacketReencode(NewExtInfo, ReadExtInfo, t, false, "CPE Test", int16(1))
	testPacketReencode(NewExtInfo, ReadExtInfo, t, true, strings.Repeat("t", packets.StringSize+1), int16(1))
}

func TestExtEntry(t *testing.T) {
	testPacketReencode(NewExtEntry, ReadExtEntry, t, false, "CPE Test", int32(1))
	testPacketReencode(NewExtEntry, ReadExtEntry, t, true, strings.Repeat("t", packets.StringSize+1), int32(1))
}

func TestSetClickDistance(t *testing.T) {
	testPacketReencode(NewSetClickDistance, ReadSetClickDistance, t, false, int16(1))
}

func TestCustomBlockSupportLevel(t *testing.T) {
	testPacketReencode(NewCustomBlockSupportLevel, ReadCustomBlockSupportLevel, t, false, byte(1))
}

func TestHoldThis(t *testing.T) {
	testPacketReencode(NewHoldThis, ReadHoldThis, t, false, byte(49), byte(0))
}

func TestSetTextHotKey(t *testing.T) {
	testPacketReencode(NewSetTextHotKey, ReadSetTextHotKey, t, false, "CPE Test", "CPE Test!", int32(113), byte(0))
	testPacketReencode(NewSetTextHotKey, ReadSetTextHotKey, t, true, strings.Repeat("t", packets.StringSize+1), "CPE Test!", int32(113), byte(0))
	testPacketReencode(NewSetTextHotKey, ReadSetTextHotKey, t, true, "CPE Test", strings.Repeat("t", packets.StringSize+1), int32(113), byte(0))
}

func TestExtAddPlayerName(t *testing.T) {
	testPacketReencode(NewExtAddPlayerName, ReadExtAddPlayerName, t, false, int16(5), "Notch", "&c[Op]Notch", "Staff", byte(0))
	testPacketReencode(NewExtAddPlayerName, ReadExtAddPlayerName, t, true, int16(5), strings.Repeat("t", packets.StringSize+1), "&c[Op]Notch", "Staff", byte(0))
}

func TestExtAddEntity(t *testing.T) {
	testPacketReencode(NewExtAddEntity, ReadExtAddEntity, t, false, byte(5), "&cNotch", "Notch")
	testPacketReencode(NewExtAddEntity, ReadExtAddEntity, t, true, byte(5), strings.Repeat("t", packets.StringSize+1), "Notch")
}

func TestExtRemovePlayerName(t *testing.T) {
	testPacketReencode(NewExtRemovePlayerName, ReadExtRemovePlayerName, t, false, int16(5))
}

func TestEnvSetColor(t *testing.T) {
	testPacketReencode(NewEnvSetColor, ReadEnvSetColor, t, false, byte(1), int16(25), int16(128), int16(0))
}

func TestMakeSelection(t *testing.T) {
	v := []interface{}{byte(0), "SomeZone", int16(1), int16(2), int16(3), int16(5), int16(6), int16(7), int16(255), int16(34), int16(128), int16(255)}
	testPacketReencode(NewMakeSelection, ReadMakeSelection, t, false, v...)
	v = []interface{}{byte(0), strings.Repeat("t", packets.StringSize+1), int16(1), int16(2), int16(3), int16(5), int16(6), int16(7), int16(255), int16(34), int16(128), int16(255)}
	testPacketReencode(NewMakeSelection, ReadMakeSelection, t, true, v...)
}

func TestRemoveSelection(t *testing.T) {
	testPacketReencode(NewRemoveSelection, ReadRemoveSelection, t, false, byte(0))
}

func TestSetBlockPermission(t *testing.T) {
	testPacketReencode(NewSetBlockPermission, ReadSetBlockPermission, t, false, byte(8), byte(0), byte(1))
}

func TestChangeModel(t *testing.T) {
	testPacketReencode(NewChangeModel, ReadChangeModel, t, false, byte(5), "spider")
	testPacketReencode(NewChangeModel, ReadChangeModel, t, true, byte(5), strings.Repeat("t", packets.StringSize+1))
}

func TestEnvSetMapAppearance(t *testing.T) {
	testPacketReencode(NewEnvSetMapAppearance, ReadEnvSetMapAppearance, t, false, "http://example.com/terrain.png", byte(7), byte(8), int16(31))
	testPacketReencode(NewEnvSetMapAppearance, ReadEnvSetMapAppearance, t, true, strings.Repeat("t", packets.StringSize+1), byte(7), byte(8), int16(31))
}

func TestEnvSetWeatherType(t *testing.T) {
	testPacketReencode(NewEnvSetWeatherType, ReadEnvSetWeatherType, t, false, byte(1))
}

func TestHackControl(t *testing.T) {
	v := []interface{}{byte(0), byte(0), byte(0), byte(1), byte(1), int16(40)}
	testPacketReencode(NewHackControl, ReadHackControl, t, false, v...)
}
