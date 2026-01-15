package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LambdaIITH/mess_registration/tests/testutils"
	"github.com/stretchr/testify/require"
)

func TestCreateSwap(t *testing.T) {
	r := testutils.Router()
	w := httptest.NewRecorder()
	body := strings.NewReader(`{"type":"veg","direction":"up"}`) // TODO: confirm payload
	req := httptest.NewRequest("POST", "/api/students/createSwap", body)
	req.Header.Set("Authorization", testutils.TestJWT("student@iith.ac.in", 1, 0))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
