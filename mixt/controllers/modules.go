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
	"views/header.html", "views/navbar-module.html",
	"views/module.html", "views/footer.html"))

var modulesTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/modules.html", "views/footer.html"))

var compareModulesTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/module-comparison.html", "views/footer.html"))

type ModulesOverview struct {
	Modules map[string][]mixt.Module
	Tissues []string
}

type Modules struct {
	Modules  []mixt.Module
	Overview ModulesOverview
}

func ModuleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mods := vars["modules"]
	tissue := vars["tissue"]

	moduleString := strings.Split(mods, "+")

	var modules []mixt.Module

	for _, module := range moduleString {
		m, err := mixt.GetModule(module, tissue)
		if err != nil {
			fmt.Println("Could not get module")
			http.Error(w, err.Error(), 503)
			return
		}
		modules = append(modules, m)
	}

	overview, err := retrieveModulesOverview()
	if err != nil {
		http.Error(w, err.Error(), 503)
	}

	moduleTemplate.Execute(w, Modules{modules, overview})
}

func ModulesHandler(w http.ResponseWriter, r *http.Request) {
	m, err := retrieveModulesOverview()
	if err != nil {
		fmt.Println("Could not get modules")
		http.Error(w, err.Error(), 503)
	}
	modulesTemplate.Execute(w, m)
}

func retrieveModulesOverview() (ModulesOverview, error) {
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("Could not get tissues")
		return ModulesOverview{}, err
	}

	modules := make(map[string][]mixt.Module)

	for _, tissue := range tissues {
		m, err := mixt.GetModules(tissue)
		if err != nil {
			return ModulesOverview{}, err
		}
		modules[tissue] = m
	}

	m := ModulesOverview{modules, tissues}
	return m, nil
}

type ModuleComparison struct {
	Modules  []mixt.Module
	Analyses mixt.Analyses
}

func CompareModulesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleA := vars["moduleA"]
	moduleB := vars["moduleB"]
	tissueA := vars["tissueA"]
	tissueB := vars["tissueB"]

	analyses, err := mixt.ModuleComparisonAnalyses(tissueA, tissueB, moduleA, moduleB)
	if err != nil {
		fmt.Println("Could not get module comparison analyses")
		http.Error(w, err.Error(), 503)
		return
	}

	ma, err := mixt.GetModule(moduleA, tissueA)
	if err != nil {
		fmt.Println("Could not get module", moduleA)
		http.Error(w, err.Error(), 503)
		return
	}

	mb, err := mixt.GetModule(moduleB, tissueB)
	if err != nil {
		fmt.Println("Could not get module", moduleB)
		http.Error(w, err.Error(), 503)
		return
	}

	reorderHeatmap, err := mixt.HeatmapReOrder(tissueB, moduleB, tissueA, moduleA)
	if err != nil {
		fmt.Println("Could not get re order heatmap")
		http.Error(w, err.Error(), 503)
		return
	}

	mb.HeatmapUrl = reorderHeatmap

	modules := []mixt.Module{ma, mb}

	compareModulesTemplate.Execute(w, ModuleComparison{modules, analyses})
}
