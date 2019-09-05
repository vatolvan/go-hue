package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var hueClient = HueClient{}

func lightsHandler(w http.ResponseWriter, r *http.Request) {
	lights := hueClient.GetLights()

	js, err := json.Marshal(lights)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func lightHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	light := hueClient.GetLight(id)

	js, err := json.Marshal(light)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func lightSetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request HueLightState
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hueClient.SetLight(id, request)
	w.WriteHeader(200)
}

func readConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // path to look for the config file in
	viper.SetDefault("PORT", "8080")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	readConfig()

	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	r.HandleFunc("/lights/{id}", lightHandler).Methods("GET")
	r.HandleFunc("/lights/{id}", lightSetHandler).Methods("PUT")
	r.HandleFunc("/lights", lightsHandler).Methods("GET")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+viper.GetString("PORT"), nil))
}
