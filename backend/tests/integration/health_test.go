package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LambdaIITH/mess_registration/router"
	"github.com/stretchr/testify/assert"
)

type Server struct {

}

func

func TestHealthEndpoint(t *testing.T) {
	r := router.SetupRouter()

	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")
	assert.Contains(t, w.Body.String(), "healthy", "Expected response body to contain 'healthy'")
}
