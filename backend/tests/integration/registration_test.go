package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LambdaIITH/mess_registration/tests/testutils"
	"github.com/stretchr/testify/require"
)

func TestRegisterMess(t *testing.T) {
	r := testutils.Router()
	w := httptest.NewRecorder()
	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("POST", "/api/students/registerMess/1", body)
	req.Header.Set("Authorization", testutils.TestJWT("student@iith.ac.in", 1, 0))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
