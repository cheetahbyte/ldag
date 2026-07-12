package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNameFromHost(t *testing.T) {
	tests := map[string]string{
		"marvin.leg.dich.aufs.gleis.rip":      "marvin",
		"marvin.leg.dich.aufs.gleis.rip:8080": "marvin",
		"gleis.rip":                           "Jemand",
		"localhost:8080":                      "Jemand",
		"127.0.0.1:8080":                      "Jemand",
		"-marvin.leg.dich.aufs.gleis.rip":     "Jemand",
	}

	for host, want := range tests {
		t.Run(host, func(t *testing.T) {
			if got := nameFromHost(host); got != want {
				t.Fatalf("nameFromHost(%q) = %q, want %q", host, got, want)
			}
		})
	}
}

func TestIndexHandlerRendersHostName(t *testing.T) {
	req := httptest.NewRequest("GET", "http://marvin.leg.dich.aufs.gleis.rip/", nil)
	rec := httptest.NewRecorder()

	indexHandler(rec, req)
	body := rec.Body.String()

	if rec.Code != 200 {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(body, "<h1>Marvin</h1>") {
		t.Fatalf("body does not contain rendered name: %s", body)
	}
	if !strings.Contains(body, "|=|=|=|=|=|=|=|=|=|") {
		t.Fatalf("body does not contain track art: %s", body)
	}
}
