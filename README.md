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
