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
every commit or two. If you look in the `examples` directory, you can see the
same script I use to test - that will always have a working example.

## Roadmap

Currently, I'm just going for semi-feature parity with the Classic client. Handling
users, messages, sending packets, teleporting around, modifying blocks, etc.
A big change that needs to happen is Notchian/~~ClassiCube~~ (kind of done) authentication.

Once that's all done, and at least remotely sane, I'd like to add some support for
[CPE](http://wiki.vg/Classic_Protocol_Extension).

Another nice idea would be a POC for a simple Minecraft Classic _server_, getting
feature parity with the vanilla server shouldn't be hard. That's for later, though.
