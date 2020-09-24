/*
 * \brief     objects in pdf
 * \author    chenxun
 */

package pdf

import (
	"fmt"
	"regexp"
	"strconv"
)

// Header declares the header of pdf files
type Header struct {
	Text string // exmaple: %PDF-1.5
}

func (h *Header) Version() (string, error) {
	reg := regexp.MustCompile("%PDF-1\\.[0-8]")

	results := reg.FindAllString(h.Text, -1)

	if len(results) != 1 {
		return "", fmt.Errorf("unknown pdf version")
	}

	return results[0], nil
}

// Body declares the body of pdf files
type Body struct {
	//TODO(chenxun): fix it
	//Objects []Object
}

// XrefTable declares the cross-reference table of the pdf files
type XrefTable struct {
	Start int64
	Count int64
	Items []XrefItem
}

type XrefItem struct {
	Offset  uint64
	Version uint32
	Valid   bool
}

type Trailer struct {
	Attributes Dictionary
	XrefOffest uint64
}

func NewTrailer(text string) *Trailer {
	reg := regexp.MustCompile("^trailer.*%%EOF$")

	if !reg.MatchString(text) {
		return nil
	}

	trailer := &Trailer{
		Attributes: Dictionary{},
	}

	reg = regexp.MustCompile("(?<=startxref\\n)\\d+")
	offset, err := strconv.Atoi(reg.FindString(text))
	if err != nil {
		return nil
	}
	trailer.XrefOffest = uint64(offset)

	//

	return trailer
}
