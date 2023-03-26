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
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jfsmig/daily/excuse"
)

var nodialy excuse.Node

func main() {
	var err error
	nodialy, err = excuse.NewJohn()
	http.HandleFunc("/excuse", doExcuse)
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("http server error:", err)
	}
}

func doExcuse(w http.ResponseWriter, req *http.Request) {
	var buf strings.Builder

	fmt.Fprint(w, buf.String())
}
