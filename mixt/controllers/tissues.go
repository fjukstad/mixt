package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/fjukstad/kvik/utils"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"
)

var tissuesTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
	"views/tissues.html", "views/footer.html"))

func TissuesHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func GetTissues(w http.ResponseWriter, r *http.Request) {
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("getting tissues went bad:", err)
		http.Error(w, err.Error(), 500)
	}
	res := utils.SearchResponse{tissues}
	b, _ := json.Marshal(res)
	w.Write(b)
}
