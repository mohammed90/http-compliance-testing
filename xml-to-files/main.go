package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var (
	input  = flag.String("input", "", "input file")
	output = flag.String("output-dir", ".", "output directory")
)

type RFCT struct {
	Req     []string `xml:"bcp14"`
	Body    string   `xml:",innerxml"`
	Section string   `xml:"pn,attr"`
}

// TODO: to be used to parse xrefs to smartly link references
type xref struct {
	XMLName        xml.Name `xml:"xref"`
	Target         string   `xml:"target,attr"`
	Section        string   `xml:"section,attr"`
	Format         string   `xml:"format,attr"`
	SectionFormat  string   `xml:"sectionFormat,attr"`
	Relative       string   `xml:"relative,attr"`
	DerivedContent string   `xml:"derivedContent,attr"`
}

func main() {
	flag.Parse()
	if *input == "" {
		log.Printf("input file is not specified")
		return
	}
	if *output == "" {
		log.Printf("output file is not specified")
		return
	}
	f, err := os.Open(*input)
	if err != nil {
		log.Printf("error opening the rfc XML file: %s", err)
		return
	}
	defer f.Close()
	if err := os.MkdirAll(*output, 0o755); err != nil {
		log.Printf("error creating output directory: %s", err)
		return
	}

	collection := make(map[string][]RFCT)
	decoder := xml.NewDecoder(f)
	var keywords []string
	for t, err := decoder.Token(); err != io.EOF; t, err = decoder.Token() {
		switch el := t.(type) {
		case xml.StartElement:
			if el.Name.Local != "t" {
				continue
			}
			var myT RFCT
			if err := decoder.DecodeElement(&myT, &el); err != nil {
				log.Printf("Error decoding: element=%s, error=%s", el, err)
			}
			if myT.Section == "section-2.2-1" {
				keywords = myT.Req
				continue
			}
			if len(myT.Req) == 0 {
				continue
			}
			for _, req := range myT.Req {
				collection[req] = append(collection[req], myT)
			}
		case xml.EndElement, xml.CharData,
			xml.Comment,
			xml.ProcInst,
			xml.Directive:
			continue
		}
	}
	sort.StringSlice(keywords).Sort()
	for _, req := range keywords {
		reqs, ok := collection[req]
		if !ok {
			log.Printf("No such req: %s", req)
			continue
		}
		outfilepath := filepath.Join(*output, req+".txt")
		out, err := os.OpenFile(outfilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			log.Printf("error opening output file: %s", err)
			return
		}
		for _, r := range reqs {
			out.WriteString(fmt.Sprintf("- %s: %s\n", r.Section, r.Body))
		}
		out.Sync()
		out.Close()
	}
}
