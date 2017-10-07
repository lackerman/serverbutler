package controllers

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestApiControllerIP(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/api/ip", nil)
	w := httptest.NewRecorder()

	api := new(ipController)
	api.ip(w, req)

	res := w.Result()

	if res.StatusCode != 200 {
		t.Errorf("The IP handler failed to complete successfully")
	}
	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("The IP handler is not returning the correct header")
	}

	body, _ := ioutil.ReadAll(res.Body)

	returned := string(body)
	expected := `{ "ip": "localhost" }`
	if returned != expected {
		t.Errorf("The IP handler failed to return the correct response\n%v\n%v", returned, expected)
	}
}
