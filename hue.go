package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Light Describes the light returned from Hue Bridge
type Light struct {
	ID   int
	Name string
	On   bool
}

// HueLight Describes the JSON object returned by Hue system
type HueLight struct {
	Name  string        `json:"name"`
	State HueLightState `json:"state"`
}

// HueLightState Describes the state of the Hue Light
type HueLightState struct {
	// "on": false,
	// "bri": 254,
	// "hue": 8402,
	// "sat": 140,
	// "effect": "none",
	// "xy": [
	// 	0.4575,
	// 	0.4099
	// ],
	// "ct": 366,
	// "alert": "select",
	// "colormode": "xy",
	// "mode": "homeautomation",
	// "reachable": true
	On         bool `json:"on"`
	Brightness int  `json:"bri"`
}

func hueBaseURL() string {
	username := viper.GetString("hue_bridge_username")
	hueBridgeIP := viper.GetString("hue_bridge_ip")
	return fmt.Sprintf("http://%s/api/%s", hueBridgeIP, username)
}

func getLights() []Light {
	URL := fmt.Sprintf("%s/lights", hueBaseURL())
	response, err := httpClient.Get(URL)

	if err != nil {
		log.Fatal(err)
	}

	var target map[string]HueLight

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		log.Fatal(err)
	}

	var lights []Light
	lights = make([]Light, 0, len(target))

	for idStr, hueLight := range target {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		lights = append(lights, Light{
			ID:   id,
			Name: hueLight.Name,
			On:   hueLight.State.On,
		})
	}

	return lights
}

func getLight(id int) Light {
	URL := fmt.Sprintf("%s/lights/%d", hueBaseURL(), id)
	fmt.Println(URL)
	response, err := httpClient.Get(URL)

	if err != nil {
		log.Fatal(err)
	}

	var target HueLight

	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		log.Fatal(err)
	}

	return Light{
		ID:   id,
		Name: target.Name,
		On:   target.State.On,
	}
}

func setLight(id int, state HueLightState) {
	URL := fmt.Sprintf("%s/lights/%d/state", hueBaseURL(), id)

	requestBody, err := json.Marshal(state)
	if err != nil {
		log.Fatal("failed to serialize request: ", err)
	}

	request, err := http.NewRequest("PUT", URL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("failed to generate http PUT request: ", err)
	}

	_, err = httpClient.Do(request)
	if err != nil {
		log.Fatal("failed to set Hue state: ", err)
	}
}
