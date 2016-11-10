package controllers

import (
	"net/http"

	"github.com/fjukstad/mixt-blood-tumor/mixt"
	"github.com/gorilla/mux"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	filetype := vars["filetype"]
	//name := vars["name"]
	body, err := mixt.Get(key, filetype)
	if err != nil {
		errorHandler(w, r, err)
		return
	}
	w.Write(body)
}
