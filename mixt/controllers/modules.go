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

	modulesTemplate.Execute(w, nil)
}

type Modules struct {
	Names []string
}

func InitModules() error {
	err := mixt.Init()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}
