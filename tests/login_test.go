package tests

import (
	"bytes"
	"encoding/json"
	"gin-rest-api/models"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	response := func(req *http.Request) {
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer res.Body.Close()

		var data models.Response
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(res.Status, data)
	}

	// Test Case 1: Successful login
	t.Run("successful login", func(t *testing.T) {
		param := models.LoginRequest{Email: "user1@mail.com", Password: "123456"}
		body, _ := json.Marshal(param)

		req, err := http.NewRequest("POST", "http://localhost:8080/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		response(req)
	})

	// Test Case 2: Invalid credentials
	t.Run("invalid credentials", func(t *testing.T) {
		param := models.LoginRequest{Email: "user0@mail.com", Password: "123456"}
		body, _ := json.Marshal(param)

		req, err := http.NewRequest("POST", "http://localhost:8080/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		response(req)
	})

	// Test Case 3: Invalid request body (e.g., malformed JSON)
	t.Run("invalid request body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "http://localhost:8080/api/login", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		response(req)
	})

	// Test Case 4: Incorrect HTTP method
	t.Run("incorrect http method", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8080/api/login", nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		response(req)
	})
}
