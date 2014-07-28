package main

import (
	"bytes"
	"fmt"
	"html"
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

func execTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	w.Header().Add("Content-Type", "text/html")

	var temp bytes.Buffer
	t, err := setupTemplate(name).ParseFiles("template/" + name)
	if err != nil {
		return err
	}
	if err = t.Execute(&temp, data); err != nil {
		return err
	}

	t, err = setupTemplate("base.html").ParseFiles("template/base.html")
	if err != nil {
		return err
	}

	return t.Execute(w, map[string]interface{}{
		"Content": string(temp.Bytes()),
	})
}

func execPage(w http.ResponseWriter, name string) {
	if err := execTemplate(w, name, map[string]interface{}{}); err != nil {
		http.Error(w, "Can't execute template: "+err.Error(), 500)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) { execPage(w, "index.html") })
	r.HandleFunc("/about", func(w http.ResponseWriter, req *http.Request) { execPage(w, "about.html") })
	r.HandleFunc("/projects", func(w http.ResponseWriter, req *http.Request) { execPage(w, "projects.html") })

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
			nl := cleanLinksRegexA.ReplaceAllString(l, "")
			nl = cleanLinksRegexB.ReplaceAllString(nl, "")
			nl = html.UnescapeString(nl)

			cleanLength := len(nl)
			overflowChars := cleanLength - 76
			if overflowChars < 0 {
				overflowChars = 0
			}
			l = l[:len(l)-overflowChars]
			for i := cleanLength; i < 76; i++ {
				l += " "
			}

			return "<span>| " + fmt.Sprintf("%s", l+" |</span>")
		},
	})
}
