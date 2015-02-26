package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"

	"github.com/fjukstad/kvik/genecards"
	"github.com/gorilla/mux"
)

var geneTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/gene.html", "views/footer.html"))

type Genes struct {
	Genes   []Gene
	Tissues []string
}

type Gene struct {
	Name    string
	Modules []string
	Summary string
}

func GeneHandler(w http.ResponseWriter, r *http.Request) {

	if !LoggedIn(r) {
		http.Redirect(w, r, "/", 302)
		return
	}

	fmt.Println("GENEHANDLER WHAT ARE YOU DOING?")
	vars := mux.Vars(r)
	term := vars["genes"]
	genes := strings.Split(term, " ")
	var result []Gene
	for _, gene := range genes {
		hits, _ := SearchForGene(gene)
		for _, h := range hits {

			modules, err := mixt.GetAllModules(h)
			if err != nil {
				fmt.Println("Could not get modules for ", h)
				http.Error(w, err.Error(), 503)
				return
			}

			fmt.Println(h, modules)

			summary := genecards.Summary(h)

			s := strings.SplitAfterN(summary, ".", 2)

			shortSummary := s[0] + ".."

			g := Gene{h, modules, shortSummary}
			result = append(result, g)
		}
	}

	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("Could not get tissues")
		http.Error(w, err.Error(), 503)
		return
	}

	res := Genes{result, tissues}
	geneTemplate.Execute(w, res)
}

func GeneSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if !LoggedIn(r) {
		http.Redirect(w, r, "/", 302)
		return
	}

	vars := mux.Vars(r)
	geneName := vars["gene"]

	var summary string
	summary = genecards.Summary(geneName)
	if summary == "" {
		summary = "no preview available..."
	}

	w.Write([]byte(summary))

}
