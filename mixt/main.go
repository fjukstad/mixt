package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"code.google.com/p/gcfg"

	"bitbucket.org/vdumeaux/mixt/mixt/controllers"
	"bitbucket.org/vdumeaux/mixt/mixt/mixt"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var outsideTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/outside-navbar.html",
	"views/outside-index.html", "views/footer.html"))

var indexTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/index.html", "views/footer.html"))

var s = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if controllers.GetUsername(r) == "" {
		outsideTemplate.Execute(w, nil)
		return
	}

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

type Config struct {
	Login struct {
		Username string
		Password string
	}
}

var cfg Config

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != "" && password != "" {
		if username == cfg.Login.Username && password == cfg.Login.Password {
			err := controllers.StartSession(username, w)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), 503)
				return
			}
			http.Redirect(w, r, "/", 302)
		}
	}
	e := Error{"Wrong username or password"}
	outsideTemplate.Execute(w, e)
}

func main() {

	err := gcfg.ReadFileInto(&cfg, "config.gcfg")
	if err != nil {
		fmt.Println(err)
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/login", LoginHandler)

	r.HandleFunc("/modules", controllers.ModulesHandler)
	r.HandleFunc("/modules/{tissue}/{modules}", controllers.ModuleHandler)

	r.HandleFunc("/public/{folder}/{file}", PublicHandler)

	r.HandleFunc("/search/{term}", controllers.SearchHandler)
	r.HandleFunc("/search/results/{terms}", controllers.SearchResultHandler)

	r.HandleFunc("/gene/summary/{gene}", controllers.GeneSummaryHandler)
	r.HandleFunc("/common/{tissue}/{module}/{geneset}/{status}", controllers.CommonGenesHandler)

	r.HandleFunc("/geneset/abstract/{geneset}", controllers.GeneSetAbstractHandler)

	r.HandleFunc("/tissues", controllers.TissuesHandler)
	r.HandleFunc("/tissues/{tissue1}/{tissue2}", controllers.TissueComparisonHandler)

	r.HandleFunc("/resources/{key}/{filetype}/{name}", controllers.FileHandler)

	addr := "opencpu:80"
	username := "biopsy@mcgill"
	password := "van-mi-ka-al"

	mixt.Init(addr, username, password)

	http.Handle("/", r)

	fmt.Println("Starting mixt app on port 8004")

	log.Fatal(http.ListenAndServe(":8004", nil))

}
