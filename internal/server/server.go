// Package server contains server-side logic
package server

import (
	"fmt"
	"net/http"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/Jasstkn/link-checker/internal/linkchecker"
)

// Output struct to pass into html as value for template
type Output struct {
	Description string
}

var tpl = template.Must(template.ParseFiles("templates/index.html"))

// CheckHandler parses form and renders reply
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	var output Output
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	url := r.FormValue("url")
	log.WithFields(log.Fields{"url": url}).Info("Starting processing url")
	result, err := linkchecker.LinkChecker(url)

	if err != nil {
		log.WithFields(log.Fields{"url": url}).Errorf("Error while processing url, %q", err)
		w.WriteHeader(http.StatusBadRequest)
		output = Output{Description: err.Error()}
	} else {
		log.WithFields(log.Fields{"url": url}).Info("Processed url")
		output = Output{Description: result}
	}
	tpl.Execute(w, output)
}

// IndexHandler serves main page with form
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// Init inits server and handlers
func Init() {
	fs := http.FileServer(http.Dir("./assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/check", CheckHandler)
	log.Info("starting server")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
