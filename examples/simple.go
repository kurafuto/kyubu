import (
	"fmt"
	"github.com/sysr-q/kyubu"
	"github.com/sysr-q/kyubu/packets"
	"github.com/sysr-q/kyubu/auth"
)

func main() {
	// Create a kyubu.Settings from a direct connect URL:
	settings := kyubu.Direct("mc://mcc.example.com:25565/notch/a2a01deadbeef26314ecafebabe8a32c")
	settings.Trickle = 25

	// OR! Conventional authentication:
	auth := auth.NewClassiCube("notch", "hunter2")
	_, err := auth.Login() // (bool, error)
	if err != nil {
		panic(err)
	}

	// Grab a server from the serverlist:
	servers, err := auth.ServerList()
	if err != nil {
		panic(err)
	}

	settings := kyubu.Settings{
		Server: servers[0],
		Auth: auth,
		Trickle: 25,
		Debug: false,
	}

	//////////

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
		fmt.Printf("<-[%#.2x] %s: recv!
", packet.Id(), packetName)
	}
}
