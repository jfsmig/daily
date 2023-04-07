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
)

type Env struct {
	Prng *rand.Rand
}

func NewEnv(seed int64) *Env { return &Env{Prng: rand.New(rand.NewSource(seed))} }

type Node interface {
	Expand(ctx context.Context, w io.StringWriter, env *Env) error
	Weight() uint64
}

type Concat struct {
	items []Node
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

func (t *Concat) Weight() uint64 {
	total := uint64(1)
	for idx, _ := range t.items {
		w := t.items[idx].Weight()
		total = total * w
	}
	return total
}

func NewSequence(items ...Node) Node { return &Concat{items: items} }

type Or struct {
	items []Node
}

func (t *Or) Expand(ctx context.Context, w io.StringWriter, env *Env) error {
	n := env.Prng.Intn(len(t.items))
	return t.items[n].Expand(ctx, w, env)
}

func (t *Or) Weight() uint64 {
	total := uint64(0)
	for idx, _ := range t.items {
		w := t.items[idx].Weight()
		total += w
	}
	return total
}

func NewChoice(items ...Node) Node { return &Or{items: items} }

type Term string

func (t *Term) Expand(ctx context.Context, w io.StringWriter, env *Env) error {
	_, err := w.WriteString(string(*t))
	return err
}

func (t *Term) Weight() uint64 {
	return 1
}

func NewTerm(s string) Node { t := Term(s); return &t }

func NewGenerator() (Node, error) {
	items := make([]Node, 0)
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
