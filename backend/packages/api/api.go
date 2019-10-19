package api

import (
	"encoding/json"
	"net/http"

	. ".."
	"../db"
	"github.com/gorilla/mux"
)

// Adds router for package managment API
func AddRouter(router *mux.Router) {
	subrouter := router.PathPrefix("/packages").Subrouter()
	subrouter.HandleFunc("", listPackages).Methods("GET")
	subrouter.HandleFunc("/{name}", getPackageInfo).Methods("GET")
	subrouter.HandleFunc("/create", createPackage).Methods("POST")
	subrouter.HandleFunc("/{name}", editPackage).Methods("PUT")
	subrouter.HandleFunc("/{name}", deletePackage).Methods("DELETE")
}

// Sends a list of all packages
func listPackages(w http.ResponseWriter, r *http.Request) {
	packages, err := db.ListPackages()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(packages)
	w.Write(data)
}

// Sends a single package's info
func getPackageInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageName := vars["name"]
	packageInfo, err := db.GetPackageInfo(packageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(packageInfo)
	w.Write(data)
}

// Creates a package and sends a response
func createPackage(w http.ResponseWriter, r *http.Request) {
	var packageInfo Package
	err := json.NewDecoder(r.Body).Decode(&packageInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.CreatePackage(packageInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Updates a packae and sends a response
func editPackage(w http.ResponseWriter, r *http.Request) {

}

// Deletes a package and sends a response
func deletePackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageName := vars["name"]
	err := db.DeletePackage(packageName)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
