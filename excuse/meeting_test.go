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
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestGenerator(t *testing.T) {
	john, err := NewGenerator()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	env := Env{Prng: rand.New(rand.NewSource(time.Now().UnixNano()))}

	var sb strings.Builder
	for i := 0; i < 10; i++ {
		sb.Reset()
		if err = john.Expand(ctx, &sb, &env); err != nil {
			t.Fatal(err)
		}
		t.Log(">", sb.String())
	}
}
