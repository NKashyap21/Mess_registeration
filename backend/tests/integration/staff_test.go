package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LambdaIITH/mess_registration/tests/testutils"
	"github.com/stretchr/testify/require"
)

func TestStaffOnlyEndpointForbiddenForStudent(t *testing.T) {
	r := testutils.Router()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/messStaff/scanning", nil)
	req.Header.Set("Authorization", testutils.TestJWT("student@iith.ac.in", 1, 0))

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestStaffOnlyEndpointAllowedForStaff(t *testing.T) {
	r := testutils.Router()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/messStaff/scanning", nil)
	req.Header.Set("Authorization", testutils.TestJWT("staff@iith.ac.in", 2, 1))

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
