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
	"github.com/jfsmig/daily/excuse"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func main() {
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

var templateIndexText = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><title>Daily Excuse</title><style>
h1 { font-family: "Comic Sans MS", "Comic Sans", "Chalkboard SE", "Comic Neue", sans-serif; }
</style></head><body itemscope itemtype="http://schema.org/WebPage"><main><h1>{{.Excuse}}</h1></main></body></html>`
