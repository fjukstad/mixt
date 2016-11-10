package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/fjukstad/mixt/mixt"

	"github.com/gorilla/mux"
)

var listSubmitTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/userlist-submit.html", "views/footer.html"))

var listResultTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/userlist-result.html", "views/footer.html"))

func UserListHandler(w http.ResponseWriter, r *http.Request) {

	listSubmitTemplate.Execute(w, nil)
}

func UserListSubmitHandler(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		errorHandler(w, r, err)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		errorHandler(w, r, err)
		return
	}
	genes := strings.Split(string(data), "\n")
	genestring := strings.Join(genes, "+")

	url := "/userlist/result/" + header.Filename + "/" + genestring
	w.Write([]byte(url))

}

type UserScore struct {
	Name    string
	Tissues []string
	Genes   []string
	Scores  map[string][]mixt.UserScore
}

func UserListResultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	genes := vars["genes"]
	listname := vars["listname"]

	genelist := strings.Split(genes, "+")

	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("getting tissues went bad:", err)
		errorHandler(w, r, err)
		return
	}

	scores := make(map[string][]mixt.UserScore, 0)

	for i, tissue := range tissues {
		if i < 2 {
			sc, err := mixt.UserEnrichmentScores(tissue, genelist)
			if err != nil {
				fmt.Println("Userlist enrichment bad", err)
				errorHandler(w, r, err)
				return
			}
			scores[tissue] = sc
		}
	}

	us := UserScore{listname, tissues, genelist, scores}

	listResultTemplate.Execute(w, us)
}

func UserListCommonGenesHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	genes := vars["genes"]
	module := vars["module"]
	tissue := vars["tissue"]
	format := vars["format"]

	genelist := strings.Split(genes, "+")

	commonGenes, err := mixt.GetCommonUserERGenes(module, tissue, genelist)
	if err != nil {
		fmt.Println("Could get common user enrichment genes")
		errorHandler(w, r, err)
		return
	}

	if format == "json" {
		common := Common{commonGenes}
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
	for _, gene := range commonGenes {
		entry := []string{gene}
		records = append(records, entry)
	}

	writer := csv.NewWriter(w)
	writer.WriteAll(records)

}
