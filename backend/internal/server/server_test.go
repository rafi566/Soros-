package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"soros/internal/api"
)

func TestHealthEndpoint(t *testing.T) {
	svc := NewAPIService()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	svc.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("unexpected health status: %s", body["status"])
	}
}

func TestCollectionsEndpoints(t *testing.T) {
	svc := NewAPIService()
	tests := []struct {
		path string
	}{{"/sources"}, {"/destinations"}, {"/connections"}}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		rec := httptest.NewRecorder()

		svc.Router().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s expected status 200, got %d", tt.path, rec.Code)
		}

		var payload []map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
			t.Fatalf("%s failed to parse response: %v", tt.path, err)
		}

		if len(payload) == 0 {
			t.Fatalf("%s returned empty payload", tt.path)
		}
	}
}

func TestJobLifecycle(t *testing.T) {
	svc := NewAPIService()

	req := httptest.NewRequest(http.MethodPost, "/jobs", nil)
	rec := httptest.NewRecorder()

	svc.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}

	var job api.Job
	if err := json.Unmarshal(rec.Body.Bytes(), &job); err != nil {
		t.Fatalf("failed to parse job response: %v", err)
	}

	if job.ID == "" || job.Status != "running" {
		t.Fatalf("unexpected job payload: %+v", job)
	}

	statusReq := httptest.NewRequest(http.MethodGet, "/jobs/"+job.ID, nil)
	statusRec := httptest.NewRecorder()
	svc.Router().ServeHTTP(statusRec, statusReq)

	if statusRec.Code != http.StatusOK {
		t.Fatalf("expected 200 when fetching job by id, got %d", statusRec.Code)
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		statusRec = httptest.NewRecorder()
		svc.Router().ServeHTTP(statusRec, statusReq)

		if statusRec.Code != http.StatusOK {
			t.Fatalf("status check returned %d", statusRec.Code)
		}

		var refreshed api.Job
		if err := json.Unmarshal(statusRec.Body.Bytes(), &refreshed); err != nil {
			t.Fatalf("failed to parse refreshed job: %v", err)
		}

		if refreshed.Status == "completed" {
			if refreshed.Progress != 100 || refreshed.FinishedAt.IsZero() {
				t.Fatalf("job did not finish cleanly: %+v", refreshed)
			}
			return
		}

		time.Sleep(1 * time.Second)
	}

	t.Fatalf("job did not complete in time")
}
