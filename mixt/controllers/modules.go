package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"

	"github.com/gorilla/mux"
)

var moduleTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
	"views/module.html", "views/footer.html"))

var modulesTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
	"views/modules.html", "views/footer.html"))

type Modules struct {
	Modules []string
	Tissues []string
}

func ModuleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	tissue := vars["tissue"]
	fmt.Println(name, tissue)

	moduleTemplate.Execute(w, nil)
}

func ModulesHandler(w http.ResponseWriter, r *http.Request) {

	modules, err := mixt.GetModules("blood")
	if err != nil {
		fmt.Println("Could not get modules")
		http.Error(w, err.Error(), 503)
		return
	}

	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("Could not get tissues")
		http.Error(w, err.Error(), 500)
		return
	}

	m := Modules{modules, tissues}
	modulesTemplate.Execute(w, m)
}

func InitModules() error {
	err := mixt.Init()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}
