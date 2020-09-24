/*
 * \brief basic types used in pdf
 * \author chenxun
 */

package pdf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type BasicType interface {
	PDFBasicType() string
	String() string
}

type Boolean bool

func (b Boolean) PDFBasicType() string {
	return "Boolean"
}

func (b Boolean) String() string {
	return fmt.Sprintf("%v", bool(b))
}

type String string

func (s String) PDFBasicType() string {
	return "String"
}

func (s String) String() string {
	return string(s)
}

func NewString(text string) String {
	text = strings.Replace(text, "\\(", "(", -1)
	text = strings.Replace(text, "\\)", ")", -1)
	text = strings.Replace(text, "\\\n", "", -1)

	return String(text)
}

type Integer int64

func (i Integer) PDFBasicType() string {
	return "Integer"
}

func (i Integer) String() string {
	return fmt.Sprintf("%v", int64(i))
}

type Real float64

func (r Real) PDFBasicType() string {
	return "Real"
}

func (r Real) String() string {
	return fmt.Sprintf("%v", float64(r))
}

type Hexadecimal string

func (h Hexadecimal) PDFBasicType() string {
	return "Hexadecimal"
}

func (h Hexadecimal) String() string {
	return string(h)
}

type Name string

func (n Name) PDFBasicType() string {
	return "Name"
}

func (n Name) String() string {
	return string(n)
}

type Array []BasicType

func (a Array) PDFBasicType() string {
	return "Array"
}

func (a Array) String() string {
	bytes := []byte{'['}
	for _, item := range a {
		bytes = append(bytes, []byte(fmt.Sprintf("%v ", item))...)
	}
	bytes[len(bytes)-1] = ']'

	return string(bytes)
}

func NewArray(text string) Array {

	items, err := SplitBasicTypes(text)
	if err != nil {
		return nil
	}

	return items
}

type Dictionary map[Name]BasicType

func (d Dictionary) PDFBasicType() string {
	return "Dictionary"
}

func (d Dictionary) String() string {
	bytes := []byte{'<', '<'}
	for key, value := range d {
		bytes = append(bytes, []byte(fmt.Sprintf("/%v %v", key, value))...)
	}
	bytes[len(bytes)-1] = '>'
	bytes = append(bytes, '>')

	return string(bytes)
}

func NewDictionary(text string) Dictionary {

	reg := regexp.MustCompile("/[^/]+")
	items := reg.FindAllString(text, -1)
	d := Dictionary{}

	for _, item := range items {
		bts, err := SplitBasicTypes(item)
		if err != nil || len(bts) != 2 {
			return nil
		}
		name := bts[0].(Name)
		d[name] = bts[1]
	}

	return d
}

// TODO(chenxun): more accurate
func NewBasicType(text string) BasicType {
	text = strings.TrimSpace(text)

	if text == "true" {
		return Boolean(true)
	} else if text == "false" {
		return Boolean(false)
	} else if val, err := strconv.Atoi(text); err == nil {
		return Integer(val)
	} else if val, err := strconv.ParseFloat(text, 64); err == nil {
		return Real(val)
	} else if len(text) > 0 && text[0] == '<' && text[len(text)-1] == '>' {
		return Hexadecimal(text[1 : len(text)-1])
	} else if len(text) > 0 && text[0] == '/' {
		return Name(text[1:])
	} else if len(text) > 0 && text[0] == '(' && text[len(text)-1] == ')' {
		return NewString(text[1 : len(text)-1])
	} else if len(text) > 0 && text[0] == '[' && text[len(text)-1] == ']' {
		return NewArray(text[1 : len(text)-1])
	} else if len(text) > 1 && text[0:2] == "<<" && text[len(text)-2:len(text)] == ">>" {
		return NewDictionary(text[2 : len(text)-2])
	}

	return nil
}

func SplitBasicTypes(text string) ([]BasicType, error) {
	text = strings.TrimSpace(text)

	cur := 0
	res := []BasicType{}

	for i := 0; i < len(text); i++ {
		// Skip spaces
		for i < len(text) && text[i] == ' ' {
			i++
		}

		cur = i
		switch text[i] {
		case '<':
			for i < len(text) && text[i] != '>' {
				i++
			}
		case '(':
			for i < len(text) && !(text[i] == ')' && text[i-1] != '\\') {
				i++
			}
		case '[':
			for i < len(text) && text[i] != ']' {
				i++
			}
		default:
			for i < len(text) && text[i] != ' ' {
				i++
			}
			i--
		}

		if bt := NewBasicType(text[cur : i+1]); bt == nil {
			return nil, fmt.Errorf("parse basic types failed")
		} else {
			res = append(res, bt)
		}
	}

	return res, nil
}
