Kyubu
=====

Minecraft protocol implementation, provides parsing, serializing and chunk
decompression.

Currently, Kyubu is separated into a few logical parts:

* `kyubu/packets` - packet parser/serializer stuff
	* `kyubu/classic` - classic packets
	* `kyubu/cpe` - [CPE](http://wiki.vg/Classic_Protocol_Extension) related packets.
* `kyubu/chunk` - chunk/block related stuff

## Enabling the CPE

If you'd like to make the Kyubu packet parser CPE-aware, it's easy! All you have
to do is `import "github.com/sysr-q/kyubu/cpe"`. That's it. Really!

`kyubu/cpe` will register all of the packets it implements, so you don't have to
do anything funky to get it to work with `kyubu/packets.Parser`.

## Extending Kyubu

Should you wish to actually _use_ Kyubu (and I can't blame you, who wouldn't!),
you can extend the functionality of the parser easily by simply registering
additional packets with the parser (`kyubu/packets.Register`).

For a working examples, see `kyubu/classic/identification.go` and perhaps
`kyubu/cpe/ext_info.go`. Those are bound to be working, they have tests after all!

The general idea is this:

1. Create a type that implements the `kyubu/packets.Packet` interface.
    * (`kyubu/cpe.ExtPacket` for CPE related packets)
    * `kyubu/packets` has a handy `ReflectBytes(Packet)` helper for adding `Bytes()`.
2. Make a function to create the packet from a `[]byte` and return the packet.
    * `kyubu/packets` has a handy `ReflectRead([]byte, *Packet)` helper for this.
3. Make a function to create a fresh packet from scratch.
4. In your file's `init()` method, `Register()` the packet with a `packets.PacketInfo{}`.
5. Write some tests to ensure your packet works as intended (and preferably won't
   break existing packets!).

## Roadmap

* `kyubu/cpe` - This is implemented at the time of writing (`05/10/2014, DD/MM/YYYY`),
  but I'd like to even stay on top of proposals, and keep them up to date as they
  get close to finalisation.
* Make sure the chunk decompression still works with modern Minecraft. I wrote
  it for Classic.

## Tests

Currently, Kyubu ships with tests for `kyubu/packets`, `kyubu/classic`,
`kyubu/chunk` and `kyubu/cpe`. I've gone for as much coverage as possible, but
as per anything, 100% is just not viable in this situation (except for
`kyubu/cpe`, whoo!).

You can run them with the standard `go test` in the above-mentioned packages.
If you make additions or revisions to the code, make sure the tests all pass, and
add any new tests for features/packets you've added.
