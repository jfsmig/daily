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
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/jfsmig/daily/excuse"
)

//go:embed shrug-emoticon.png
var icon []byte

//go:embed robots.txt
var robots []byte

//go:embed index.html
var templateIndexText string

type HandlerFunc func(http.ResponseWriter, *http.Request)

func main() {
	http.HandleFunc("/favicon.png", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Write(icon)
		}
	}())

	http.HandleFunc("/robots.txt", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Write(robots)
		}
	}())

	http.HandleFunc("/", func() HandlerFunc {
		type Args struct {
			Excuse string
		}
		nodaily, err := excuse.NewJohn()
		if err != nil {
			log.Fatalln("excuse init error: ", err)
		}
		tpl := template.Must(template.New("index").Parse(templateIndexText))
		return func(w http.ResponseWriter, req *http.Request) {
			env := excuse.NewEnv()
			var sb strings.Builder
			_ = nodaily.Expand(req.Context(), &sb, env)
			args := Args{Excuse: sb.String()}
			if err := tpl.Execute(w, args); err != nil {
				log.Println("Template rendering error:", err)
				w.WriteHeader(500)
			}
		}
	}())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("http server error:", err)
	}
}
