package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	input  = flag.String("input", "", "input file")
	output = flag.String("output-dir", ".", "output directory")
)

type Section struct {
	XMLName     xml.Name  `xml:"section"`
	Anchor      string    `xml:"anchor,attr"`
	PN          string    `xml:"pn,attr"`
	Subsections []Section `xml:"section"`
}

func processSections(sections map[string]Section, section Section) {
	sections[section.PN] = section
	for _, subsection := range section.Subsections {
		processSections(sections, subsection)
	}
}

type rfct struct {
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

	bs, err := os.ReadFile(*input)
	if err != nil {
		log.Printf("error reading the rfc XML file: %s", err)
		return
	}
	sections := make(map[string]Section)

	decoder := xml.NewDecoder(bytes.NewReader(bs))
	for t, err := decoder.Token(); err != io.EOF; t, err = decoder.Token() {
		switch el := t.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			default:
				continue
			case "section":
				var section Section
				if err := decoder.DecodeElement(&section, &el); err != nil {
					log.Printf("Error decoding: element=%s, error=%s", el, err)
				}
				processSections(sections, section)
			}
		}
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
	decoder = xml.NewDecoder(f)
	var keywords []string
	collection := make(map[string][]rfct)

	for t, err := decoder.Token(); err != io.EOF; t, err = decoder.Token() {
		switch el := t.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			default:
				continue
			case "t":
				var myT rfct
				if err := decoder.DecodeElement(&myT, &el); err != nil {
					log.Printf("Error decoding: element=%s, error=%s", el, err)
				}
				if myT.Section == "section-2.2-1" {
					keywords = myT.Req
					continue
				}

				// Section 16 has guideline for extending HTTP, which is
				// a more process rather than computational.
				if strings.HasPrefix(myT.Section, "section-16") {
					continue
				}

				if len(myT.Req) == 0 {
					continue
				}
				for _, req := range myT.Req {
					collection[req] = append(collection[req], myT)
				}
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
			continue
		}

		for _, r := range reqs {
			sectionTag := r.Section
			if n := strings.Count(r.Section, "-"); n > 1 {
				sectionTag = r.Section[:strings.LastIndex(r.Section, "-")]
			}

			prefixedMajorNumber := strings.Split(sectionTag, ".")[0]
			majorSectionNumber := strings.TrimPrefix(prefixedMajorNumber, "section-")
			majorName := sections[prefixedMajorNumber].Anchor

			rootDir := filepath.Join(*output, fmt.Sprintf("%s-%s", majorSectionNumber, strings.ReplaceAll(majorName, ".", "-")))
			sectionDir := filepath.Join(rootDir, fmt.Sprintf("%s-%s", r.Section, strings.ReplaceAll(sections[sectionTag].Anchor, ".", "-")))
			if err := os.MkdirAll(sectionDir, 0o755); err != nil {
				log.Printf("error creating output directory for section: %s", err)
				return
			}
			outfilepath := filepath.Join(sectionDir, req+".txt")
			out, err := os.OpenFile(outfilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
			if err != nil {
				log.Printf("error opening output file: %s", err)
				return
			}
			out.WriteString(fmt.Sprintf("- %s: %s\n", r.Section, r.Body))
			out.Sync()
			out.Close()
		}
	}
}
