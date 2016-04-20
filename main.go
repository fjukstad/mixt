package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/scalingdata/gcfg"

	"bitbucket.org/vdumeaux/mixt/controllers"
	"bitbucket.org/vdumeaux/mixt/mixt"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"

	"github.com/codegangsta/negroni"
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

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// any request to /public is open to anyone
	if strings.HasPrefix(r.URL.Path, "/public") {
		next(w, r)
		return
	}

	// don't have to be logged in to log in...
	if strings.HasPrefix(r.URL.Path, "/login") {
		next(w, r)
		return
	}

	// auth user
	if controllers.LoggedIn(r) {
		next(w, r)
		return
	} else {
		e := Error{"Please log in."}
		outsideTemplate.Execute(w, e)
	}
}

func main() {

	err := gcfg.ReadFileInto(&cfg, "config.gcfg")
	if err != nil {
		fmt.Println(err)
		return
	}
	port := flag.String("port", ":8004", "port to start the app on")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/login", LoginHandler)

	r.HandleFunc("/modules", controllers.ModulesHandler)
	r.HandleFunc("/modules/{tissue}/{modules}", controllers.ModuleHandler)
	
	r.HandleFunc("/network", controllers.NetworkHandler)

	r.HandleFunc("/public/{folder}/{file}", PublicHandler)

	r.HandleFunc("/search/{term}", controllers.SearchHandler)
	r.HandleFunc("/search/results/{terms}", controllers.SearchResultHandler)

	r.HandleFunc("/gene/summary/{gene}", controllers.GeneSummaryHandler)
	r.HandleFunc("/common/{tissue}/{module}/{geneset}/{status}/{output}", controllers.CommonGenesHandler)

	r.HandleFunc("/geneset/abstract/{geneset}", controllers.GeneSetAbstractHandler)

	r.HandleFunc("/compare/{tissueA}/{tissueB}/{moduleA}/{moduleB}", controllers.CompareModulesHandler)

	r.HandleFunc("/clinical-comparison", controllers.ModuleClinicalHandler)
	r.HandleFunc("/clinical-comparison/{tissue}/{analysis}", controllers.ModuleClinicalAnalysisHandler)

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

	r.HandleFunc("/tomgraph/{tissue}/{what}", controllers.TOMGraphHandler)

	addr := "gor:8181"
	username := "biopsy@mcgill"
	password := "van-mi-ka-al"

	mixt.Init(addr, username, password)

	http.Handle("/", r)

	n := negroni.New(negroni.HandlerFunc(AuthMiddleware))
	n.UseHandler(r)

	fmt.Println("Starting mixt app on", *port)
	log.Fatal(http.ListenAndServe(*port, n))

}
