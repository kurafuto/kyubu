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
	"go/parser"
	"go/token"
	"log"
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

var (
	file      = flag.String("file", "packets.go", "The file containing packet definitions")
	direction = flag.String("direction", "Anomalous", "The packet direction: serverbound, clientbound, anomalous")
	state     = flag.String("state", "", "Hint of state that the packet is used in")
)

func main() {
	flag.Parse()

	fset := token.NewFileSet()
	pf, err := parser.ParseFile(fset, *file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var (
		idPrefix   = "Packet ID: 0x"
		Endianness = "binary.BigEndian"
	)

	var (
		packets []packet
		imports = map[string]struct{}{}
		buf     bytes.Buffer
	)

	// Parse information about the packets we're gonna make parsers for.
	for _, d := range pf.Decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

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

	imports["io"] = struct{}{}
	imports["binary"] = struct{}{}

	errWrap := func(x string, y ...interface{}) string {
		z := fmt.Sprintf(x, y...)
		// err has to be predefined.
		return fmt.Sprintf("if err = %s; err != nil {\nreturn\n}\n", z)
	}

	for _, p := range packets {
		t := "t" // func (t *T) ...

		// Id() byte
		fmt.Fprintf(&buf, "func (%s *%s) Id() byte {\nreturn %d;\n}\n\n", t, p.name, p.id)

		// Encode(io.Writer) error
		fmt.Fprintf(&buf, "func (%s *%s) Encode(ww io.Writer) (err error) {\n", t, p.name)
		for _, field := range p.spec.Fields.List {
			var ident *ast.Ident

			if ide, ok := field.Type.(*ast.Ident); ok {
				ident = ide
			} else if selx, ok := field.Type.(*ast.SelectorExpr); ok {
				ident = selx.Sel
			}

			// TODO: panic?
			name := field.Names[0].Name
			fmt.Fprintf(&buf, "// Encoding: %s (%s)\n", name, ident.Name)

			switch ident.Name {
			case "bool":
				tmp := tmpVar()
				fmt.Fprintf(&buf, "%s := byte(0)\n", tmp)
				fmt.Fprintf(&buf, "if %s.%s {\n", t, name)
				fmt.Fprintf(&buf, "\t%s = byte(1)\n", tmp)
				fmt.Fprintf(&buf, "}\n")
				fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s)", Endianness, tmp))
			case "int8", "uint8", "int16", "uint16", "int32", "int64", "float32", "float64":
				fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s.%s)", Endianness, t, name))
			case "string":
				x := tmpVar() // []byte for varint
				b := tmpVar() // []byte(string)
				n := tmpVar() // num of varint bytes

				// TODO: Is making a new []byte every time efficient?
				fmt.Fprintf(&buf, "%s := make([]byte, binary.MaxVarintLen32)\n", x)
				fmt.Fprintf(&buf, "%s := []byte(*%s)\n", b, t)
				fmt.Fprintf(&buf, "%s := binary.PutVarint(%s, int64(len(%s)))\n", n, x, b)

				fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
				fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s)", Endianness, b))
			}

			fmt.Fprintf(&buf, "\n")
		}
		fmt.Fprintf(&buf, "return\n}\n\n")

		// Decode(io.Writer) error
	}

	fmt.Print(buf.String())
}
