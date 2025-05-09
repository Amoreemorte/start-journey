package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	Router *mux.Router
	Db     Dbdocuments
}

func (api *API) RegisterEndpoints() {
	api.Router.Use(logMiddleware)
	api.Router.Use(authMiddleware)
	api.Router.HandleFunc("/api/v1/docs/delete/{id}", api.deleteDocumentHandler).Methods(http.MethodDelete)
	api.Router.HandleFunc("/api/v1/docs/search/{id}", api.getDocumentHandler).Methods(http.MethodGet)
	api.Router.HandleFunc("/api/v1/docs/create/{author}", api.createDocumentHandler).Methods(http.MethodPost)
	api.Router.HandleFunc("/api/v1/docs/update/{id}", api.updateDocumentHandler).Methods(http.MethodPatch)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (api *API) getDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	document, err := api.Db.SearchDocumentById(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
	} else {
		json.NewEncoder(w).Encode(document)
	}

	log.Println("Search id is:", id)
}

func (api *API) createDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	document, err := api.Db.CreateDocument(name, "storage")
	if err != nil {
		http.Error(w, err.Error(), 400)
	} else {
		json.NewEncoder(w).Encode(document)
	}
	log.Println("Create id is:", document.Id)
}

func (api *API) deleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := api.Db.DeleteDocumentById(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
	} else {
		w.Write([]byte(fmt.Sprintf("Delete id is: %v", id)))
	}
}

func (api *API) updateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	document, err := api.Db.SearchDocumentById(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
	} else {
		json.NewEncoder(w).Encode(document)
	}
	log.Println("Update id is:", id)
}

func main() {
	api := API{Router: mux.NewRouter(), Db: InitDB("C:/learning_go/api/apiTask/documents.csv")}
	api.RegisterEndpoints()

	http.ListenAndServe(":8080", api.Router)
}
