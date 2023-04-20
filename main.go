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
	ht "html/template"
	"log"
	"net/http"
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
	excuseAny, err := excuse.NewGenerator()
	if err != nil {
		log.Fatalln("excuse init error: ", err)
	}
	excuseOOO, err := excuse.NewOOO()
	if err != nil {
		log.Fatalln("excuse init error: ", err)
	}
	excuseMeeting, err := excuse.NewNoMeeting()
	if err != nil {
		log.Fatalln("excuse init error: ", err)
	}

	generateExcuseRaw := func(w http.ResponseWriter, req *http.Request, e excuse.Generator) {
		var sb strings.Builder
		env := excuse.NewEnv(time.Now().UnixNano())
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		if err := e.Expand(req.Context(), &sb, env); err != nil {
			log.Println("Template rendering error:", err)
		} else {
			s := strings.Trim(sb.String(), " ")
			w.Write([]byte(s + "\n"))
		}
	}

	tplMain := ht.Must(ht.New("index").Parse(templateIndexText))

	generateExcuseHtml := func(w http.ResponseWriter, req *http.Request, gen excuse.Generator) {
		type Args struct {
			Excuse  string
			Refresh int64
		}
		// This will change the excuse each hour
		env := excuse.NewEnv(time.Now().Truncate(defaultTimeSlotRegen).UnixNano())
		var sb strings.Builder
		_ = gen.Expand(req.Context(), &sb, env)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		args := Args{
			Excuse:  sb.String(),
			Refresh: int64(defaultTimeSlotRefresh.Seconds()),
		}
		if err := tplMain.Execute(w, args); err != nil {
			log.Println("Template rendering error:", err)
		}
	}

	// A set of routes providing the excuse as a simple string
	http.HandleFunc("/raw/ooo", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuseRaw(w, req, excuseOOO) }
	}())
	http.HandleFunc("/raw/meeting", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuseRaw(w, req, excuseMeeting) }
	}())
	http.HandleFunc("/raw/any", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuseRaw(w, req, excuseAny) }
	}())

	// A set of routes providing the excuse as a "splash" html page
	http.HandleFunc("/html/ooo", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuseHtml(w, req, excuseOOO) }
	}())
	http.HandleFunc("/html/meeting", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuseHtml(w, req, excuseMeeting) }
	}())
	http.HandleFunc("/html/any", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateExcuseHtml(w, req, excuseAny) }
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
		return func(w http.ResponseWriter, req *http.Request) { generateExcuseHtml(w, req, excuseAny) }
	}())

	// Generate a long sentence that is helpful to diagnose design quirks.
	generateLoremHtml := func(w http.ResponseWriter, req *http.Request) {
		type Args struct {
			Excuse  string
			Refresh int64
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		args := Args{
			Excuse:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			Refresh: 3600,
		}
		if err := tplMain.Execute(w, args); err != nil {
			log.Println("Template rendering error:", err)
		}
	}

	http.HandleFunc("/w/lorem", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) { generateLoremHtml(w, req) }
	}())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("http server error:", err)
	}
}
