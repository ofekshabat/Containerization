package api

import (
	"encoding/json"
	"net/http"

	. ".."
	"../db"
	"../fs"
	"github.com/gorilla/mux"
)

// CreateRequest contains the information needed for creating images
type CreateRequest struct {
	Image         Image  `json:"image"`
	BaseImageName string `json:"baseImageName"`
}

// Adds router for image managment API
func AddRouter(router *mux.Router) {
	subrouter := router.PathPrefix("/images").Subrouter()
	subrouter.HandleFunc("", listImages).Methods("GET")
	subrouter.HandleFunc("/{name}", getImageInfo).Methods("GET")
	subrouter.HandleFunc("/create", createImage).Methods("POST")
	subrouter.HandleFunc("/{name}", deleteImage).Methods("DELETE")
}

// Sends a single image's info
func getImageInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["name"]
	imageInfo, err := db.GetInfo(imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(imageInfo)
	w.Write(data)
}

// Sends a list of all images
func listImages(w http.ResponseWriter, r *http.Request) {
	images, err := db.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	data, _ := json.Marshal(images)
	w.Write(data)
}

// Creates an image and sends a response
func createImage(w http.ResponseWriter, r *http.Request) {
	var createRequest CreateRequest
	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the image in the database
	err = db.Create(createRequest.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create the image in the filesystem
	err = fs.Create(createRequest.Image, createRequest.BaseImageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Deletes an image and sends a response
func deleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["name"]

	err, deletedCount := db.Delete(imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = fs.Delete(imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
