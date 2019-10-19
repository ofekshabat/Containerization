package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	containersapi "./containers/api"
	"./db"
	imagesapi "./images/api"
	packagesapi "./packages/api"

	"github.com/gorilla/mux"
)

const (
	serverPort  = 8000
	MecholaPath = "mechola"
)

var wg sync.WaitGroup

// Starts the server and listens for client requests
func Serve() {
	homePath := os.Getenv("HOME")
	os.Chdir(filepath.Join(homePath, MecholaPath))

	router := mux.NewRouter()
	containersapi.AddRouter(router)
	imagesapi.AddRouter(router)
	packagesapi.AddRouter(router)

	db.Connect()

	fmt.Printf("Listening on port %d... ", serverPort)
	address := fmt.Sprintf(":%d", serverPort)
	log.Fatal(http.ListenAndServe(address, router))
	fmt.Println()
}
