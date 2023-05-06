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

import "testing"

func test_parser2(t *testing.T, source, encoded string) {
	gen, err := ParseExpression(source)
	if err != nil {
		t.Fatal(err)
	}
	if gen == nil {
		t.Fatal("nil generator")
	}
	recoded := EncodeGenerator(gen)
	t.Logf("gen: '%v' -> %s", gen, JsonGenerator(gen))
	t.Logf("source: \"%s\" encoded: \"%v\"", source, recoded)
	if encoded != recoded {
		t.Fatal("unexpected result")
	}
}

func test_parser(t *testing.T, source string) {
	test_parser2(t, source, source)
}

func TestParser_SimpleTerm(t *testing.T) {
	test_parser(t, "plop mmlm mxls")
}

func TestParser_SimpleChoice(t *testing.T) {
	test_parser(t, "<plop|mmlm|mxls>")
	test_parser2(t, "<>", "")
	test_parser2(t, "< >", " ")
	test_parser2(t, "<plop|>", "plop")
	test_parser2(t, "<plop| >", "<plop| >")
}

func TestParser_SimpleConcat(t *testing.T) {
	test_parser(t, "plop <mmlm|mxls>")
}
