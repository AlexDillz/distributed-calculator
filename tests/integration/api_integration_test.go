package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexDillz/distributed-calculator/internal/server"
	"github.com/AlexDillz/distributed-calculator/internal/storage"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *httptest.Server {
	store, _ := storage.NewStorage(":memory:")
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", server.RegisterHandler(store))
	mux.HandleFunc("/api/v1/login", server.LoginHandler(store))
	mux.HandleFunc("/api/v1/calculate",
		server.AuthMiddleware(server.CalculateHandler(store)),
	)
	mux.HandleFunc("/api/v1/expressions",
		server.AuthMiddleware(server.ListExpressionsHandler(store)),
	)
	mux.HandleFunc("/api/v1/expressions/",
		server.AuthMiddleware(server.GetExpressionHandler(store)),
	)
	return httptest.NewServer(mux)
}

func TestRegisterLoginCalculate(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/v1/register", "application/json",
		strings.NewReader(`{"login":"u","password":"p"}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	buf := bytes.NewBufferString(`{"login":"u","password":"p"}`)
	resp, err = http.Post(srv.URL+"/api/v1/login", "application/json", buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var payload struct{ Token string }
	err = json.NewDecoder(resp.Body).Decode(&payload)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload.Token)

	reqBody := `{"expression":"2+2*2"}`
	req, _ := http.NewRequest("POST", srv.URL+"/api/v1/calculate", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+payload.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]float64
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, 6.0, result["result"])
}

func TestListExpressions(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	http.Post(srv.URL+"/api/v1/register", "application/json",
		strings.NewReader(`{"login":"u","password":"p"}`))
	loginResp, _ := http.Post(srv.URL+"/api/v1/login", "application/json",
		strings.NewReader(`{"login":"u","password":"p"}`))
	var payload struct{ Token string }
	json.NewDecoder(loginResp.Body).Decode(&payload)

	for _, expr := range []string{"1+1", "3*3"} {
		req, _ := http.NewRequest("POST", srv.URL+"/api/v1/calculate",
			strings.NewReader(fmt.Sprintf(`{"expression":"%s"}`, expr)))
		req.Header.Set("Authorization", "Bearer "+payload.Token)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	req, _ := http.NewRequest("GET", srv.URL+"/api/v1/expressions", nil)
	req.Header.Set("Authorization", "Bearer "+payload.Token)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResp struct {
		Expressions []struct {
			ID     int     `json:"id"`
			Status string  `json:"status"`
			Result float64 `json:"result"`
		} `json:"expressions"`
	}
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	assert.NoError(t, err)
	assert.Len(t, listResp.Expressions, 2)
}

func TestGetExpressionByID(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	http.Post(srv.URL+"/api/v1/register", "application/json",
		strings.NewReader(`{"login":"u","password":"p"}`))
	loginResp, _ := http.Post(srv.URL+"/api/v1/login", "application/json",
		strings.NewReader(`{"login":"u","password":"p"}`))
	var payload struct{ Token string }
	json.NewDecoder(loginResp.Body).Decode(&payload)

	reqCalc, _ := http.NewRequest("POST", srv.URL+"/api/v1/calculate",
		strings.NewReader(`{"expression":"5-2"}`))
	reqCalc.Header.Set("Authorization", "Bearer "+payload.Token)
	reqCalc.Header.Set("Content-Type", "application/json")
	respCalc, _ := http.DefaultClient.Do(reqCalc)
	assert.Equal(t, http.StatusOK, respCalc.StatusCode)

	reqList, _ := http.NewRequest("GET", srv.URL+"/api/v1/expressions", nil)
	reqList.Header.Set("Authorization", "Bearer "+payload.Token)
	respList, _ := http.DefaultClient.Do(reqList)
	var listResp struct {
		Expressions []struct {
			ID int `json:"id"`
		} `json:"expressions"`
	}
	json.NewDecoder(respList.Body).Decode(&listResp)
	id := listResp.Expressions[0].ID

	url := fmt.Sprintf("%s/api/v1/expressions/%d", srv.URL, id)
	reqGet, _ := http.NewRequest("GET", url, nil)
	reqGet.Header.Set("Authorization", "Bearer "+payload.Token)
	respGet, err := http.DefaultClient.Do(reqGet)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGet.StatusCode)

	var singleResp struct {
		Expression struct {
			ID     int     `json:"id"`
			Status string  `json:"status"`
			Result float64 `json:"result"`
		} `json:"expression"`
	}
	err = json.NewDecoder(respGet.Body).Decode(&singleResp)
	assert.NoError(t, err)
	assert.Equal(t, id, singleResp.Expression.ID)
}
