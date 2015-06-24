package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"
)

type Tissues struct {
	T []Tissue
}
type Tissue struct {
	Name        string
	Comparisons []string
}

var tissuesTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/tissues.html", "views/footer.html"))

var tissueComparisonTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/tissue-comparison.html", "views/footer.html"))

func TissuesHandler(w http.ResponseWriter, r *http.Request) {
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("getting tissues went bad:", err)
		http.Error(w, err.Error(), 500)
	}

	var result []Tissue
	for i, t := range tissues {
		comp := make([]string, len(tissues))
		for j, u := range tissues {
			if i != j {
				comp[j] = t + "/" + u
			}
		}
		tissue := Tissue{t, comp}
		result = append(result, tissue)
	}
	res := Tissues{result}
	tissuesTemplate.Execute(w, res)
}

func TissueComparisonHandler(w http.ResponseWriter, r *http.Request) {
	tissueComparisonTemplate.Execute(w, nil)
}
