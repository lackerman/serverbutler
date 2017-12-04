package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestApiControllerIP(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/api/ip", nil)
	w := httptest.NewRecorder()

	api := new(ipController)
	api.get(w, req)

	res := w.Result()

	if res.StatusCode != 200 {
		t.Errorf("The IP handler failed to complete successfully")
	}
	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("The IP handler is not returning the correct header")
	}

	body, _ := ioutil.ReadAll(res.Body)
	var info map[string]interface{}
	if err := json.Unmarshal(body, &info); err != nil {
		t.Errorf("The IP handler failed to return a valid JSON response. %v\n%v", err.Error(), string(body))
	}
	if _, present := info["ip"]; !present {
		t.Errorf("The IP handler failed to return an IP")
	}
}
