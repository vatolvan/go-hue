package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ServerMock() (baseURL string, mux *http.ServeMux, teardownFn func()) {
	mux = http.NewServeMux()
	srv := httptest.NewServer(mux)
	return srv.URL, mux, srv.Close
}

func TestGetLights(t *testing.T) {
	baseURL, mux, teardown := ServerMock()
	defer teardown()

	var reqNum int
	mux.HandleFunc("/lights", func(w http.ResponseWriter, r *http.Request) {
		reqNum++
		// inspect request

		var resp map[string]HueLight
		resp = make(map[string]HueLight)
		resp["1"] = HueLight{Name: "light", State: HueLightState{On: true, Brightness: 254}}

		js, err := json.Marshal(resp)
		if err != nil {
			t.Errorf("Failed")
		}

		w.Write([]byte(js))
	})

	c := HueClient{BaseURL: baseURL}

	c.GetLights()

	if reqNum != 1 {
		t.Errorf("Did not call /lights")
	}
}

/*
func TestAbs(t *testing.T) {
	got := Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}
*/
