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
	"context"
	"errors"
	"io"
	"math/rand"
	"strings"
)

type Env struct {
	Prng *rand.Rand
}

func NewEnv(seed int64) *Env { return &Env{Prng: rand.New(rand.NewSource(seed))} }

type Generator interface {
	// Expand populates the w output with the runes of a random string
	Expand(ctx context.Context, w io.StringWriter, env *Env) error
	// Encode dumps a representation of the Generator that can be parsed in another Generator producing the same output
	Encode(out *strings.Builder)
	// Json dumps a JSON representation of the current Generator
	Json(out *strings.Builder)
	// MaxLength returns the longets possible string the current Generator can produce
	MaxLength() int
	// Count returns the number of possible choices at this level of the tree
	Count() int
}

type Concat struct {
	items []Generator
	count int
}

// Expand sequentially forward the call to each of its component, in order.
func (t *Concat) Expand(ctx context.Context, w io.StringWriter, env *Env) error {
	for idx, _ := range t.items {
		if err := t.items[idx].Expand(ctx, w, env); err != nil {
			return err
		}
	}
	return nil
}

func (t *Concat) Encode(out *strings.Builder) {
	if len(t.items) > 0 {
		t.items[0].Encode(out)
		if len(t.items) > 1 {
			for _, t := range t.items[1:] {
				t.Encode(out)
			}
		}
	}
}

func (t *Concat) Json(out *strings.Builder) {
	out.WriteString("[")
	if len(t.items) > 0 {
		t.items[0].Json(out)
		if len(t.items) > 1 {
			for _, t := range t.items[1:] {
				out.WriteRune(',')
				t.Json(out)
			}
		}
	}
	out.WriteString("]")
}

func (t *Concat) MaxLength() int {
	var total int
	for i, _ := range t.items {
		total += t.items[i].MaxLength()
	}
	return total
}

func (t *Concat) Count() int { return t.count }

func NewSequence(items ...Generator) Generator {
	out := &Concat{items: items}
	// Recompute the count of possibilities here
	out.count = 1
	for i, _ := range out.items {
		out.count *= out.items[i].Count()
	}
	return out
}

type Choice struct {
	items []Generator
	count int
}

// Expand randomly picks a component and forwards the call to it.
// The function respects a weight set to the number of possibilities behind each component of the choice.
func (t *Choice) Expand(ctx context.Context, w io.StringWriter, env *Env) error {
	needle := env.Prng.Intn(t.Count())
	cursor := 0
	for _, x := range t.items {
		cursor += x.Count()
		if cursor >= needle {
			return x.Expand(ctx, w, env)
		}
	}
	return errors.New("bug")
}

func (t *Choice) Encode(out *strings.Builder) {
	out.WriteRune('<')
	if len(t.items) > 0 {
		t.items[0].Encode(out)
		if len(t.items) > 1 {
			for _, t := range t.items[1:] {
				out.WriteRune('|')
				t.Encode(out)
			}
		}
	}
	out.WriteRune('>')
}

func (t *Choice) Json(out *strings.Builder) {
	out.WriteString("[\"?\"")
	for _, t := range t.items {
		out.WriteRune(',')
		t.Json(out)
	}
	out.WriteString("]")
}

func (t *Choice) MaxLength() int {
	var max int
	for i, _ := range t.items {
		if l := t.items[i].MaxLength(); l > max {
			max = l
		}
	}
	return max
}

func (t *Choice) Count() int { return t.count }

func NewChoice(items ...Generator) Generator {
	out := &Choice{items: items}
	out.count = 0
	for i, _ := range out.items {
		out.count += out.items[i].Count()
	}
	return out
}

type Term string

// Expand generates its input string
func (t *Term) Expand(ctx context.Context, w io.StringWriter, env *Env) error {
	_, err := w.WriteString(string(*t))
	return err
}

func (t *Term) Encode(out *strings.Builder) { out.WriteString(string(*t)) }

func (t *Term) Json(out *strings.Builder) {
	out.WriteRune('"')
	out.WriteString(string(*t))
	out.WriteRune('"')
}

func (t *Term) MaxLength() int { return len(*t) }

func (t *Term) Count() int { return 1 }

func NewTerm(s string) Generator { t := Term(s); return &t }

type Empty struct{}

// Expand generates nothing
func (t *Empty) Expand(_ context.Context, _ io.StringWriter, _ *Env) error { return nil }

func (t *Empty) Encode(out *strings.Builder) {}

func (t *Empty) Json(out *strings.Builder) { out.WriteString("nil") }

func (t *Empty) MaxLength() int { return 0 }

func (t *Empty) Count() int { return 1 }

func NewEmpty() Generator { return &Empty{} }

func EncodeGenerator(g Generator) string {
	str := strings.Builder{}
	g.Encode(&str)
	return str.String()
}

func JsonGenerator(g Generator) string {
	str := strings.Builder{}
	g.Json(&str)
	return str.String()
}
