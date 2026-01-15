package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LambdaIITH/mess_registration/tests/testutils"
	"github.com/stretchr/testify/require"
)

func TestLoginEndpoint(t *testing.T) {
	r := testutils.Router()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/login", nil) // TODO: verify route

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
