package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliworkshop/pgsql_slowest_queries/application"
	"github.com/aliworkshop/pgsql_slowest_queries/handler"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestApplication(t *testing.T) {
	app := application.NewApp()
	handlers := handler.NewHandler(app.GetDB())

	app.RegisterRoutes(handlers.Hello, "/", application.GET)
	app.RegisterRoutes(handlers.SlowestConnectionHandler, "/slowest_connection", application.GET)

	go func() {
		if err := app.Start(); err != nil {
			log.Fatal("error happened on app starting... : " + err.Error())
		}
	}()

	baseUrl := "http://localhost:3000"
	statusCode, _, err := Call("GET", baseUrl, nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	statusCode, _, err = Call("POST", baseUrl, nil)
	assert.NoError(t, err)
	assert.Equal(t, 405, statusCode)

	statusCode, slow, err := Call("GET", fmt.Sprintf("%s/%s", baseUrl, "slowest_connection"), map[string]string{
		"page":     "1",
		"per_page": "10",
		"order":    "desc",
	})
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Greater(t, slow[0].StateChange.Sub(*slow[0].QueryStart).Microseconds(), slow[1].StateChange.Sub(*slow[1].QueryStart).Microseconds())

	statusCode, slow, err = Call("GET", fmt.Sprintf("%s/%s", baseUrl, "slowest_connection"), map[string]string{
		"page":     "1",
		"per_page": "2",
		"order":    "desc",
	})
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, 2, len(slow))

	statusCode, slow, err = Call("GET", fmt.Sprintf("%s/%s", baseUrl, "slowest_connection"), map[string]string{
		"page":     "2",
		"per_page": "10",
		"order":    "desc",
	})
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, 0, len(slow))

	statusCode, slow, err = Call("GET", fmt.Sprintf("%s/%s", baseUrl, "slowest_connection"), map[string]string{
		"page":      "1",
		"per_page":  "10",
		"order":     "desc",
		"statement": "CREATE",
	})
	assert.NoError(t, err)
	assert.Equal(t, statusCode, 200)
	assert.Equal(t, 0, len(slow))

}

func Call(method string, path string, params map[string]string) (int, []handler.SlowQuery, error) {
	u, _ := url.Parse(path)
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	payload, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return 500, nil, err
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return 500, nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	cl := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := cl.Do(req)
	if err != nil {
		return 500, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp.StatusCode, nil, nil
	}

	data := []handler.SlowQuery{}
	json.NewDecoder(resp.Body).Decode(&data)
	return 200, data, nil
}
