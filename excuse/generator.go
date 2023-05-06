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
	"io"
	"math/rand"
	"strings"
)

type Env struct {
	Prng *rand.Rand
}

func NewEnv(seed int64) *Env { return &Env{Prng: rand.New(rand.NewSource(seed))} }

type Generator interface {
	Expand(ctx context.Context, w io.StringWriter, env *Env) error
	Encode(out *strings.Builder)
	Json(out *strings.Builder)
	MaxLength() int
}

type Concat struct {
	items []Generator
}

// Expand sequentially forward the call to each of its component, in order.
func (t *Concat) Expand(ctx context.Context, w io.StringWriter, env *Env) error {
	for idx, _ := range t.items {
		if err := t.items[idx].Expand(ctx, w, env); err != nil {
			return err
		}
		_, _ = w.WriteString(" ")
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
	out.WriteString("{\"type\":\"sequence\",\"items\":[")
	if len(t.items) > 0 {
		t.items[0].Json(out)
		if len(t.items) > 1 {
			for _, t := range t.items[1:] {
				out.WriteRune(',')
				t.Json(out)
			}
		}
	}
	out.WriteString("]}")
}

func (t *Concat) MaxLength() int {
	var total int
	for i, _ := range t.items {
		total += t.items[i].MaxLength()
	}
	return total
}

func NewSequence(items ...Generator) Generator { return &Concat{items: items} }

type Choice struct {
	items []Generator
}

// Expand randomly picks a component and forwards the call to it.
func (t *Choice) Expand(ctx context.Context, w io.StringWriter, env *Env) error {
	n := env.Prng.Intn(len(t.items))
	return t.items[n].Expand(ctx, w, env)
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
	out.WriteString("{\"type\":\"choice\",\"items\":[")
	if len(t.items) > 0 {
		t.items[0].Json(out)
		if len(t.items) > 1 {
			for _, t := range t.items[1:] {
				out.WriteRune(',')
				t.Json(out)
			}
		}
	}
	out.WriteString("]}")
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

func NewChoice(items ...Generator) Generator { return &Choice{items: items} }

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

func NewTerm(s string) Generator { t := Term(s); return &t }

type Empty struct{}

// Expand generates nothing
func (t *Empty) Expand(_ context.Context, _ io.StringWriter, _ *Env) error { return nil }

func (t *Empty) Encode(out *strings.Builder) {}

func (t *Empty) Json(out *strings.Builder) { out.WriteString("nil") }

func (t *Empty) MaxLength() int { return 0 }

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
