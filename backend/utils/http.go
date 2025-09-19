package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateSession(c *gin.Context) uint {
	userID, ok := GetUserID(c.Request.Context())
	if !ok {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return 0
	}
	return userID
}

func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value("user_id").(uint)
	return userID, ok
}

func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	c.Writer.Write(response)
}

func RespondWithError(c *gin.Context, code int, message string) {
	RespondWithJSON(c, code, map[string]string{"error": message})
}

func ParseJSONRequest(c *gin.Context, v interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	defer c.Request.Body.Close()
	return decoder.Decode(v)
}

func ParseJSONResponse(resp *http.Response, v interface{}) error {
	// Don't close here - let the caller handle closing
	return json.NewDecoder(resp.Body).Decode(v)
}

func CreateHTTPRequest(url string, method string, body interface{}, headers map[string]string) (*http.Request, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func SendHTTPRequest(req *http.Request, w http.ResponseWriter) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}
	// Don't defer close here - let the caller handle closing the response body

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Received non-success status code: %d", resp.StatusCode)
		// Return the response as-is, don't create a new one
		return resp, nil
	}

	return resp, nil
}

func AddQueryParams(req *http.Request, params map[string]string) {
	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
}
