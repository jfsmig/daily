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
	Debug(out *strings.Builder)
}

type Concat struct {
	items []Generator
}

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

func (t *Concat) Debug(out *strings.Builder) {
	out.WriteString(" Sequence(")
	if len(t.items) > 0 {
		t.items[0].Debug(out)
		if len(t.items) > 1 {
			for _, t := range t.items[1:] {
				out.WriteRune(',')
				t.Debug(out)
			}
		}
	}
	out.WriteRune(')')
}

func NewSequence(items ...Generator) Generator { return &Concat{items: items} }

type Choice struct {
	items []Generator
}

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

func (t *Choice) Debug(out *strings.Builder) {
	out.WriteString(" Choice(")
	if len(t.items) > 0 {
		t.items[0].Debug(out)
		if len(t.items) > 1 {
			for _, t := range t.items[1:] {
				out.WriteRune('|')
				t.Debug(out)
			}
		}
	}
	out.WriteRune(')')
}

func NewChoice(items ...Generator) Generator { return &Choice{items: items} }

type Term string

func (t *Term) Expand(ctx context.Context, w io.StringWriter, env *Env) error {
	_, err := w.WriteString(string(*t))
	return err
}

func (t *Term) Encode(out *strings.Builder) { out.WriteString(string(*t)) }

func (t *Term) Debug(out *strings.Builder) {
	out.WriteString(" Term(")
	out.WriteString(string(*t))
	out.WriteRune(')')
}

func NewTerm(s string) Generator { t := Term(s); return &t }

type Empty struct{}

func (t *Empty) Expand(_ context.Context, _ io.StringWriter, _ *Env) error { return nil }

func (t *Empty) Encode(out *strings.Builder) {}

func (t *Empty) Debug(out *strings.Builder) { out.WriteString(" Empty()") }

func NewEmpty() Generator { return &Empty{} }

func NewGenerator() (Generator, error) {
	items := make([]Generator, 0)
	if n, err := NewNoMeeting(); err != nil {
		return nil, err
	} else {
		items = append(items, n)
	}
	if n, err := NewOOO(); err != nil {
		return nil, err
	} else {
		items = append(items, n)
	}
	return NewChoice(items...), nil
}

func EncodeGenerator(g Generator) string {
	str := strings.Builder{}
	g.Encode(&str)
	return str.String()
}

func DebugGenerator(g Generator) string {
	str := strings.Builder{}
	g.Debug(&str)
	return str.String()
}
