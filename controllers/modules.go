package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"

	"bitbucket.org/vdumeaux/mixt/mixt"
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

var clinicalComparisonTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/module-clinical-comparison.html", "views/footer.html"))

type ModulesOverview struct {
	Modules map[string][]mixt.Module
	Tissues []string
	Cohorts []string
}

type Modules struct {
	Modules  []mixt.Module
	Overview ModulesOverview
}

func ModuleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mods := vars["modules"]
	tissue := vars["tissue"]
	cohort := vars["cohort"]

	moduleString := strings.Split(mods, "+")

	fmt.Println("Selected cohort", cohort)

	var modules []mixt.Module

	for _, module := range moduleString {
		m, err := mixt.GetModule(module, tissue, cohort)
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

	var cleanTissues []string
	for _, i := range m.Tissues {
		if i != "nblood" && i != "bnblood" {
			cleanTissues = append(cleanTissues, i)
		}
	}
	m.Tissues = cleanTissues

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

	cohorts, err := mixt.GetCohorts()
	if err != nil {
		fmt.Println("Could not get cohorts")
		return ModulesOverview{}, err
	}

	m := ModulesOverview{modules, tissues, cohorts}
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
	cohort := vars["cohort"]

	analyses, err := mixt.ModuleComparisonAnalyses(tissueA, tissueB, moduleA, moduleB)
	if err != nil {
		fmt.Println("Could not get module comparison analyses")
		http.Error(w, err.Error(), 503)
		return
	}

	ma, err := mixt.GetModule(moduleA, tissueA, cohort)
	if err != nil {
		fmt.Println("Could not get module", moduleA)
		http.Error(w, err.Error(), 503)
		return
	}

	mb, err := mixt.GetModule(moduleB, tissueB, cohort)
	if err != nil {
		fmt.Println("Could not get module", moduleB)
		http.Error(w, err.Error(), 503)
		return
	}

	reorderHeatmap, err := mixt.HeatmapReOrder(tissueB, moduleB, tissueA, moduleA, cohort)
	if err != nil {
		fmt.Println("Could not get re order heatmap")
		http.Error(w, err.Error(), 503)
		return
	}

	mb.HeatmapUrl = reorderHeatmap

	scatterPlot, err := mixt.CohortScatterplot(tissueA, tissueB, moduleA, moduleB, cohort)
	if err != nil {
		fmt.Println("Could not get cohort scatterplot")
		http.Error(w, err.Error(), 503)
		return
	}

	ma.ScatterplotUrl = scatterPlot

	scatterPlot, err = mixt.CohortScatterplot(tissueB, tissueA, moduleB, moduleA, cohort)
	if err != nil {
		fmt.Println("Could not get cohort scatterplot", err)
		http.Error(w, err.Error(), 503)
		return
	}
	mb.ScatterplotUrl = scatterPlot

	ma.AlternativeHeatmapUrl, err = mixt.HeatmapReOrder(tissueA, moduleA, tissueB, moduleB, cohort)
	if err != nil {
		fmt.Println("Could not get re order heatmap")
		http.Error(w, err.Error(), 503)
		return
	}

	// Get alternative heatmaps for both modules
	mb.AlternativeHeatmapUrl, err = mixt.CohortHeatmap(tissueB, moduleB, cohort)
	if err != nil {
		fmt.Println("Could not get heatmap", tissueB, moduleB)
		http.Error(w, err.Error(), 503)
		return
	}

	modules := []mixt.Module{ma, mb}

	compareModulesTemplate.Execute(w, ModuleComparison{modules, analyses})
}

type Clinical struct {
	Tissues []string
}

func ModuleClinicalHandler(w http.ResponseWriter, r *http.Request) {
	tissues, err := mixt.GetTissues()

	if err != nil {
		fmt.Println("Could not get tissues")
		http.Error(w, err.Error(), 503)
		return
	}

	c := Clinical{tissues}
	clinicalComparisonTemplate.Execute(w, c)
}

func ModuleClinicalAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tissue := vars["tissue"]
	analysis := vars["analysis"]
	//cohort := vars["cohort"]

	var res []byte
	var err error

	if analysis == "eigengene" {
		res, err = mixt.ClinicalEigengene(tissue)
	} else if analysis == "eigengene" {
		res, err = mixt.ClinicalROI(tissue)
	} else if analysis == "ranksum" {
		res, err = mixt.ClinicalRanksum(tissue)
	}

	fmt.Println("Clinical results", string(res))

	if err != nil {
		fmt.Println("Could not run analysis", tissue, analysis)
		http.Error(w, err.Error(), 503)
	}

	w.Write(res)
}
