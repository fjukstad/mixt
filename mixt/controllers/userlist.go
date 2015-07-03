package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

var listSubmitTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/userlist-submit.html", "views/footer.html"))

var listResultTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/userlist-result.html", "views/footer.html"))

func UserListHandler(w http.ResponseWriter, r *http.Request) {

	listSubmitTemplate.Execute(w, nil)
}

func UserListSubmitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling result")

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 503)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 503)
	}
	genes := strings.Split(string(data), "\n")
	genestring := strings.Join(genes, "+")

	url := "/userlist/result/" + header.Filename + "/" + genestring
	w.Write([]byte(url))

}

func UserListResultHandler(w http.ResponseWriter, r *http.Request) {
	listResultTemplate.Execute(w, nil)
}
