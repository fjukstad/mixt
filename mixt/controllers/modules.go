package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"bitbucket.org/vdumeaux/mixt/mixt/models"

	"github.com/fjukstad/kvik/dataset"
	"github.com/gorilla/mux"
)

// Dataset where do our analyses
var d *dataset.Dataset

var moduleTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
	"views/module.html", "views/footer.html"))
var moduleTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/panels.html",
	"views/moduleshtml", "views/footer.html"))

type Module struct {
	Name string
	URL  string
}

func loadModule(name string) Module {
	imgurl, err := models.Hist(d)
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

type Modules struct {
	Names []string
}

func ListModules(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tissue := vars["tissue"]

	fmt.Println("Listing modules for", tissue)

	m := Modules{Names: []string{"module a", "module b", "module c"}}
	modules, _ := json.Marshal(m)

	res, _ := models.Add(d, 2, 3)
	log.Println(res)

	w.Write(modules)
}

func InitModules() error {
	var err error
	d, err = models.Init()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}
