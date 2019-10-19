package api

import (
	"encoding/json"
	"net/http"

	. ".."
	imagesfs "../../images/fs"
	"../db"
	"../fs"
	"../runtime"
	"github.com/gorilla/mux"
)

// ImportRequest is a struct for import request variables.
type ImportRequest struct {
	Path          string `json:"path"`
	ContainerName string `json:"containerName"`
}

// ExportRequest is a struct for export variables.
type ExportRequest struct {
	Path string `json:"path"`
}

// Adds router for container managment API
func AddRouter(router *mux.Router) {
	subrouter := router.PathPrefix("/containers").Subrouter()
	subrouter.HandleFunc("", listContainers).Methods("GET")
	subrouter.HandleFunc("/{name}", getInfo).Methods("GET")
	subrouter.HandleFunc("/create", createContainer).Methods("POST")
	subrouter.HandleFunc("/{name}", editContainer).Methods("PUT")
	subrouter.HandleFunc("/{name}", deleteContainer).Methods("DELETE")
	subrouter.HandleFunc("/{name}/start", startContainer).Methods("POST")
	subrouter.HandleFunc("/{name}/stop", stopContainer).Methods("POST")
	subrouter.HandleFunc("/{name}/restart", restartContainer).Methods("POST")
	subrouter.HandleFunc("/import", importContainer).Methods("POST")
	subrouter.HandleFunc("/{name}/export", exportCotainer).Methods("POST")
}

// Sends a list of all containers
func listContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := db.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	data, _ := json.Marshal(containers)
	w.Write(data)
}

// Sends a single container's info
func getInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["name"]
	containerInfo, err := db.GetInfo(containerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(containerInfo)
	w.Write(data)
}

// Creates a container and sends a response
func createContainer(w http.ResponseWriter, r *http.Request) {
	var containerInfo ContainerInfo
	err := json.NewDecoder(r.Body).Decode(&containerInfo)
	containerInfo.State = "stopped"
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.Create(containerInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	imagePath := imagesfs.GetImagePath(containerInfo.BaseImageName)
	err = fs.Create(containerInfo, imagePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Updates a container and send a response
func editContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["name"]
	var updateFields map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updateFields)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, updatedCount := db.Update(containerName, updateFields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if updatedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	newContainerName := updateFields["containerName"].(string)
	// Rename conatainer directory if necessary
	if newContainerName != containerName {
		fs.Rename(containerName, newContainerName)
	}
}

// Deletes a container and sends a response
func deleteContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["name"]

	err, deletedCount := db.Delete(containerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = fs.Delete(containerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Starts a container and sends a response
func startContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["name"]
	containerInfo, err := db.GetInfo(containerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = runtime.Start(containerInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Stop a container and sends a response
func stopContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["name"]
	err := runtime.Stop(containerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Restart a container and sends a response
func restartContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["name"]

	containerInfo, err := db.GetInfo(containerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = runtime.Restart(containerInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Imports a container from a tarball and sends a response
func importContainer(w http.ResponseWriter, r *http.Request) {
	var importRequest ImportRequest
	err := json.NewDecoder(r.Body).Decode(&importRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	containerInfo, err := fs.Import(importRequest.ContainerName, importRequest.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Create(containerInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(containerInfo)
	w.Write(data)
}

// Exports a container into a tarball and sends a response
func exportCotainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerName := vars["name"]
	containerInfo, err := db.GetInfo(containerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var exportRequest ExportRequest
	err = json.NewDecoder(r.Body).Decode(&exportRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = fs.Export(containerInfo, exportRequest.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
