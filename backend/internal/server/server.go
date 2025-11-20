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
	fanouts      []api.Fanout
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
		fanouts: []api.Fanout{
			{ID: "fan-1", SourceID: "src-1", DestinationIDs: []string{"dst-1", "dst-2"}, Status: "ready"},
			{ID: "fan-2", SourceID: "src-2", DestinationIDs: []string{"dst-2"}, Status: "ready"},
		},
		jobs: make(map[string]api.Job),
	}

	svc.mux.HandleFunc("/health", withTimeout(svc.handleHealth))
	svc.mux.HandleFunc("/sources", withTimeout(svc.handleSources))
	svc.mux.HandleFunc("/destinations", withTimeout(svc.handleDestinations))
	svc.mux.HandleFunc("/connections", withTimeout(svc.handleConnections))
	svc.mux.HandleFunc("/fanouts", withTimeout(svc.handleFanouts))
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

func (s *APIService) handleFanouts(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, s.fanouts)
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
		var req struct {
			SourceID string   `json:"sourceId"`
			DestIDs  []string `json:"destinationIds"`
		}

		if r.Body != nil {
			defer r.Body.Close()
			_ = json.NewDecoder(r.Body).Decode(&req)
		}

		if req.SourceID == "" {
			req.SourceID = s.defaultSourceID()
		}

		if len(req.DestIDs) == 0 {
			req.DestIDs = s.defaultDestinations(req.SourceID)
		}

		if len(req.DestIDs) == 0 {
			s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no destinations provided"})
			return
		}

		job := s.startJob(req.SourceID, req.DestIDs)
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

func (s *APIService) startJob(sourceID string, destIDs []string) api.Job {
	id := fmt.Sprintf("job-%d", s.jobCounter.Add(1))
	now := time.Now().UTC()
	job := api.Job{ID: id, Status: "running", Progress: 0, StartedAt: now, SourceID: sourceID, DestIDs: destIDs}

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

func (s *APIService) defaultSourceID() string {
	if len(s.fanouts) > 0 {
		return s.fanouts[0].SourceID
	}

	if len(s.connections) > 0 {
		return s.connections[0].SourceID
	}

	return ""
}

func (s *APIService) defaultDestinations(sourceID string) []string {
	for _, fanout := range s.fanouts {
		if fanout.SourceID == sourceID {
			return append([]string(nil), fanout.DestinationIDs...)
		}
	}

	for _, conn := range s.connections {
		if conn.SourceID == sourceID {
			return []string{conn.DestinationID}
		}
	}

	if len(s.destinations) > 0 {
		return []string{s.destinations[0].ID}
	}

	return nil
}

func (s *APIService) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
