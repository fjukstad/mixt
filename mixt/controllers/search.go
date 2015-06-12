package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"

	"github.com/gorilla/mux"
)

type SearchResponse struct {
	Terms []string
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if !LoggedIn(r) {
		http.Redirect(w, r, "/", 302)
		return
	}

	vars := mux.Vars(r)
	term := vars["term"]
	result, err := SearchForGene(term)
	if err != nil {
		fmt.Println("Search went bad:", err)
		http.Error(w, err.Error(), 500)
	}

	setres, err := SearchForGeneSet(term)
	if err != nil {
		fmt.Println("Searching for gene sets bad", err)
		http.Error(w, err.Error(), 500)
	}

	result = append(result, setres...)

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

	result := inSlice(searchTerm, geneSetNames)

	return result, nil
}

func inSlice(s string, words []string) []string {
	var result []string
	for _, word := range words {
		a := strings.ToLower(word)
		b := strings.ToLower(s)
		if strings.Contains(a, b) {
			wordFmt := strings.Trim(word, "\"")
			result = append(result, wordFmt)
		}
	}
	return result
}
