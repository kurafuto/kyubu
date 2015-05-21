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
	"os"
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

var tmpCount = 0

func tmpVar() string {
	tmpCount++
	return fmt.Sprintf("tmp%d", tmpCount-1)
}

var (
	file        = flag.String("file", "packets.go", "The file containing packet definitions")
	output      = flag.String("output", "", "Output file. If empty, ${file}_proto.go")
	direction   = flag.String("direction", "Anomalous", "The packet direction: serverbound, clientbound, anomalous")
	state       = flag.String("state", "", "Hint of state that the packet is used in")
	varIntLen   = flag.Int("varIntLen", 32, "Max length of VarInts: 16, 32, 64")
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
		*output = strings.TrimSuffix(*file, ".go") + "_proto.go"
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
	// TODO: Make sure we need this
	imports["github.com/sysr-q/kyubu/packets"] = struct{}{}

	errWrap := func(x string, y ...interface{}) string {
		z := fmt.Sprintf(x, y...)
		// err has to be predefined.
		return fmt.Sprintf("if err = %s; err != nil {\nreturn\n}\n", z)
	}

	for _, p := range packets {
		t := "t" // func (t *T) ...

		// Id() byte
		fmt.Fprintf(&buf, "func (%s *%s) Id() byte {\nreturn %#.2x // %d\n}\n\n", t, p.name, p.id, p.id)

		// Encode(io.Writer) error
		fmt.Fprintf(&buf, "func (%s *%s) Encode(ww io.Writer) (err error) {\n", t, p.name)
		for _, field := range p.spec.Fields.List {
			var typeName string

			if ide, ok := field.Type.(*ast.Ident); ok {
				typeName = ide.Name
			} else if selx, ok := field.Type.(*ast.SelectorExpr); ok {
				typeName = selx.X.(*ast.Ident).Name
				typeName = typeName + "." + selx.Sel.Name
			}

			// TODO: panic?
			name := field.Names[0].Name
			fmt.Fprintf(&buf, "// Encoding: %s (%s)\n", name, typeName)

			switch typeName {
			case "bool":
				tmp := tmpVar() // byte value for bool
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
				fmt.Fprintf(&buf, "%s := make([]byte, binary.MaxVarintLen%d)\n", x, *varIntLen)
				fmt.Fprintf(&buf, "%s := []byte(%s.%s)\n", b, t, name)
				fmt.Fprintf(&buf, "%s := binary.PutVarint(%s, int64(len(%s)))\n", n, x, b)
				fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
				fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s)", Endianness, b))
			case "packets.VarInt", "packets.VarLong":
				x := tmpVar() // []byte for varint
				n := tmpVar() // num of varint bytes

				// TODO: Is making a new []byte every time efficient?
				fmt.Fprintf(&buf, "%s := make([]byte, binary.MaxVarintLen%d)\n", x, *varIntLen)
				fmt.Fprintf(&buf, "%s := binary.PutVarint(%s, int64(%s.%s))\n", n, x, t, name)
				fmt.Fprintf(&buf, errWrap("binary.Write(ww, %s, %s[:%s])", Endianness, x, n))
			case "packets.Chunk":
			case "packets.Metadata":
			case "packets.Slot":
			case "packets.ObjectData":
			case "packets.NBT":
			case "packets.Position":
			case "packets.Angle":
			case "packets.UUID":
			}

			fmt.Fprintf(&buf, "\n")
		}
		fmt.Fprintf(&buf, "return\n}\n\n")

		// Decode(io.Writer) error
		fmt.Fprintf(&buf, "func (%s *%s) Decode(rr io.Reader) (err error) {\n", t, p.name)
		for _, field := range p.spec.Fields.List {
			var typeName string

			if ide, ok := field.Type.(*ast.Ident); ok {
				typeName = ide.Name
			} else if selx, ok := field.Type.(*ast.SelectorExpr); ok {
				typeName = selx.X.(*ast.Ident).Name
				typeName = typeName + "." + selx.Sel.Name
			}

			// TODO: panic?
			name := field.Names[0].Name
			fmt.Fprintf(&buf, "// Decoding: %s (%s)\n", name, typeName)

			// TODO: ALL OF THIS
			switch typeName {
			case "bool":
				tmp := tmpVar() // byte value for bool
				fmt.Fprintf(&buf, "%s := make([]byte, 1)\n", tmp)
				fmt.Fprintf(&buf, errWrap("_, err = rr.Read(b)\n"))
				fmt.Fprintf(&buf, "if %s[0] == 0x01 {\n", tmp)
				fmt.Fprintf(&buf, "%s.%s = true\n", t, name)
				fmt.Fprintf(&buf, "} else {\n")
				fmt.Fprintf(&buf, "%s.%s = false\n", t, name)
				fmt.Fprintf(&buf, "}")
			case "int8", "uint8", "int16", "uint16", "int32", "int64", "float32", "float64":
				fmt.Fprintf(&buf, errWrap("binary.Read(rr, %s, %s.%s)", Endianness, t, name))
			case "string":
				imports["errors"] = struct{}{}

				n := tmpVar() // num of varint bytes
				b := tmpVar() // []byte to read into
				x := tmpVar() // num of bytes read

				fmt.Fprintf(&buf, "%s, err := binary.ReadVarint(rr)\n", n)
				fmt.Fprintf(&buf, "if err != nil {\nreturn\n}\n")
				fmt.Fprintf(&buf, "%s := make([]byte, %s)\n", b, n)
				fmt.Fprintf(&buf, "%s, err := rr.Read(%s)\n", x, b)
				fmt.Fprintf(&buf, "if err != nil {\nreturn\n} ")
				fmt.Fprintf(&buf, "else if int64(%s) != %s {\n", x, n)
				fmt.Fprintf(&buf, `return errors.New("didn't read enough bytes for string")`+"\n")
				fmt.Fprintf(&buf, "}\n")
				fmt.Fprintf(&buf, "%s.%s = string(%s)\n", t, name, b)
			case "packets.VarInt":
				n := tmpVar()
				fmt.Fprintf(&buf, "%s, err := binary.ReadVarint(rr)\n", n)
				fmt.Fprintf(&buf, "if err != nil {\nreturn\n}\n")
				fmt.Fprintf(&buf, "%s.%s = int32(%s)\n", t, name, n)
			case "packets.VarLong":
				n := tmpVar()
				fmt.Fprintf(&buf, "%s, err := binary.ReadVarint(rr)\n", n)
				fmt.Fprintf(&buf, "if err != nil {\nreturn\n}\n")
				fmt.Fprintf(&buf, "%s.%s = %s\n", t, name, n)
			case "packets.Chunk":
			case "packets.Metadata":
			case "packets.Slot":
			case "packets.ObjectData":
			case "packets.NBT":
			case "packets.Position":
			case "packets.Angle":
			case "packets.UUID":
			}

			fmt.Fprintf(&buf, "\n")
		}
		fmt.Fprintf(&buf, "return\n}\n\n")
	}

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
	fmt.Fprintf(&header, "func init() {\n")
	// TODO: for blah { Register() }
	fmt.Fprintf(&header, "}\n\n")

	// TODO: Write these to a file.
	fmt.Print(header.String())
	fmt.Print(buf.String())
}
