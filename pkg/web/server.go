package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/borud/simpleton/pkg/store"
	"github.com/gorilla/mux"
)

// Server implements the system's webserver.
type Server struct {
	listenAddr    string
	staticDir     string
	maxUploadSize int64
	db            *store.SqliteStore
}

const (
	defaultOffset = 0
	defaultLimit  = 10
)

// New creates a new webserver instance
func New(store *store.SqliteStore, listen string, staticDir string) *Server {
	return &Server{
		listenAddr: listen,
		staticDir:  staticDir,
		db:         store,
	}
}

func (s *Server) dataHandler(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = defaultOffset
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = defaultLimit
	}

	dataArray, err := s.db.ListData(offset, limit)
	if err != nil {
		http.Error(w, "Error listing documents", http.StatusNotFound)
		log.Printf("Unable to list docs offset = %d, limit = %d: %v", offset, limit, err)
		return
	}

	json, err := json.MarshalIndent(dataArray, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting data to JSON: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

// dataPayloadHandler returns just the payload and sets the MIME type
// to application/octet-stream
func (s *Server) dataPayloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	data, err := s.db.Get(int64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Data %s not found", idString), http.StatusNotFound)
		log.Printf("Data '%s' not found: %v", idString, err)
		return
	}

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Write(data.Payload)
}

// ListenAndServe ...
func (s *Server) ListenAndServe() {
	m := mux.NewRouter().StrictSlash(true)

	// Data access
	m.HandleFunc("/data", s.dataHandler).Methods("GET")
	m.HandleFunc("/data/{id}", s.dataPayloadHandler).Methods("GET")

	// Serve static files
	m.PathPrefix("/").Handler(http.FileServer(http.Dir(s.staticDir)))

	server := &http.Server{
		Handler: m,
		Addr:    s.listenAddr,
	}

	log.Printf("Webserver listening to '%s'", s.listenAddr)
	log.Printf("Webserver terminated: '%v'", server.ListenAndServe())
}
