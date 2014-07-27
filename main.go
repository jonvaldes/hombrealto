package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		t, err := setupTemplate("index.html").ParseFiles("template/index.html")
		if err != nil {
			http.Error(w, "Can't parse template: "+err.Error(), 500)
			return
		}
		err = t.Execute(w, map[string]interface{}{
			"Lines": []string{
				"moooo",
				"wweeeeeee",
				"<a href=\"https://twitter.com/jon_valdes\">@jon_valdes</a> is sit amet vitae augue. Nam tincidunt congue enim, ut porta lorm",
			},
		})

		if err != nil {
			http.Error(w, "Can't execute template: "+err.Error(), 500)
			return
		}
	})

	n := negroni.New()
	n.Use(negroni.NewRecovery())

	n.Use(negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	log.Fatal(http.ListenAndServe(":8090", n))
}

var cleanLinksRegexA = regexp.MustCompile("<a href=\".*\">")
var cleanLinksRegexB = regexp.MustCompile("</a>")

func setupTemplate(name string) *template.Template {
	return template.New(name).Funcs(map[string]interface{}{
		"hostname":    func() string { h, _ := os.Hostname(); return h },
		"machineInfo": func() string { return runtime.GOOS + " " + runtime.GOARCH },
		"big": func(t string) string {
			chars := make([]string, len(t))
			for i := 0; i < len(t); i++ {
				chars[i] = string(t[i])
			}

			return strings.Join(chars, " ")
		},

		"center": func(l string) string {
			for len(l) <= 74 {
				l = " " + l + " "
			}
			return l
		},
		"line": func(l string) string {
			l = cleanLinksRegexA.ReplaceAllString(l, "")
			l = cleanLinksRegexB.ReplaceAllString(l, "")
			length := len(l)
			if length > 76 {
				length = 76
			}
			return "<span>\n│ " + fmt.Sprintf("%-76s", l[:length]) + " │</span>"
		},
	})
}
