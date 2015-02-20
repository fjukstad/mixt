package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/vdumeaux/mixt/mixt/models"

	"github.com/fjukstad/kvik/dataset"
	"github.com/fjukstad/kvik/utils"
	"github.com/gorilla/mux"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	term := vars["term"]

	genes, err := models.GetGenes(d)
	if err != nil {
		fmt.Println("Search went bad:", err)
		http.Error(w, err.Error(), 500)
	}

	var result []string
	for _, gene := range genes {
		a := strings.ToLower(gene)
		b := strings.ToLower(term)
		if strings.Contains(a, b) {
			geneFmt := strings.Trim(gene, "\"")
			result = append(result, geneFmt)
		}
	}

	fmt.Println(genes, term)

	res := utils.SearchResponse{result}
	b, _ := json.Marshal(res)
	w.Write(b)

}

func InitSearch(dataset *dataset.Dataset) {
	d = dataset
}
