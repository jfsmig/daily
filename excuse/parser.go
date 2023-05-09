// Copyright (C) 2023 Jean-Francois Smigielski
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package excuse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type frameType uint32

const (
	frameTypeConcat frameType = iota
	frameTypeChoice
)

type frame interface {
	// The sequence is over. Sanity checks must be performed and the final Generator returned.
	finish() Generator
	//
	parse(in io.RuneReader) error
}

var errRune = errors.New("unexpected token")
var errClosed = errors.New("sequence end")

type frameChoice struct {
	parts []Generator
}

func newFrameChoice() frame { return &frameChoice{parts: make([]Generator, 0)} }

func (f *frameChoice) finish() Generator {
	switch {
	case len(f.parts) < 1:
		return nil
	case len(f.parts) == 1:
		return f.parts[0]
	default:
		return NewChoice(f.parts...)
	}
}

func (f *frameChoice) parse(in io.RuneReader) error {
	for {
		next := newFrameConcat()
		err := next.parse(in)
		if g := next.finish(); g != nil {
			f.parts = append(f.parts, g)
		}
		switch err {
		case io.EOF, errClosed:
			return err
		case nil:
			continue
		default:
			return err
		}
	}
}

type frameConcat struct {
	str   strings.Builder
	parts []Generator
}

func newFrameConcat() frame { return &frameConcat{parts: make([]Generator, 0)} }

func (f *frameConcat) finishInner() {
	if f.str.Len() > 0 {
		f.parts = append(f.parts, NewTerm(f.str.String()))
		f.str.Reset()
	}
}

func (f *frameConcat) finish() Generator {
	f.finishInner()
	switch {
	case len(f.parts) < 1:
		return nil
	case len(f.parts) == 1:
		return f.parts[0]
	default:
		return NewSequence(f.parts...)
	}
}

func (f *frameConcat) parse(in io.RuneReader) error {
	for {
		c, _, err := in.ReadRune()
		if err != nil {
			return err
		}

		switch c {
		case '<':
			f.finishInner()
			sub := newFrameChoice()
			err = sub.parse(in)
			if g := sub.finish(); g != nil {
				f.parts = append(f.parts, g)
			}
			if err == io.EOF {
				return err
			} else {
				continue
			}
		case '>': // finishes the upper choice (item + sequence)
			f.finishInner()
			return errClosed
		case '|': // finishes the current item in the upper choice
			f.finishInner()
			return nil
		default:
			f.str.WriteRune(c)
		}
	}
}

// "plop" -> Term("plop")
// "plip plop <a,e,i,o,u> pouet" ->
// Concat(Term("plip plop"), Choice(Term("a"), Term("e"), Term("i"), Term("o"), Term("u")), Term("pouet"))
func ParseExpression(encoded string) (Generator, error) {
	in := strings.NewReader(encoded)
	top := newFrameConcat()
	if err := top.parse(in); err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	if in.Len() > 0 {
		return nil, errRune
	} else if g := top.finish(); g != nil {
		return g, nil
	} else {
		// Ugly corner case: in intermediate computations we use that nil return to discard
		// the related generator parser. But the top-level parsing requires an output even
		// for empty (collapsed) pattern.
		return NewEmpty(), nil
	}
}

func ParseStream(r io.Reader) (Generator, error) {
	in := bufio.NewScanner(r)
	out := &Choice{items: make([]Generator, 0)}
	for lineNum := 1; in.Scan(); lineNum++ {
		line := strings.Trim(in.Text(), "\n\t\r ")
		if strings.HasPrefix(line, "#") {
			continue
		}
		if item, err := ParseExpression(line); err != nil {
			return nil, fmt.Errorf("Error at line %d: %w", lineNum, err)
		} else {
			out.items = append(out.items, item)
		}
	}
	if err := in.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func ParseStreamString(encoded string) (Generator, error) {
	return ParseStream(strings.NewReader(encoded))
}

func ParseStreamBytes(encoded []byte) (Generator, error) {
	return ParseStream(bytes.NewReader(encoded))
}
