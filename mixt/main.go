package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"code.google.com/p/gcfg"

	"bitbucket.org/vdumeaux/mixt/mixt/controllers"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var outsideTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/outside-navbar.html", "views/panels.html",
	"views/outside-index.html", "views/footer.html"))

var indexTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
	"views/index.html", "views/footer.html"))

var s = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if getUsername(r) == "" {
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

func startSession(username string, w http.ResponseWriter) error {
	val := map[string]string{"username": username}
	encoded, err := s.Encode("session", val)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cookie := &http.Cookie{
		Name:   "session",
		Value:  encoded,
		Path:   "/",
		MaxAge: 86400, // max age one day
	}
	http.SetCookie(w, cookie)
	return nil
}

func getUsername(r *http.Request) string {
	cookie, err := r.Cookie("session")

	if err != nil {
		fmt.Println(err)
		return ""
	}

	var val map[string]string

	err = s.Decode("session", cookie.Value, &val)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	username := val["username"]
	return username

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
			err := startSession(username, w)
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

	r.HandleFunc("/modules/{tissue}", controllers.ModulesHandler)
	r.HandleFunc("/modules/{tissue}/{modules}", controllers.ModuleHandler)

	r.HandleFunc("/public/{folder}/{file}", PublicHandler)

	r.HandleFunc("/search/{term}", controllers.SearchHandler)

	r.HandleFunc("/gene/{genes}", controllers.GeneHandler)
	r.HandleFunc("/gene/summary/{gene}", controllers.GeneSummaryHandler)

	r.HandleFunc("/tissues", controllers.TissuesHandler)
	r.HandleFunc("/tissues/{tissue1}/{tissue2}", controllers.TissueComparisonHandler)

	err = controllers.InitModules()
	if err != nil {
		log.Panic(err)
	}

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
