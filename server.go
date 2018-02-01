package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	Manager *Manager
	Port    int
}

func (server *Server) Listen() {
	// register handlers
	http.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		httpHandler(w, r, server.Manager)
	})
	// launch server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil))
}

func httpHandler(w http.ResponseWriter, r *http.Request, m *Manager) {
	var project struct{ Path string }

	switch r.Method {
	case "GET":
		projects := []struct{ Path string }{}
		for path := range m.projects {
			projects = append(projects, struct{ Path string }{Path: path})
		}
		json.NewEncoder(w).Encode(projects)
	case "POST":
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		log.Debug("Received POST for ", project.Path)
		m.Add(project.Path)
		w.WriteHeader(204)
	case "DELETE":
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		log.Debug("Received DELETE for ", project.Path)
		m.Remove(project.Path)
		w.WriteHeader(204)
	default:
		http.Error(w, "Request can not be processed", 400)
	}
}
