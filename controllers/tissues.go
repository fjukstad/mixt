package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"

	"github.com/fjukstad/mixt/mixt"
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
		errorHandler(w, r, err)
		return
	}

	var result []Tissue
	for _, t := range tissues {
		comp := make([]string, len(tissues))
		for j, u := range tissues {
			//if i != j {
			comp[j] = t + "/" + u
			//}
		}
		tissue := Tissue{t, comp}
		result = append(result, tissue)
	}
	res := Tissues{result}
	tissuesTemplate.Execute(w, res)
}

type TissueComparison struct {
	TissueA string
	TissueB string
	Cohorts []string
}

func TissueComparisonHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tissueA := vars["tissueA"]
	tissueB := vars["tissueB"]

	cohorts, err := mixt.GetCohorts()
	if err != nil {
		fmt.Println("Could not get cohorts")
		errorHandler(w, r, err)
		return
	}

	tissueComparisonTemplate.Execute(w, TissueComparison{tissueA, tissueB, cohorts})
}

func AnalysisHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tissueA := vars["tissueA"]
	tissueB := vars["tissueB"]

	analysis := vars["analysis"]
	cohort := vars["cohort"]

	fmt.Println("analysis,cohort", analysis, cohort)

	var result []byte
	var err error

	if analysis == "ranksum" {
		result, err = mixt.PatientRankSum(tissueA, tissueB, cohort)
		if err != nil {
			fmt.Println("Could not get patient rank sum correlation")
			errorHandler(w, r, err)
			return
		}
	}

	if analysis == "eigengene" {
		result, err = mixt.EigengeneCorrelation(tissueA, tissueB)
		if err != nil {
			fmt.Println("Could not do eigen gene correlation", err)
			errorHandler(w, r, err)
			return
		}
	}

	if analysis == "overlap" {
		result, err = mixt.GeneOverlapTest(tissueA, tissueB)
		if err != nil {
			errorHandler(w, r, err)
			return
		}
	}

	if analysis == "roi" {
		result, err = mixt.ROITest(tissueA, tissueB)
		if err != nil {
			errorHandler(w, r, err)
			return
		}
	}

	if analysis == "patientrank" {
		result, err = mixt.PatientRankCorrelation(tissueA, tissueB)
		if err != nil {
			errorHandler(w, r, err)
			return
		}
	}

	w.Write(result)

}

func TOMGraphHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tissue := vars["tissue"]
	what := vars["what"]

	res, err := mixt.GetTOMGraph(tissue, what)
	if err != nil {
		errorHandler(w, r, err)
		return
	}
	w.Write(res)

}
