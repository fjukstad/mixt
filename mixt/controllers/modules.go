package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"

	"github.com/gorilla/mux"
)

var moduleTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/module.html", "views/footer.html"))

var modulesTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
	"views/modules.html", "views/footer.html"))

type ModulesOverview struct {
	Modules []mixt.Module
	Tissues []string
}

type Modules struct {
	Modules []mixt.Module
}

func ModuleHandler(w http.ResponseWriter, r *http.Request) {
	if !LoggedIn(r) {
		http.Redirect(w, r, "/", 302)
		return
	}

	vars := mux.Vars(r)
	mods := vars["modules"]
	tissue := vars["tissue"]

	modules := strings.Split(mods, "+")

	var result []mixt.Module

	for _, module := range modules {
		m, err := mixt.GetModule(module, tissue)
		if err != nil {
			fmt.Println("Could not get module")
			http.Error(w, err.Error(), 503)
			return
		}
		result = append(result, m)
	}

	moduleTemplate.Execute(w, Modules{result})
}

func ModulesHandler(w http.ResponseWriter, r *http.Request) {
	if !LoggedIn(r) {
		http.Redirect(w, r, "/", 302)
		return
	}

	vars := mux.Vars(r)
	tissue := vars["tissue"]

	modules, err := mixt.GetModules(tissue)
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

	m := ModulesOverview{modules, tissues}
	modulesTemplate.Execute(w, m)
}
