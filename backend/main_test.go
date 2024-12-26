package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name         string
	method       string
	path         string
	body         map[string]string
	status       int
	disconnectDB bool
}{
	{
		name:   "Test create player",
		method: "POST",
		path:   "/player",
		body:   map[string]string{"name": "Player1"},
		status: http.StatusCreated,
	},
	{
		name:   "Test create player name missing",
		method: "POST",
		path:   "/player",
		body:   map[string]string{},
		status: http.StatusBadRequest,
	},
	{
		name:   "Test get all players",
		method: "GET",
		path:   "/players",
		body:   nil,
		status: http.StatusOK,
	},
	{
		name:   "Test delete player",
		method: "DELETE",
		path:   "/player/",
		body:   nil,
		status: http.StatusOK,
	},
	{
		name:   "Test delete not existing player",
		method: "DELETE",
		path:   "/player/999",
		body:   nil,
		status: http.StatusNotFound,
	},
	{
		name:         "Test create player db disconnected",
		method:       "POST",
		path:         "/player",
		body:         map[string]string{"name": "Player1"},
		status:       http.StatusInternalServerError,
		disconnectDB: true,
	},
	{
		name:   "Test get all players db disconnected",
		method: "GET",
		path:   "/players",
		body:   nil,
		status: http.StatusInternalServerError,
	},
}

func handleBody(body map[string]string) io.Reader {
	if b, err := json.Marshal(body); err == nil {
		return bytes.NewBuffer(b)
	}
	return nil
}

func handlePath(app *App, method string, path string) string {
	if method == "DELETE" {
		var player Player
		app.DB.First(&player)
		return path + fmt.Sprintf("%d", player.ID)
	}
	return path
}

func TestEndPoints(t *testing.T) {

	os.Setenv("DATABASE_URL", "postgres://myuser:mypassword@localhost:5432/mydb")

	app, err := createApp()
	if err != nil {
		t.Fatalf("Error creating test app: %v", err)
	}
	router := setupRouter(app)

	for _, tc := range tests {
		if tc.disconnectDB {
			dbInstance, _ := app.DB.DB()
			_ = dbInstance.Close()
		}
		t.Run(tc.name, func(t *testing.T) {
			body := handleBody(tc.body)
			path := handlePath(app, tc.method, tc.path)
			req, _ := http.NewRequest(tc.method, path, body)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Sprintln(w)
			assert.Equal(t, tc.status, w.Code)
		})
	}
	os.Unsetenv("DATABASE_URL")
}

func TestCreateAppNoURL(t *testing.T) {
	_, err := createApp()
	assert.Equal(t, err.Error(), "DATABASE_URL environment variable is not set")
}

func TestCreateAppWrongURL(t *testing.T) {
	os.Setenv("DATABASE_URL", "NOURL")
	_, err := createApp()
	assert.Equal(t, err.Error(), "failed to connect to the database: cannot parse `NOURL`: failed to parse as DSN (invalid dsn)")
	os.Unsetenv("DATABASE_URL")
}
