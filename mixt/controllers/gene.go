package controllers

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/fjukstad/kvik/genecards"
	"github.com/gorilla/mux"
)

var geneTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
	"views/gene.html", "views/footer.html"))

type Genes struct {
	Genes []string
}

func GeneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	term := vars["genes"]
	genes := strings.Split(term, " ")
	var result []string
	for _, gene := range genes {
		hits, _ := SearchForGene(gene)
		result = append(result, hits...)
	}
	res := Genes{result}
	geneTemplate.Execute(w, res)
}

func GeneSummaryHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	geneName := vars["gene"]

	summary := genecards.Summary(geneName)

	w.Write([]byte(summary))

}
