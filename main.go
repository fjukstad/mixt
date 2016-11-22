package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/fjukstad/mixt/controllers"
	"github.com/fjukstad/mixt/mixt"

	"github.com/gorilla/mux"
)

var outsideTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/outside-navbar.html",
	"views/outside-index.html", "views/footer.html"))

var indexTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/index.html", "views/footer.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, nil)
}

func PublicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["folder"]
	file := vars["file"]

	base := "public/"
	filename := base + folder + "/" + file

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		w.Write([]byte("Could not find file " + filename))
	} else {
		w.Write(f)
	}
}

type Error struct {
	Error string
}

func main() {

	port := flag.String("port", ":8004", "port to start the app on")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/modules", controllers.ModulesHandler)
	r.HandleFunc("/modules/{tissue}/{modules}/genes", controllers.GeneList)
	r.HandleFunc("/modules/{tissue}/{modules}/cohort/{cohort}", controllers.ModuleHandler)

	r.HandleFunc("/network", controllers.NetworkHandler)

	r.HandleFunc("/public/{folder}/{file}", PublicHandler)

	r.HandleFunc("/search/{term}", controllers.SearchHandler)
	r.HandleFunc("/search/results/{terms}", controllers.SearchResultHandler)

	r.HandleFunc("/gene/summary/{gene}", controllers.GeneSummaryHandler)
	r.HandleFunc("/common/{tissue}/{module}/{geneset}/{status}/{output}", controllers.CommonGenesHandler)

	r.HandleFunc("/geneset/abstract/{geneset}", controllers.GeneSetAbstractHandler)

	r.HandleFunc("/compare/{tissueA}/{tissueB}/{moduleA}/{moduleB}/cohort/{cohort}", controllers.CompareModulesHandler)

	r.HandleFunc("/clinical-comparison", controllers.ModuleClinicalHandler)
	r.HandleFunc("/clinical-comparison/{tissue}/{analysis}/{cohort}", controllers.ModuleClinicalAnalysisHandler)

	r.HandleFunc("/tissues", controllers.TissuesHandler)
	r.HandleFunc("/tissues/{tissueA}/{tissueB}", controllers.TissueComparisonHandler)
	r.HandleFunc("/tissues/{tissueA}/{tissueB}/{analysis}/{cohort}", controllers.AnalysisHandler)

	r.HandleFunc("/resources/{key}/{filetype}/{name}", controllers.FileHandler)

	r.HandleFunc("/userlist", controllers.UserListHandler)
	r.HandleFunc("/userlist/submit", controllers.UserListSubmitHandler)
	r.HandleFunc("/userlist/result/{listname}/{genes}", controllers.UserListResultHandler)
	r.HandleFunc("/userlist/common/{tissue}/{module}/{genes}/{format}", controllers.UserListCommonGenesHandler)

	r.HandleFunc("/common-go/{tissue}/{module}/{id}/{format}", controllers.CommonGOTermGenesHandler)

	r.HandleFunc("/suggest-papers/{tissue}/{module}", controllers.PaperSuggestionHandler)

	r.HandleFunc("/tomgraph/{tissue}/{component}/{format}", controllers.TOMGraphHandler)

	addr := "mixt-blood-tumor.bci.mcgill.ca:8787"
	username := ""
	password := ""

	mixt.Init(addr, username, password)

	http.Handle("/", r)

	fmt.Println("Starting mixt app on", *port)
	log.Fatal(http.ListenAndServe(*port, r))

}
