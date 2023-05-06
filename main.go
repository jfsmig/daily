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
	"fmt"
	ht "html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	tt "text/template"
	"time"

	"github.com/jfsmig/daily/excuse"
)

const defaultTimeSlotRegen time.Duration = 5 * time.Minute
const defaultTimeSlotRefresh time.Duration = defaultTimeSlotRegen + time.Second

//go:embed shrug-emoticon.png
var icon []byte

//go:embed robots.txt
var robots []byte

//go:embed index.html
var templateIndexText string

//go:embed sitemap.xml
var templateSitemapText string

type HandlerFunc func(http.ResponseWriter, *http.Request)

func main() {
	excuseAny, err := newGenerator()
	if err != nil {
		log.Fatalln("excuse init error: ", err)
	}
	excuseOOO, err := newOOO()
	if err != nil {
		log.Fatalln("excuse init error: ", err)
	}
	excuseMeeting, err := newNoMeeting()
	if err != nil {
		log.Fatalln("excuse init error: ", err)
	}

	tplMain := ht.Must(ht.New("index").Parse(templateIndexText))

	generateExcuse := func(w http.ResponseWriter, req *http.Request, gen excuse.Generator) {
		type Args struct {
			Excuse  string
			Refresh int64
			Seed    int64
		}

		seed := int64(0)
		if req.URL.Query().Has("seed") {
			if s, err := strconv.ParseInt(req.URL.Query().Get("seed"), 10, 63); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				seed = s
			}
		} else {
			if req.URL.Query().Has("raw") {
				seed = time.Now().UnixNano()
			} else {
				seed = time.Now().Truncate(defaultTimeSlotRegen).UnixNano()
			}
		}

		env := excuse.NewEnv(seed)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		if req.URL.Query().Has("debug") {
			var sb strings.Builder
			gen.Json(&sb)
			w.Header().Add("Content-Type", "text/plain")
			if _, err := w.Write([]byte(sb.String())); err != nil {
				log.Println("Json rendering error:", err)
			}
		} else {
			var sb strings.Builder
			_ = gen.Expand(req.Context(), &sb, env)
			if req.URL.Query().Has("raw") {
				w.Header().Add("Content-Type", "text/plain")
				if _, err := fmt.Fprintf(w, "#%d\n%s", seed, sb.String()); err != nil {
					log.Println("Raw rendering error:", err)
				}
			} else {
				args := Args{
					Excuse:  sb.String(),
					Refresh: int64(defaultTimeSlotRefresh.Seconds()),
					Seed:    seed,
				}
				w.Header().Add("Content-Type", "text/html")
				if err := tplMain.Execute(w, args); err != nil {
					log.Println("Template rendering error:", err)
				}
			}
		}
	}

	// A set of routes providing the excuse as a "splash" html page
	http.HandleFunc("/ooo", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuse(w, req, excuseOOO) }
	}())
	http.HandleFunc("/meeting", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuse(w, req, excuseMeeting) }
	}())
	http.HandleFunc("/any", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuse(w, req, excuseAny) }
	}())

	http.HandleFunc("/sitemap.xml", func() HandlerFunc {
		type Args struct {
			Date string
		}
		tpl := tt.Must(tt.New("sitemap").Parse(templateSitemapText))
		args := Args{Date: time.Now().Truncate(25 * time.Hour).Format(time.RFC3339)}
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "application/xml")
			w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
			if err := tpl.Execute(w, args); err != nil {
				log.Println("Sitemap rendering error:", err)
			}
		}
	}())
	http.HandleFunc("/favicon.png", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "image/png")
			w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week
			if _, err := w.Write(icon); err != nil {
				log.Println("Icon reply error:", err)
			}
		}
	}())
	http.HandleFunc("/robots.txt", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week
			if _, err := w.Write(robots); err != nil {
				log.Println("Robots reply error:", err)
			}
		}
	}())

	// By default, the landing page proposes
	http.HandleFunc("/", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuse(w, req, excuseAny) }
	}())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("http server error:", err)
	}
}
