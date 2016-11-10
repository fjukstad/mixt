package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fjukstad/kvik/genecards"
	"github.com/fjukstad/mixt/mixt"
	"github.com/gorilla/mux"
)

type Genes struct {
	Genes   []Gene
	Tissues []string
}

type Gene struct {
	Name    string
	Modules []string
	Summary string
}

func GeneResults(genes []string) ([]Gene, []string, error) {
	var result []Gene
	for _, gene := range genes {
		hits, _ := SearchForGene(gene)
		fmt.Println(hits)
		for _, h := range hits {

			modules, err := mixt.GetAllModuleNames(h)
			if err != nil {
				fmt.Println("Could not get modules for ", h)
				return []Gene{}, []string{}, err
			}

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
		//http.Error(w, err.Error(), 503)
		return []Gene{}, []string{}, err
	}

	return result, tissues, nil
}

func GeneSummaryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	geneName := vars["gene"]

	var summary string
	summary = genecards.Summary(geneName)
	if summary == "" {
		summary = "no preview available..."
	}

	w.Write([]byte(summary))

}

type Common struct {
	Genes []string
}

func CommonGenesHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	module := vars["module"]
	tissue := vars["tissue"]
	geneset := vars["geneset"]
	status := vars["status"]
	output := vars["output"]

	commonGenes, err := mixt.GetCommonGenes(tissue, module, geneset, status)
	if err != nil {
		fmt.Println("Could not get common genes", err)
		errorHandler(w, r, err)
		return
	}

	common := Common{commonGenes}

	if output == "json" {
		b, err := json.Marshal(common)
		if err != nil {
			fmt.Println("Could not marshal common genes", err)
			errorHandler(w, r, err)
			return
		}

		w.Write(b)
		return
	}

	header := []string{"Gene"}
	var records [][]string
	records = append(records, header)
	for _, gene := range commonGenes {
		entry := []string{gene}
		records = append(records, entry)
	}

	writer := csv.NewWriter(w)
	writer.WriteAll(records)

}

func GeneList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modules := vars["modules"]
	tissue := vars["tissue"]

	parsedModules := strings.Split(modules, "+")
	csv, err := mixt.GeneListCSV(parsedModules, tissue)
	if err != nil {
		errorHandler(w, r, err)
		return
	}
	filename := tissue + "-" + modules + ".csv"
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Write(csv)
	return
}
