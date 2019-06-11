//+build ignore

package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/vinothzomato/radius/dictionary"
	"github.com/vinothzomato/radius/dictionarygen"
)

func main() {
	resp, err := http.Get(`http://support.arubanetworks.com/ToolsResources/tabid/76/DMXModule/514/Command/Core_Download/Method/attachment/Default.aspx?EntryId=31244`)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("main.dictionary", body, 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove("main.dictionary")

	parser := dictionary.Parser{
		Opener: restrictedOpener{
			"main.dictionary": body,
		},
	}
	dict, err := parser.ParseFile("main.dictionary")
	if err != nil {
		log.Fatal(err)
	}

	gen := dictionarygen.Generator{
		Package: "aruba",
	}
	generated, err := gen.Generate(dict)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("generated.go", generated, 0644); err != nil {
		log.Fatal(err)
	}
}

type restrictedOpener map[string][]byte

func (r restrictedOpener) OpenFile(name string) (dictionary.File, error) {
	contents, ok := r[name]
	if !ok {
		return nil, errors.New("unknown file " + name)
	}
	return &restrictedFile{
		Reader:    bytes.NewReader(contents),
		NameValue: name,
	}, nil
}

type restrictedFile struct {
	io.Reader
	NameValue string
}

func (r *restrictedFile) Name() string {
	return r.NameValue
}

func (r *restrictedFile) Close() error {
	return nil
}
