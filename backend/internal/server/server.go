package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"soros/internal/api"
)

type APIService struct {
	mux          *http.ServeMux
	sources      []api.Source
	destinations []api.Destination
	connections  []api.Connection
	jobs         map[string]api.Job
	mu           sync.RWMutex
	jobCounter   atomic.Uint64
}

func NewAPIService() *APIService {
	svc := &APIService{
		mux: http.NewServeMux(),
		sources: []api.Source{
			{ID: "src-1", Name: "Postgres", Status: "ready"},
			{ID: "src-2", Name: "Stripe", Status: "ready"},
		},
		destinations: []api.Destination{
			{ID: "dst-1", Name: "BigQuery", Status: "ready"},
			{ID: "dst-2", Name: "Snowflake", Status: "ready"},
		},
		connections: []api.Connection{
			{ID: "con-1", SourceID: "src-1", DestinationID: "dst-1", Status: "scheduled"},
			{ID: "con-2", SourceID: "src-2", DestinationID: "dst-2", Status: "running"},
		},
		jobs: make(map[string]api.Job),
	}

	svc.mux.HandleFunc("/health", withTimeout(svc.handleHealth))
	svc.mux.HandleFunc("/sources", withTimeout(svc.handleSources))
	svc.mux.HandleFunc("/destinations", withTimeout(svc.handleDestinations))
	svc.mux.HandleFunc("/connections", withTimeout(svc.handleConnections))
	svc.mux.HandleFunc("/jobs", withTimeout(svc.handleJobs))
	svc.mux.HandleFunc("/jobs/", withTimeout(svc.handleJobByID))

	return svc
}

func (s *APIService) Router() http.Handler {
	return s.mux
}

func withTimeout(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()
		handler(w, r.WithContext(ctx))
	}
}

func (s *APIService) handleHealth(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, api.Health{Status: "ok"})
}

func (s *APIService) handleSources(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, s.sources)
}

func (s *APIService) handleDestinations(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, s.destinations)
}

func (s *APIService) handleConnections(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, s.connections)
}

func (s *APIService) handleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method { //nolint:exhaustive // only care about GET and POST
	case http.MethodGet:
		s.mu.RLock()
		jobs := make([]api.Job, 0, len(s.jobs))
		for _, job := range s.jobs {
			jobs = append(jobs, job)
		}
		s.mu.RUnlock()
		s.writeJSON(w, http.StatusOK, jobs)
	case http.MethodPost:
		job := s.startJob()
		s.writeJSON(w, http.StatusAccepted, job)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *APIService) handleJobByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/jobs/"):]
	if id == "" {
		http.NotFound(w, r)
		return
	}

	s.mu.RLock()
	job, ok := s.jobs[id]
	s.mu.RUnlock()
	if !ok {
		http.NotFound(w, r)
		return
	}

	s.writeJSON(w, http.StatusOK, job)
}

func (s *APIService) startJob() api.Job {
	id := fmt.Sprintf("job-%d", s.jobCounter.Add(1))
	now := time.Now().UTC()
	job := api.Job{ID: id, Status: "running", Progress: 0, StartedAt: now}

	s.mu.Lock()
	s.jobs[id] = job
	s.mu.Unlock()

	go s.runJob(id)

	return job
}

func (s *APIService) runJob(id string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for progress := 10; progress <= 100; progress += 15 {
		<-ticker.C

		s.mu.Lock()
		job := s.jobs[id]
		job.Progress = progress
		if progress == 100 {
			job.Status = "completed"
			job.FinishedAt = time.Now().UTC()
		}
		s.jobs[id] = job
		s.mu.Unlock()
	}
}

func (s *APIService) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
