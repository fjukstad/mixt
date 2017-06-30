package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/fjukstad/mixt/mixt"
	"github.com/gorilla/mux"
)

var networkTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html", "views/loader.html",
	"views/network.html", "views/footer.html"))

func NetworkHandler(w http.ResponseWriter, r *http.Request) {
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("Could not get tissues")
		errorHandler(w, r, err)
		return
	}

	var t []string
	for _, tissue := range tissues {
		if tissue != "bnblood" && tissue != "nblood" {
			t = append(t, tissue)
		}
	}

	data := struct {
		Tissues []string
	}{
		t,
	}

	networkTemplate.Execute(w, data)
}

func TOMGraphHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tissue := vars["tissue"]
	component := vars["component"]
	format := vars["format"]

	res, err := mixt.GetTOMGraph(tissue, component, format)
	if err != nil {
		errorHandler(w, r, err)
		return
	}

	// If user wants a csv set appropriate filename and extensions.
	if format == "csv" {
		filename := tissue + "-network.csv"
		w.Header().Set("Trailer", "Content-Disposition")
		w.Header().Add("Content-Disposition", "attachment; filename=\""+filename+"\"")
	}

	w.Write(res)

}
