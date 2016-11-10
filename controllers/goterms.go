package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fjukstad/mixt/mixt"
	"github.com/gorilla/mux"
)

type GOTerm struct {
	Name      string
	Id        string
	TissueSet map[string][]mixt.GOTerm
}

func GOTermResults(searchTerms []string) ([]GOTerm, error) {
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("getting tissues went bad:", err)
		return []GOTerm{}, err
	}

	var GOTermNames []string
	for _, term := range searchTerms {
		hits, err := SearchForGOTerms(term)
		if err != nil {
			return []GOTerm{}, err
		}
		GOTermNames = append(GOTermNames, hits...)
	}

	var goterms []GOTerm

	for _, GOTermName := range GOTermNames {
		set := make(map[string][]mixt.GOTerm)
		var id string
		for _, tissue := range tissues {
			if tissue == "nblood" {
				continue
			}
			score, err := mixt.GetGOScoresForTissue(tissue, GOTermName)
			if err != nil {
				//return []GOTerm{}, err
				continue
			}
			set[tissue] = score
			//GOTermScore{score.Module, score.ClassicFisher,
			//	score.Weight01Fisher}
			id = score[0].GOId
		}
		goterms = append(goterms, GOTerm{GOTermName, id, set})
	}

	return goterms, nil
}

func CommonGOTermGenesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	module := vars["module"]
	tissue := vars["tissue"]
	id := vars["id"]
	format := vars["format"]

	genes, err := mixt.GetCommonGOTermGenes(module, tissue, id)
	if err != nil {
		fmt.Println("Could not common go term genes")
		errorHandler(w, r, err)
		return
	}

	if format == "json" {
		common := Common{genes}
		b, err := json.Marshal(common)
		if err != nil {
			fmt.Println("Could not marshal common go term genes")
			errorHandler(w, r, err)
			return
		}

		w.Write(b)
		return
	}

	header := []string{"Gene"}
	var records [][]string
	records = append(records, header)
	for _, gene := range genes {
		entry := []string{gene}
		records = append(records, entry)
	}

	writer := csv.NewWriter(w)
	writer.WriteAll(records)

}
