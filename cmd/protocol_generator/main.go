// protocol_generator is a tool that we use with `go generate` to turn the
// useless structs for protocol packets into the reader/writer functions.
//
// Heavily based off of the protocol_generator in thinkofdeath's project steven,
// but designed to work with Kyubu's new packet registration system.
//   https://github.com/thinkofdeath/steven/blob/master/cmd/protocol_builder/builder.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type packet struct {
	// The ID of the packet we're gonna register.
	id int
	// The name of the packet we're working with.
	// Used to register and generate parse/serialize code.
	name string
	//
	spec *ast.StructType
}

func errWrap(x string, y ...interface{}) string {
	z := fmt.Sprintf(x, y...)
	// err has to be predefined.
	return fmt.Sprintf("if err = %s; err != nil {\nreturn\n}\n", z)
}

var (
	file        = flag.String("file", "packets.go", "The file containing packet definitions")
	output      = flag.String("output", "", "Output file. If empty, ${file}_proto.go")
	direction   = flag.String("direction", "anomalous", "The packet direction: serverbound, clientbound, anomalous")
	state       = flag.String("state", "", "State that the packet is used in")
	packageName = flag.String("package", "", "The name of the package the output will end up in")
)

var (
	idPrefix   = "Packet ID: 0x"
	Endianness = "binary.BigEndian"
)

var (
	packets []packet
	imports = map[string]struct{}{}
	buf     bytes.Buffer
)

func main() {
	flag.Parse()

	if *output == "" {
		*output = (*file)[:len(*file)-len(filepath.Ext(*file))] + "_proto.go"
	}

	fset := token.NewFileSet()
	pf, err := parser.ParseFile(fset, *file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Parse information about the packets we're gonna make parsers for.
	for _, d := range pf.Decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		// Make sure we have 1 type declaration
		if decl.Tok != token.TYPE || len(decl.Specs) != 1 {
			continue
		}

		spec, ok := decl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		// We only want: type X struct{}
		if _, ok := spec.Type.(*ast.StructType); !ok {
			continue
		}

		// Find the packet ID
		doc := decl.Doc.Text()
		pos := strings.Index(doc, idPrefix)

		if pos == -1 {
			log.Printf("Couldn't find packet ID for type: %s\n", spec.Name.Name)
			continue
		}

		idString := strings.TrimSpace(doc[pos+len(idPrefix):])

		id, err := strconv.ParseInt(idString, 16, 32)
		if err != nil {
			log.Printf("Error parsing id %q for type: %s\n", idString, spec.Name.Name)
			continue
		}

		packets = append(packets, packet{
			id:   int(id),
			name: spec.Name.Name,
			spec: spec.Type.(*ast.StructType),
		})
	}

	for _, p := range packets {
		t := "t" // func (t *T) ...

		en := Encoder{p: p, t: t, tmpPrefix: "tmp"}
		de := Decoder{p: p, t: t, tmpPrefix: "tmp"}

		// Id() byte
		fmt.Fprintf(&buf, "func (%s *%s) Id() byte {\nreturn %#.2x // %d\n}\n\n", t, p.name, p.id, p.id)

		// Encode(io.Writer) error
		fmt.Fprintf(&buf, "func (%s *%s) Encode(ww io.Writer) (err error) {\n", t, p.name)
		en.Write()
		fmt.Fprintf(&buf, "return\n}\n\n")

		// Decode(io.Writer) error
		fmt.Fprintf(&buf, "func (%s *%s) Decode(rr io.Reader) (err error) {\n", t, p.name)
		de.Write()
		fmt.Fprintf(&buf, "return\n}\n\n")
	}

	// TODO: Make sure we need all these.
	imports["encoding/binary"] = struct{}{}
	imports["github.com/sysr-q/kyubu/packets"] = struct{}{}
	imports["io"] = struct{}{}

	var header bytes.Buffer

	// General information
	//header.WriteString("// Generated by protocol_generator\n")
	header.WriteString("// ARR, HERE BE DRAGONS! DO NOT EDIT\n") // Pirates are hip, right?
	header.WriteString("// " + strings.Join(os.Args, " "))
	fmt.Fprintf(&header, "\n\npackage %s\n\n", *packageName)

	// Imports
	fmt.Fprint(&header, "import (\n")
	for imp := range imports {
		fmt.Fprintf(&header, "%q\n", imp)
	}
	fmt.Fprint(&header, ")\n\n")

	// init() -- register all the packets
	var states = map[string]string{
		"handshake": "packets.Handshake",
		"status":    "packets.Status",
		"login":     "packets.Login",
		"play":      "packets.Play",
	}
	var directions = map[string]string{
		"serverbound": "packets.ServerBound",
		"clientbound": "packets.ClientBound",
		"anomalous":   "packets.Anomalous",
	}

	fmt.Fprintf(&header, "func init() {\n")
	for _, p := range packets {
		fs := "packets.Register(%s, %s, %#.2x, func() packets.Packet { return &%s{} })\n"
		fmt.Fprintf(&header, fs, states[*state], directions[*direction], p.id, p.name)
	}
	fmt.Fprintf(&header, "}\n\n")

	buf.WriteTo(&header)

	b, err := format.Source(header.Bytes())
	if err != nil {
		log.Println(header.String())
		log.Fatalf("format error: %s", err)
	}

	o, err := os.Create(*output)
	if err != nil {
		log.Fatalln(err)
	}
	defer o.Close()
	o.Write(b)
}
