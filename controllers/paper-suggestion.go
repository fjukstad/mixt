package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/fjukstad/mixt-blood-tumor/mixt"

	"github.com/fjukstad/kvik/eutils"
	"github.com/gorilla/mux"
)

var paperSuggestionTemplate = template.Must(template.ParseFiles("views/base.html", "views/header.html", "views/outside-navbar.html",
	"views/paper-suggestion.html", "views/footer.html"))

type Results struct {
	Papers []eutils.DocumentSummary
}

func PaperSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tissue := vars["tissue"]
	module := vars["module"]

	m, err := mixt.GetModule(module, tissue, "all")
	if err != nil {
		fmt.Println("Could not get module", err)
		errorHandler(w, r, err)
		return
	}

	genelist := []string{}
	for _, gene := range m.Genes {
		genelist = append(genelist, gene.Name)
	}

	genelist = genelist[0:10]

	fmt.Println(genelist, module, tissue)

	searchres, err := eutils.Search(genelist)
	if err != nil {
		fmt.Println(err)
		return
	}

	//.PmcRefCount .ViewCount

	fmt.Println("SEARCHRES:", searchres)
	res, err := eutils.Summary(searchres)

	for _, paper := range res.DocumentSummarySet.DocumentSummary {
		fmt.Println(paper)
	}

	paperSuggestionTemplate.Execute(w, Results{res.DocumentSummarySet.DocumentSummary})

}
