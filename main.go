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
	"embed"
	"fmt"
	ht "html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jfsmig/daily/excuse"
)

const defaultTimeSlotRegen time.Duration = 5 * time.Minute
const defaultTimeSlotRefresh time.Duration = defaultTimeSlotRegen + time.Second

//go:embed templates/*
//go:embed data/*
var payload embed.FS

type HandlerFunc func(http.ResponseWriter, *http.Request)

func must[T any](b T, e error) T {
	if e != nil {
		panic(e)
	} else {
		return b
	}
}

func load(name string) []byte { return must(payload.ReadFile(name)) }

func loadS(name string) string { return string(load(name)) }

// newOOO returns <Statement,Cause> plus an optional error
func newOOO() (excuse.Generator, excuse.Generator) {
	statement := excuse.NewChoice(
		excuse.NewTerm("I'm going to be OOO,"),
		excuse.NewTerm("I need to be OOO today,"),
		excuse.NewTerm("I can't show up today,"))
	cause := must(excuse.ParseStreamString(loadS("data/ooo.txt")))
	return statement, cause
}

// newNoMeeting returns <Statement,Cause> plus an optional error
func newNoMeeting() (excuse.Generator, excuse.Generator) {
	statement := excuse.NewChoice(
		excuse.NewTerm("I cannot attend the daily,"),
		excuse.NewTerm("going to miss the meeting,"),
		excuse.NewTerm("gonna miss the meeting,"),
		excuse.NewTerm("No daily meeting for me,"))
	cause := must(excuse.ParseStreamString(loadS("data/meeting.txt")))
	return statement, cause
}

func initHttp() http.Handler {
	// Load all the structure coming from embedded files
	statementOOO, excuseOOO := newOOO()
	statementMeeting, excuseMeeting := newNoMeeting()
	tplIndex := must(ht.New("index").Parse(loadS("templates/index.html")))
	tplSplash := must(ht.New("splash").Parse(loadS("templates/splash.html")))
	tplSitemap := must(ht.New("sitemap").Parse(loadS("templates/sitemap.xml")))
	robotsBytes := load("templates/robots.txt")
	iconBytes := load("templates/shrug-emoticon.png")

	generateExcuse := func(w http.ResponseWriter, req *http.Request, genStatement, genCause excuse.Generator) {
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

		var sbCause, sbStatement strings.Builder
		env := excuse.NewEnv(seed)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("X-daily-seed", strconv.FormatInt(seed, 10))
		if req.URL.Query().Has("debug") {
			genStatement.Json(&sbStatement)
			genCause.Json(&sbCause)
			w.Header().Add("Content-Type", "text/plain")
			_, _ = fmt.Fprintf(w, "%v\n%s\n", genStatement.Count(), sbStatement.String())
			_, _ = fmt.Fprintf(w, "%v\n%s\n", genCause.Count(), sbCause.String())
		} else {
			_ = genStatement.Expand(req.Context(), &sbStatement, env)
			_ = genCause.Expand(req.Context(), &sbCause, env)
			if req.URL.Query().Has("raw") {
				w.Header().Add("Content-Type", "text/plain")
				if _, err := w.Write([]byte(sbStatement.String())); err != nil {
					log.Println("Raw rendering error:", err)
				}
				if _, err := w.Write([]byte(sbCause.String())); err != nil {
					log.Println("Raw rendering error:", err)
				}
			} else {
				type Args struct {
					Seed      int64
					Refresh   int64
					Statement string
					Excuse    string
				}
				args := Args{
					Seed:      seed,
					Refresh:   int64(defaultTimeSlotRefresh.Seconds()),
					Statement: sbStatement.String(),
					Excuse:    sbCause.String(),
				}
				w.Header().Add("Content-Type", "text/html")
				if err := tplSplash.Execute(w, args); err != nil {
					log.Println("Template rendering error:", err)
				}
			}
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ooo", func() HandlerFunc {
		// A set of routes providing the excuse as a "splash" html page
		return func(w http.ResponseWriter, req *http.Request) {
			generateExcuse(w, req, statementOOO, excuseOOO)
		}
	}())

	mux.HandleFunc("/meeting", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			generateExcuse(w, req, statementMeeting, excuseMeeting)
		}
	}())

	mux.HandleFunc("/sitemap.xml", func() HandlerFunc {
		type Args struct {
			Date string
		}
		args := Args{Date: time.Now().Truncate(25 * time.Hour).Format(time.RFC3339)}
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "application/xml")
			w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
			if err := tplSitemap.Execute(w, args); err != nil {
				log.Println("Sitemap rendering error:", err)
			}
		}
	}())

	mux.HandleFunc("/favicon.png", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "image/png")
			w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week
			if _, err := w.Write(iconBytes); err != nil {
				log.Println("Icon reply error:", err)
			}
		}
	}())

	mux.HandleFunc("/robots.txt", func() HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Set("Cache-Control", "public, max-age=604800") // 1 week
			if _, err := w.Write(robotsBytes); err != nil {
				log.Println("Robots reply error:", err)
			}
		}
	}())

	mux.HandleFunc("/", func() HandlerFunc {
		// By default, the landing page proposes
		type Args struct {
			Date string
		}
		return func(w http.ResponseWriter, req *http.Request) {
			args := Args{Date: time.Now().Truncate(25 * time.Hour).Format(time.RFC3339)}
			w.Header().Add("Content-Type", "text/html")
			if err := tplIndex.Execute(w, args); err != nil {
				log.Println("Template rendering error:", err)
			}
		}
	}())

	return mux
}

func main() {
	if err := http.ListenAndServe(":8080", initHttp()); err != nil {
		log.Fatalln("http server error:", err)
	}
}
