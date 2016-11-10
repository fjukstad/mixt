package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"text/template"

	"github.com/fjukstad/mixt/mixt"

	"github.com/gorilla/mux"
)

var searchResultTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/search-result.html", "views/footer.html"))

type SearchResponse struct {
	Terms []string
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	term := vars["term"]
	result, err := SearchForGene(term)
	if err != nil {
		fmt.Println("Search went bad:", err)
		errorHandler(w, r, err)
		return
	}

	setres, err := SearchForGeneSet(term)
	if err != nil {
		fmt.Println("Searching for gene sets bad", err)
		errorHandler(w, r, err)
		return
	}

	result = append(result, setres...)

	GOTermNames, err := SearchForGOTerms(term)
	if err != nil {
		fmt.Println("Searching for go terms went bad", err)
		errorHandler(w, r, err)
		return
	}

	result = append(result, GOTermNames...)

	res := SearchResponse{result}
	b, _ := json.Marshal(res)
	w.Write(b)
}

var genes []string

func SearchForGene(searchTerm string) ([]string, error) {
	var err error
	if len(genes) < 1 {
		genes, err = mixt.GetGenes()
		if err != nil {
			return []string{}, err
		}
	}

	result := inSlice(searchTerm, genes)

	return result, nil
}

var geneSetNames []string

func SearchForGeneSet(searchTerm string) ([]string, error) {
	var err error
	if len(geneSetNames) < 1 {
		geneSetNames, err = mixt.GetGeneSetNames()
		if err != nil {
			return []string{}, err
		}
	}

	results := inSlice(searchTerm, geneSetNames)

	return results, nil
}

var GOTermNames []string

func SearchForGOTerms(searchTerm string) ([]string, error) {
	var err error
	if len(GOTermNames) < 1 {
		GOTermNames, err = mixt.GetGOTermNames()
		if err != nil {
			return []string{}, err
		}
	}
	results := inSlice(searchTerm, GOTermNames)
	return results, nil

}

func inSlice(s string, words []string) []string {
	var result []string
	for _, word := range words {
		a := strings.ToLower(word)
		b := strings.ToLower(s)
		if strings.Contains(a, b) {
			wordFmt := strings.Trim(word, "\"")

			jresult := strings.Join(result, " ")
			if !strings.Contains(jresult, wordFmt) {
				result = append(result, wordFmt)
			}
		}
	}
	return result
}

type SearchResults struct {
	Genes    []Gene
	Tissues  []string
	GeneSets []GeneSet
	GOTerms  []GOTerm
}

func SearchResultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	term := vars["terms"]
	searchTerms := strings.Split(term, " ")

	genes, tissues, err := GeneResults(searchTerms)
	if err != nil {
		fmt.Println("gene results error", err)
		errorHandler(w, r, err)
		return
	}

	geneSets, err := SetResults(searchTerms)

	if err != nil {
		fmt.Println("set results", err)
		errorHandler(w, r, err)
		return
	}

	goterms, err := GOTermResults(searchTerms)
	if err != nil {
		fmt.Println("go term results", err)
		errorHandler(w, r, err)
		return
	}

	res := SearchResults{genes, tissues, geneSets, goterms}
	searchResultTemplate.Execute(w, res)
}
