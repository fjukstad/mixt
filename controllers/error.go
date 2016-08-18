package controllers

import (
	"net/http"
	"text/template"
)

type Message struct {
	Error string
}

var errorTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/500.html", "views/footer.html"))

func errorHandler(w http.ResponseWriter, r *http.Request, e error) {
	message := Message{e.Error()}
	errorTemplate.Execute(w, message)
}
