package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"

	"github.com/gorilla/mux"
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
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 503)
		return
	}
	genes := strings.Split(string(data), "\n")
	genestring := strings.Join(genes, "+")

	url := "/userlist/result/" + header.Filename + "/" + genestring
	w.Write([]byte(url))

}

type UserScore struct {
	Name    string
	Tissues []string
	Genes   []string
	Scores  map[string][]mixt.UserScore
}

func UserListResultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genes := vars["genes"]
	listname := vars["listname"]

	fmt.Println(genes, listname)

	genelist := strings.Split(genes, "+")

	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("getting tissues went bad:", err)
		http.Error(w, err.Error(), 503)
		return
	}

	scores := make(map[string][]mixt.UserScore, 0)

	for _, tissue := range tissues {
		sc, err := mixt.UserEnrichmentScores(tissue, genelist)
		fmt.Println(sc, err)

		scores[tissue] = sc

	}

	us := UserScore{listname, tissues, genelist, scores}

	listResultTemplate.Execute(w, us)
}
