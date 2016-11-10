package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"

	"github.com/fjukstad/mixt-blood-tumor/mixt"
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
			fmt.Println("Could not get module", err)
			errorHandler(w, r, err)
			return
		}
		modules = append(modules, m)
	}

	overview, err := retrieveModulesOverview()
	if err != nil {
		fmt.Println("Modules overview error", err)
		errorHandler(w, r, err)
		return
	}

	err = moduleTemplate.Execute(w, Modules{modules, overview})
	if err != nil {
		fmt.Println(err)
	}

}

func ModulesHandler(w http.ResponseWriter, r *http.Request) {
	m, err := retrieveModulesOverview()
	if err != nil {
		fmt.Println("Could not get modules")
		errorHandler(w, r, err)
		return
	}

	var cleanTissues []string
	for _, i := range m.Tissues {
		if i != "nblood" && i != "bnblood" {
			cleanTissues = append(cleanTissues, i)
		}
	}
	m.Tissues = cleanTissues

	err = modulesTemplate.Execute(w, m)
	if err != nil {
		fmt.Println(err)
	}

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
		fmt.Println("Could not get cohorts", err)
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
		errorHandler(w, r, err)
		return
	}

	ma, err := mixt.GetModule(moduleA, tissueA, cohort)
	if err != nil {
		fmt.Println("Could not get module", moduleA)
		errorHandler(w, r, err)
		return
	}

	mb, err := mixt.GetModule(moduleB, tissueB, cohort)
	if err != nil {
		fmt.Println("Could not get module", moduleB)
		errorHandler(w, r, err)
		return
	}

	reorderHeatmap, err := mixt.HeatmapReOrder(tissueB, moduleB, tissueA, moduleA, cohort)
	if err != nil {
		fmt.Println("Could not get re order heatmap")
		errorHandler(w, r, err)
		return
	}

	mb.HeatmapUrl = reorderHeatmap

	scatterPlot, err := mixt.CohortScatterplot(tissueA, tissueB, moduleA, moduleB, cohort)
	if err != nil {
		fmt.Println("Could not get cohort scatterplot")
		errorHandler(w, r, err)
		return
	}

	ma.ScatterplotUrl = scatterPlot

	scatterPlot, err = mixt.CohortScatterplot(tissueB, tissueA, moduleB, moduleA, cohort)
	if err != nil {
		fmt.Println("Could not get cohort scatterplot", err)
		errorHandler(w, r, err)
		return
	}
	mb.ScatterplotUrl = scatterPlot

	ma.AlternativeHeatmapUrl, err = mixt.HeatmapReOrder(tissueA, moduleA, tissueB, moduleB, cohort)
	if err != nil {
		fmt.Println("Could not get re order heatmap")
		errorHandler(w, r, err)
		return
	}

	// Get alternative heatmaps for both modules
	mb.AlternativeHeatmapUrl, err = mixt.CohortHeatmap(tissueB, moduleB, cohort)
	if err != nil {
		fmt.Println("Could not get heatmap", tissueB, moduleB)
		errorHandler(w, r, err)
		return
	}

	modules := []mixt.Module{ma, mb}

	err = compareModulesTemplate.Execute(w, ModuleComparison{modules, analyses})
	if err != nil {
		fmt.Println(err)
	}
}

type Clinical struct {
	Tissues []string
	Cohorts []string
}

func ModuleClinicalHandler(w http.ResponseWriter, r *http.Request) {
	tissues, err := mixt.GetTissues()

	if err != nil {
		fmt.Println("Could not get tissues")
		errorHandler(w, r, err)
		return
	}

	cohorts, err := mixt.GetCohorts()
	if err != nil {
		fmt.Println("Could not get cohorts", err)
		errorHandler(w, r, err)
		return
	}

	c := Clinical{tissues, cohorts}
	err = clinicalComparisonTemplate.Execute(w, c)
	if err != nil {
		fmt.Println(err)
	}

}

func ModuleClinicalAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tissue := vars["tissue"]
	analysis := vars["analysis"]
	cohort := vars["cohort"]

	if cohort == "" {
		cohort = "all"
	}

	var res []byte
	var err error

	if analysis == "eigengene" {
		res, err = mixt.ClinicalEigengene(tissue)
	} else if analysis == "eigengene" {
		res, err = mixt.ClinicalROI(tissue)
	} else if analysis == "ranksum" {
		res, err = mixt.ClinicalRanksum(tissue, cohort)
	}

	if err != nil {
		fmt.Println("Could not run analysis", tissue, analysis)
		errorHandler(w, r, err)
		return
	}

	w.Write(res)
}
