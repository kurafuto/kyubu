Kyubu (キューブ)
================

[Minecraft Classic protocol](http://wiki.vg/Classic_Protocol) implementation in
Go, primarily to get the grips with networking stuff.

Currently, Kyubu is separated into a few logical parts (two of which are just
internals):

* `kyubu` - basic Minecraft Classic client implementation
* `kyubu/packets` - packet parser/serializer stuff
* `kyubu/chunk` - chunk/block related stuff
* `kyubu/auth`- ClassiCube authentication

The client impl. is kind of messy right now, and goes through sweeping changes
every commit or two. Right now, you can use it like:

```go
import (
	"fmt"
	"github.com/sysr-q/kyubu"
	"github.com/sysr-q/kyubu/packets"
	"github.com/sysr-q/kyubu/auth"
)

func main() {
	// NOTE! ClassiCube auth is not 100% implemented, so this will fail to
	// connect to verify-names=true servers.
	settings := kyubu.Settings{
		Server: auth.Server{Address: "mcc.example.com", Port: 25565, MpPass: "validationkey"},
		Auth: auth.NewClassiCube("notch", "password"),
		Trickle: 25,
		Debug: false,
	}
	k, err := kyubu.New(settings)
	if err != nil {
		panic(err)
	}

	saidHello := false
	for {
		packet := <-k.In
		if packet == nil {
			// Server disconnect, etc.
			break
		}

		if !saidHello && k.LoggedIn {
			// You should check the error for new packets, but this is just an example.
			mesg, _ := packets.NewMessage(127, "Hello, world!")
			k.Out <- mesg
			saidHello = true
		}

		packetName := packets.Packets[packet.Id()].Ident
		fmt.Printf("<-[%#.2x] %s: recv!\n", packet.Id(), packetName)
	}
}
```

## Roadmap

Currently, I'm just going for semi-feature parity with the Classic client. Handling
users, messages, sending packets, teleporting around, modifying blocks, etc.
A big change that needs to happen is Notchian/~~ClassiCube~~ (kind of done) authentication.

Once that's all done, and at least remotely sane, I'd like to add some support for
[CPE](http://wiki.vg/Classic_Protocol_Extension).

Another nice idea would be a POC for a simple Minecraft Classic _server_, getting
feature parity with the vanilla server shouldn't be hard. That's for later, though.
