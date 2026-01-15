package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LambdaIITH/mess_registration/tests/testutils"
	"github.com/stretchr/testify/require"
)

func TestGetUserInfo(t *testing.T) {
	r := testutils.Router()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/getUser", nil)
	req.Header.Set("Authorization", testutils.TestJWT("student@iith.ac.in", 1, 0))

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
