Kyubu (キューブ)
================

[Minecraft Classic protocol](http://wiki.vg/Classic_Protocol) implementation in
Go, primarily to get the grips with networking stuff.

Currently, Kyubu is separated into a few logical parts (two of which are just
internals):

* `kyubu/client` - basic Minecraft Classic client implementation
* `kyubu/packets` - packet parser/serializer stuff
* `kyubu/chunk` - chunk/block related stuff
* `kyubu/auth`- Authentication, currently only supports ClassiCube

The client impl. is kind of messy right now, and goes through sweeping changes
every commit or two. If you look in the `examples` directory, you can see the
same script I use to test - that will always have a working example.

## Extending Kyubu

Should you wish to actually _use_ Kyubu (and I can't blame you, who wouldn't!),
you can extend the functionality of the parser easily by simply registering
additional packets with the parser. For example, here's how you'd implement the
__ExtInfo__ packet from the [Classic Protocol Extension](http://wiki.vg/Classic_Protocol_Extension#ExtInfo_Packet):

```go
package main

import (
	"fmt"
	"github.com/sysr-q/kyubu/packets"
)

var ExtInfoSize = packets.ByteSize + packets.StringSize + packets.ShortSize

type ExtInfo struct {
	PacketId       byte
	AppName        string
	ExtensionCount int16 // short
}

func (p *ExtInfo) Id() byte {
	return p.PacketId
}

func (p *ExtInfo) Size() int {
	return ExtInfoSize
}

func (p *ExtInfo) Bytes() []byte {
	// Use the lovely packets.PacketWrapper.
	raw := packets.NewPacketWrapper([]byte{})
	raw.WriteByte(p.PacketId)
	raw.WriteString(p.AppName)
	raw.WriteShort(p.ExtensionCount)
	return raw.Buffer.Bytes()
}

func ReadExtInfo(b []byte) (packets.Packet, error) {
	// This is kind of boilerplate, but it's just how it's done, sorry!
	p := ExtInfo{}
	raw := packets.NewPacketWrapper(b)
	if packetId, err := raw.ReadByte(); err != nil {
		return nil, err
	} else {
		p.PacketId = packetId
	}
	if appName, err := raw.ReadString(); err != nil {
		return nil, err
	} else {
		p.AppName = appName
	}
	if extCount, err := raw.ReadShort(); err != nil {
		return nil, err
	} else {
		p.ExtensionCount = extCount
	}
	return &p, nil
}

func NewExtInfo(appName string, extCount int16) (p *ExtInfo, err error) {
	if len(appName) > packets.StringSize {
		err = fmt.Errorf("cpe: Can't write string longer than %d", packets.StringSize)
		return
	}
	p = &ExtInfo{
		PacketId:       0x10,
		AppName:        appName,
		ExtensionCount: extCount,
	}
	return
}

func init() {
	packets.Register(&packets.PacketInfo{
		Id:   0x10,
		Read: ReadExtInfo,
		Size: ExtInfoSize,
		Type: packets.Both, // C<>S
		Name: "Ext Info [CPE]",
	})
}
```

The general idea is this:

1. Create a type that implements the `kyubu/packets.Packet` interface.
2. Make a function to create the packet from a `[]byte` and return the packet.
3. Make a function to create a fresh packet from scratch.
4. In your file's `init()` method, `Register()` the packet with a `packets.PacketInfo{}`.
5. Write some tests to ensure your packet works as intended (and preferably won't
   break existing packets!).

__Hints:__

* The `Register()` function returns a `(bool, error)` tuple, explaining whether
  your packet registration was successful, and if not, why it failed.
* Want to override existing packets? Set `kyubu/packets.AllowOverride` to `true`.
  You'll then be able to replace existing packets via the `Register()` function.

## Roadmap

Currently, I'm just going for semi-feature parity with the Classic client. Handling
users, messages, sending packets, teleporting around, modifying blocks, etc.
A nice addition that could happen is Notchian/~~ClassiCube~~ authentication.

Once that's all done, and at least remotely sane, I'd like to add some support for
[CPE](http://wiki.vg/Classic_Protocol_Extension).

## Tests

Currently, Kyubu ships with tests for `kyubu/packets` and `kyubu/chunk`. I've
gone for as much coverage as possible, but as per anything, 100% is just not
viable in this situation.

You can run them with the standard `go test` in `kyubu/packets` and `kyubu/chunk`.
If you make additions or revisions to the code, make sure the tests all pass, and
add any new tests for features/packets you've added.
