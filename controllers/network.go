package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"bitbucket.org/vdumeaux/mixt/mixt"
)

var networkTemplate = template.Must(template.ParseFiles("views/base.html",
	"views/header.html", "views/navbar.html",
	"views/network.html", "views/footer.html"))

func NetworkHandler(w http.ResponseWriter, r *http.Request) {
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("Could not get tissues")
		http.Error(w, err.Error(), 500)
		return
	}

	//remove nblood bnblood
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
