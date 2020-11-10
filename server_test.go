package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETSubscriptions(t *testing.T) {
	t.Run("Returns 200 OK", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		SubscriptionServer(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}