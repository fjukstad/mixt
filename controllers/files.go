package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
"bitbucket.org/vdumeaux/mixt/mixt"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	filetype := vars["filetype"]
	//name := vars["name"]
	body, err := mixt.Get(key, filetype)
	if err != nil {
		fmt.Println("Could not read all", err)
		return
	}
	w.Write(body)
}
