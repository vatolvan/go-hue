package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ServerMock() (baseURL string, mux *http.ServeMux, teardownFn func()) {
	mux = http.NewServeMux()
	srv := httptest.NewServer(mux)
	return srv.URL, mux, srv.Close
}

func findFirst(lights []Light, predicate func(Light) bool) (Light, error) {
	for _, light := range lights {
		if predicate(light) {
			return light, nil
		}
	}
	return Light{}, errors.New("No elements found")
}

func TestGetLights(t *testing.T) {
	baseURL, mux, teardown := ServerMock()
	defer teardown()

	var reqNum int

	var expectedHueLights map[string]HueLight
	expectedHueLights = make(map[string]HueLight)
	expectedLight1 := HueLight{Name: "light", State: HueLightState{On: true}}
	expectedLight2 := HueLight{Name: "light2", State: HueLightState{On: false}}
	expectedHueLights["1"] = expectedLight1
	expectedHueLights["12"] = expectedLight2

	mux.HandleFunc("/lights", func(w http.ResponseWriter, r *http.Request) {
		reqNum++

		js, err := json.Marshal(expectedHueLights)
		if err != nil {
			t.Errorf("Failed")
		}

		w.Write([]byte(js))
	})

	c := HueClient{BaseURL: baseURL}

	lights := c.GetLights()

	assert.Equal(t, 1, reqNum, "Wrong number of calls to /ligts!")
	assert.NotNil(t, lights, "Nil lights returned")
	assert.Equal(t, 2, len(lights), "Wrong number of lights returned")

	firstLight, err := findFirst(lights, func(light Light) bool {
		return light.ID == 1
	})
	assert.Nil(t, err)

	assert.Equal(t, expectedLight1.Name, firstLight.Name, "Wrong name for first light!")
	assert.True(t, firstLight.On, "First light should be on!")
	assert.Equal(t, 1, firstLight.ID, "Wrong ID for first light!")

	secondLight, err := findFirst(lights, func(light Light) bool {
		return light.ID == 12
	})
	assert.Nil(t, err)

	assert.Equal(t, expectedLight2.Name, secondLight.Name, "Wrong name for the second light!")
	assert.False(t, secondLight.On, "Second light should be off!")
	assert.Equal(t, 12, secondLight.ID, "Wrong ID for second light!")
}

func TestGetLight(t *testing.T) {
	baseURL, mux, teardown := ServerMock()
	defer teardown()

	var reqNum int

	expectedLight := HueLight{Name: "light", State: HueLightState{On: true}}

	mux.HandleFunc("/lights/1", func(w http.ResponseWriter, r *http.Request) {
		reqNum++

		js, err := json.Marshal(expectedLight)
		if err != nil {
			t.Errorf("Failed")
		}

		w.Write([]byte(js))
	})

	c := HueClient{BaseURL: baseURL}

	light := c.GetLight(1)

	assert.Equal(t, 1, reqNum, "Wrong number of calls to /lights/1")
	assert.NotNil(t, light)
	assert.Equal(t, expectedLight.Name, light.Name)
	assert.Equal(t, expectedLight.State.On, light.On)
	assert.Equal(t, 1, light.ID)
}

func TestSetLight(t *testing.T) {
	baseURL, mux, teardown := ServerMock()
	defer teardown()

	var state HueLightState

	var reqNum int
	mux.HandleFunc("/lights/1/state", func(w http.ResponseWriter, r *http.Request) {
		reqNum++
		err := json.NewDecoder(r.Body).Decode(&state)
		if err != nil {
			t.Errorf("Failed")
		}

		w.WriteHeader(200)
	})

	c := HueClient{BaseURL: baseURL}

	expectedState := HueLightState{On: true, Brightness: 100}
	c.SetLight(1, expectedState)

	assert.Equal(t, 1, reqNum, "lights/1/state not called")
	assert.NotNil(t, state, "Light state not received")
	assert.Equal(t, expectedState.On, state.On)
	assert.Equal(t, expectedState.Brightness, state.Brightness)
}
