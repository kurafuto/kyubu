package main

import (
	"fmt"
	"github.com/sysr-q/kyubu/auth"
	"github.com/sysr-q/kyubu/client"
	"github.com/sysr-q/kyubu/packets"
)

func main() {
	//// Create a kyubu.Settings from a direct connect URL:
	settings, err := client.Direct("mc://mcc.example.com:25565/notch/a2a01deadbeef26314ecafebabe8a32c")
	if err != nil {
		panic(err)
	}
	settings.Trickle = 25

	//// OR! Conventional authentication:
	auth := auth.NewClassiCube("notch", "hunter2")
	authed, err := auth.Login() // (bool, error)
	if err != nil {
		panic(err)
	} else if !authed {
		panic("Authentication failed!")
	}

	// Then you grab a server from the serverlist:
	servers, err := auth.ServerList()
	if err != nil {
		panic(err)
	}

	settings := kyubu.Settings{
		Server:  servers[0],
		Auth:    auth,
		Trickle: 25,
		Debug:   false,
	}

	//// Create a new client with the settings

	k, err := client.New(settings)
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

		packetName := packets.Packets[packet.Id()].Name
		fmt.Printf("<-[%#.2x] %s: recv!", packet.Id(), packetName)
	}
}
