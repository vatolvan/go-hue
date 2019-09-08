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

// HueClient Handles the Hue system communication
type HueClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

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
	On         bool `json:"on"`
	Brightness int  `json:"bri"`
}

func hueBaseURL() string {
	username := viper.GetString("hue_bridge_username")
	hueBridgeIP := viper.GetString("hue_bridge_ip")
	return fmt.Sprintf("http://%s/api/%s", hueBridgeIP, username)
}

var defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}

// GetLights Return lights in the Hue system
func (c *HueClient) GetLights() []Light {
	httpClient, baseURL := c.getHTTPClient()

	URL := fmt.Sprintf("%s/lights", baseURL)
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

// GetLight Return Status of the light with specified id
func (c *HueClient) GetLight(id int) Light {
	httpClient, baseURL := c.getHTTPClient()

	URL := fmt.Sprintf("%s/lights/%d", baseURL, id)
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

// SetLight Sets the state of a Hue light with specified id
func (c *HueClient) SetLight(id int, state HueLightState) {
	httpClient, baseURL := c.getHTTPClient()

	URL := fmt.Sprintf("%s/lights/%d/state", baseURL, id)

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

func (c *HueClient) getHTTPClient() (*http.Client, string) {
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}

	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = hueBaseURL()
	}

	return httpClient, baseURL
}
