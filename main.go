package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

func execTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	w.Header().Add("Content-Type", "text/html")

	filename := "template/" + name
	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	inputData, err := ioutil.ReadFile(absFilename)
	if err != nil {
		return err
	}

	output := blackfriday.MarkdownCommon(inputData)

	t, err := setupTemplate("base.html").ParseFiles("template/base.html")
	if err != nil {
		return err
	}

	return t.Execute(w, map[string]interface{}{
		"Content": string(output),
	})
}

func execPage(w http.ResponseWriter, name string) {
	if err := execTemplate(w, name, map[string]interface{}{}); err != nil {
		http.Error(w, "Can't execute template: "+err.Error(), 500)
	}
}

func execBlog(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	articleName := vars["name"]
	fmt.Println("Article name: ", articleName)
	if strings.Contains(articleName, "/") { // Avoid people hacking the FS by loading urls like blog/../../.../etc/passwd
		http.Error(w, "Blog article not found: "+articleName, 404)
		return
	}
	fpath := "blog/" + articleName

	f, err := os.Open("template/" + fpath)
	f.Close()
	if err != nil {
		http.Error(w, "Blog article not found: "+articleName, 404)
		return
	}

	if err := execTemplate(w, fpath, map[string]interface{}{}); err != nil {
		http.Error(w, "Can't execute template: "+err.Error(), 500)
	}
}

func sendCSS(w http.ResponseWriter) {

	var internal = func(w http.ResponseWriter) error {
		inputData, err := ioutil.ReadFile("template/main.css")
		if err != nil {
			return err
		}

		colorsData, err := ioutil.ReadFile("template/colors.json")
		if err != nil {
			return err
		}

		var colors map[string]string

		if err := json.Unmarshal(colorsData, &colors); err != nil {
			return err
		}

		output := string(inputData)
		for k, v := range colors {
			output = strings.Replace(output, k, v, -1)
		}
		fmt.Fprint(w, output)
		return nil
	}
	if err := internal(w); err != nil {
		http.Error(w, "Can't execute template: "+err.Error(), 500)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) { execPage(w, "index.html") })
	r.HandleFunc("/about", func(w http.ResponseWriter, req *http.Request) { execPage(w, "about.html") })
	r.HandleFunc("/projects", func(w http.ResponseWriter, req *http.Request) { execPage(w, "projects.html") })
	r.HandleFunc("/thoughts", func(w http.ResponseWriter, req *http.Request) { execPage(w, "thoughts.html") })
	r.HandleFunc("/blog", func(w http.ResponseWriter, req *http.Request) { execPage(w, "articles.html") })
	r.HandleFunc("/blog/", func(w http.ResponseWriter, req *http.Request) { execPage(w, "articles.html") })
	r.HandleFunc("/blog/{name}/", execBlog)
	r.HandleFunc("/blog/{name}", execBlog)
	r.HandleFunc("/main.css", func(w http.ResponseWriter, req *http.Request) { sendCSS(w) })

	n := negroni.New()
	n.Use(negroni.NewRecovery())

	n.Use(negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	log.Fatal(http.ListenAndServe(":8080", n))
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
