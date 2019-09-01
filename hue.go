package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Light Describes the light returned from Hue Bridge
type Light struct {
	ID   string
	Name string
}

func getLights() []Light {
	username := viper.GetString("hue_bridge_username")
	baseURL := viper.GetString("hue_bridge_ip")

	URL := fmt.Sprintf("http://%s/api/%s/lights", baseURL, username)
	fmt.Println(URL)
	response, err := httpClient.Get(URL)

	if err != nil {
		log.Fatal(err)
	}

	var target map[string]interface{}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&target)
	// err = json.Unmarshal(response.Body, &target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(target)

	n := len(target)

	fmt.Println("Number of lights: ", n)

	var lights []Light
	lights = make([]Light, n)

	iter := 0
	for id, value := range target {
		light := (value.(map[string]interface{}))
		lights[iter] = Light{
			ID:   id,
			Name: light["name"].(string),
		}
		iter++
	}

	return lights
}
