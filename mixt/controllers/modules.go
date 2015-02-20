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

type Module struct {
	Name string
	URL  string
}

type Modules struct {
	Modules []Module
	Tissues []string
}

func loadModule(name string) Module {
	imgurl, err := mixt.Hist()
	if err != nil {
		return Module{"", ""}
	}
	return Module{Name: name, URL: imgurl}
}

func ModuleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	tissue := vars["tissue"]

	fmt.Println(name, tissue)

	module := loadModule(name)

	fmt.Println(module)
	moduleTemplate.Execute(w, module)
}

func ModulesHandler(w http.ResponseWriter, r *http.Request) {
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("Could not get modules")
		http.Error(w, err.Error(), 500)
	}

	m := Modules{nil, tissues}

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
