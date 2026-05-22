package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/internal/healthcheck"
	"github.com/cronwatch/internal/metrics"
)

func newStore(t *testing.T) *metrics.Store {
	t.Helper()
	return metrics.New()
}

func TestHandleHealth_Returns200(t *testing.T) {
	s := healthcheck.New(":0", newStore(t))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	s.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "ok" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestHandleStatus_EmptyStore(t *testing.T) {
	s := healthcheck.New(":0", newStore(t))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	s.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp healthcheck.StatusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("expected ok, got %q", resp.Status)
	}
}

func TestHandleStatus_DegradedOnFailure(t *testing.T) {
	store := newStore(t)
	store.Record("backup", false, 2*time.Second)
	s := healthcheck.New(":0", store)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	s.ServeHTTP(rec, req)
	var resp healthcheck.StatusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Status != "degraded" {
		t.Fatalf("expected degraded, got %q", resp.Status)
	}
}

func TestHandleStatus_ContainsJobMetrics(t *testing.T) {
	store := newStore(t)
	store.Record("cleanup", true, 500*time.Millisecond)
	s := healthcheck.New(":0", store)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	s.ServeHTTP(rec, req)
	var resp healthcheck.StatusResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if _, ok := resp.Jobs["cleanup"]; !ok {
		t.Fatal("expected cleanup job in response")
	}
}
