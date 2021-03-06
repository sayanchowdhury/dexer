package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/farhaanbukhsh/file-indexer/lib/conf"
	"github.com/farhaanbukhsh/file-indexer/lib/indexer"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var templates []string

func init() {
	templates = []string{
		"ui/index.html",
		"ui/layout/header.html",
		"ui/layout/footer.html",
		"ui/layout/search.html",
	}
}

// Server Data Structure for holding the configuration and logger
type Server struct {
	conf conf.Configuration
}

// RootHandler is the controller responsible for the frontend
func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(templates...)
	s.conf.LogMan.Must(err, "Template Parsed Successfully")
	t.ExecuteTemplate(w, "index", nil)
}

// SearchFile is the controller that helps with indexing the file
func (s *Server) SearchFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchResult := indexer.Search(s.conf.IndexFilename, params["query"])
	json.NewEncoder(w).Encode(searchResult.Hits)
}

// Start function starts the server
func (s *Server) Start() {
	s.conf.LogMan.Must(nil, "Serving on "+s.conf.Port)
	router := mux.NewRouter()
	router.HandleFunc("/", s.RootHandler)
	router.HandleFunc("/search/{query}", s.SearchFile).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/"))))
	s.conf.LogMan.Must(http.ListenAndServe(s.conf.Port, handlers.LoggingHandler(s.conf.LogMan, router)), "")
}

// NewServer function creates a new server and return a pointer to it
func NewServer(c conf.Configuration) *Server {
	temp := Server{
		conf: c,
	}
	return &temp
}
