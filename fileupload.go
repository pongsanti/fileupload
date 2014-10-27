package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"html/template"
)

var templates = template.Must(template.ParseFiles("./template/upload.html"))

type Image struct {
	Filename string
}

func main() {
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./uploaded"))))

	http.HandleFunc("/form/", formHandler)
	http.HandleFunc("/upload/", uploadHandler)
	http.ListenAndServe(":8080", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "upload", "")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	out, err := os.Create("./uploaded/" + header.Filename)
	if err != nil {
		fmt.Fprintln(w, "Unalble to create file.")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	renderTemplate(w, "upload", header.Filename)
}

func renderTemplate(w http.ResponseWriter, tmpl string, filename string) {
    err := templates.ExecuteTemplate(w, tmpl+".html", Image{filename})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}