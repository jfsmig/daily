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

package main

import (
	"context"
	"github.com/jfsmig/daily/excuse"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func testGenerator(t *testing.T, genStatement, genCause excuse.Generator) {
	if maxLength := genStatement.MaxLength(); maxLength >= 32 {
		t.Fatalf("Too long (max=32 gen=%d)", maxLength)
	}
	if maxLength := genCause.MaxLength(); maxLength >= 128 {
		t.Fatalf("Too long (max=128 gen=%d)", maxLength)
	}

	ctx := context.TODO()
	env := excuse.Env{Prng: rand.New(rand.NewSource(time.Now().UnixNano()))}

	var sb strings.Builder
	for i := 0; i < 10; i++ {
		sb.Reset()
		if err := genStatement.Expand(ctx, &sb, &env); err != nil {
			t.Fatal(err)
		}
		t.Log(">", sb.String())
	}
	for i := 0; i < 10; i++ {
		sb.Reset()
		if err := genCause.Expand(ctx, &sb, &env); err != nil {
			t.Fatal(err)
		}
		t.Log(">", sb.String())
	}
}

func TestMeeting(t *testing.T) {
	statement, cause := newNoMeeting()
	testGenerator(t, statement, cause)
}

func TestOOO(t *testing.T) {
	statement, cause := newOOO()
	testGenerator(t, statement, cause)
}

func TestLoad(t *testing.T) {
	mux := initHttp()
	if mux == nil {
		t.Fatal("invalid http multiplexer")
	}
}
