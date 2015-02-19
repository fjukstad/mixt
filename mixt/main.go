package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"bitbucket.org/vdumeaux/mixt/mixt/controllers"

	"github.com/gorilla/mux"
)

var indexTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
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

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/modules/{tissue}", controllers.ListModules)
	r.HandleFunc("/module/{tissue}/{name}", controllers.ModuleHandler)

	r.HandleFunc("/public/{folder}/{file}", PublicHandler)

	err := controllers.InitModules()
	if err != nil {
		log.Panic(err)
	}

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
