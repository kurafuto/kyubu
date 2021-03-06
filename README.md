Kyubu
=====

Minecraft protocol implementation, provides parsing, serializing and chunk
decompression.

Currently, Kyubu is separated into a few logical parts:

* `kyubu/cmd` - Where `protocol_generator` et al live
* `kyubu/format` - the Chat format, lifted from [thinkofdeath/steven](https://github.com/thinkofdeath/steven)
* `kyubu/modern/minimal` - Minimal modern (1.8.3/47) packet definitions
	Enough for [Kurafuto](https://github.com/kurafuto/kurafuto) to function, but not everything
* `kyubu/packets` - packet parser/serializer stuff

## Extending Kyubu

Should you wish to actually _use_ Kyubu (and I can't blame you, who wouldn't!),
you can extend the functionality of the parser easily by simply registering
additional packets with the parser (`kyubu/packets.Register`).

For a working example, see `kyubu/cmd/protocol_generatortest/test_proto.go`.

The general idea is this:

1. Create a type that implements the `kyubu/packets.Packet` interface.
2. Add a `go generate` line to your file:
	* `//go:generate protocol_generator -file=$GOFILE -direction=x -state=x -package=x`
	* `go get github.com/kurafuto/kurafuto/cmd/protocol_generator` to install it.
3. Run `go generate` and ensure the `x_proto.go` file looks like you expect.
4. Write some tests to ensure your packet works as intended (and preferably won't
	break existing packets!).

## protocol_generator

To save on time, rather than the old method of reflection, Kyubu now ships with
a custom `go generate`-compatible command that makes your packet parser/serializer
cruft for you. Look at `cmd/protocol_generator/test/test.go` for an example.

NOTE: It's heavily based off of the protocol_generator from
[thinkofdeath/steven](https://github.com/thinkofdeath/steven).

## Tests

Tests are shipped for `kyubu/packets`, but as Kyubu is going through breaking
changes to work with Modern, it may not pass.

You can run them with the standard `go test` in the above-mentioned packages.
If you make additions or revisions to the code, make sure the tests all pass, and
add any new tests for features/packets you've added.
